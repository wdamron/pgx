package pgx

import (
	"time"
)

// NextColumn returns the ValueReader for the next column in the current row, and
// a bool to signal whether the operation succeeded.
func (rows *Rows) NextColumn() (vr *ValueReader, ok bool) {
	return rows.nextColumn()
}

// A ParamsEncoder is used by the "raw" variants of Query/Exec, allowing the normal
// switch-based encoding process to be bypassed.
//
// A ParamsEncoder should be able to able to encode parameter formats and values for
// a given PreparedStatement into a WriteBuf.
type ParamsEncoder interface {
	// Encode the parameter formats for the given PreparedStatement (i.e. text/binary)
	// into the provided WriteBuf
	EncodeParamFormats(*PreparedStatement, *WriteBuf) error
	// Encode the parameters for the given PreparedStatement into the provided WriteBuf
	EncodeParams(*PreparedStatement, *WriteBuf) error
}

// ExecRaw executes sql, with parameters encoded by the given ParamsEncoder.
func (c *Conn) ExecRaw(sql string, pe ParamsEncoder) (commandTag CommandTag, err error) {
	startTime := time.Now()
	c.lastActivityTime = startTime

	if c.logger != dlogger {
		defer func() {
			if err == nil {
				endTime := time.Now()
				c.logger.Info("Exec", "sql", sql, "args", logQueryArgs(nil), "time", endTime.Sub(startTime), "commandTag", commandTag)
			} else {
				c.logger.Error("Exec", "sql", sql, "args", logQueryArgs(nil), "error", err)
			}
		}()
	}

	if err = c.sendQueryRaw(sql, pe); err != nil {
		return
	}

	var softErr error

	for {
		var t byte
		var r *msgReader
		t, r, err = c.rxMsg()
		if err != nil {
			return commandTag, err
		}

		switch t {
		case readyForQuery:
			c.rxReadyForQuery(r)
			return commandTag, softErr
		case rowDescription:
		case dataRow:
		case bindComplete:
		case commandComplete:
			commandTag = CommandTag(r.readCString())
		default:
			if e := c.processContextFreeMsg(t, r); e != nil && softErr == nil {
				softErr = e
			}
		}
	}
}

// QueryRaw executes sql, with parameters encoded by the given ParamsEncoder.
// If there is an error the returned *Rows will be returned in an error state,
// so it is allowed to ignore the error returned from Query and handle it in *Rows.
func (c *Conn) QueryRaw(sql string, pe ParamsEncoder) (*Rows, error) {
	c.lastActivityTime = time.Now()
	rows := &Rows{conn: c, startTime: c.lastActivityTime, sql: sql, args: nil, logger: c.logger}

	ps, ok := c.preparedStatements[sql]
	if !ok {
		var err error
		ps, err = c.Prepare("", sql)
		if err != nil {
			rows.abort(err)
			return rows, rows.err
		}
	}

	rows.fields = ps.FieldDescriptions
	err := c.sendPreparedQueryRaw(ps, pe)
	if err != nil {
		rows.abort(err)
	}
	return rows, rows.err
}

// QueryRowRaw is a convenience wrapper over QueryRaw.
// Any error that occurs while querying is deferred until calling Scan on the returned *Row.
// That *Row will error with ErrNoRows if no rows are returned.
func (c *Conn) QueryRowRaw(sql string, pe ParamsEncoder) *Row {
	rows, _ := c.QueryRaw(sql, pe)
	return (*Row)(rows)
}

func (c *Conn) sendQueryRaw(sql string, pe ParamsEncoder) (err error) {
	if ps, present := c.preparedStatements[sql]; present {
		return c.sendPreparedQueryRaw(ps, pe)
	} else {
		return c.sendSimpleQueryRaw(sql, pe)
	}
}

func (c *Conn) sendSimpleQueryRaw(sql string, pe ParamsEncoder) error {
	if pe == nil {
		wbuf := newWriteBuf(c.wbuf[0:0], 'Q')
		wbuf.WriteCString(sql)
		wbuf.closeMsg()

		_, err := c.conn.Write(wbuf.buf)
		if err != nil {
			c.die(err)
			return err
		}

		return nil
	}

	ps, err := c.Prepare("", sql)
	if err != nil {
		c.die(err)
		return err
	}

	return c.sendPreparedQueryRaw(ps, pe)
}

func (c *Conn) sendPreparedQueryRaw(ps *PreparedStatement, pe ParamsEncoder) error {
	// bind
	wbuf := newWriteBuf(c.wbuf[0:0], 'B')
	wbuf.WriteByte(0)
	wbuf.WriteCString(ps.Name)
	argLen := len(ps.ParameterOids)

	// (write format codes)
	wbuf.WriteInt16(int16(argLen))
	if err := pe.EncodeParamFormats(ps, wbuf); err != nil {
		c.die(err)
		return err
	}

	// (write parameters)
	wbuf.WriteInt16(int16(argLen))
	if err := pe.EncodeParams(ps, wbuf); err != nil {
		c.die(err)
		return err
	}

	wbuf.WriteInt16(int16(len(ps.FieldDescriptions)))
	for _, fd := range ps.FieldDescriptions {
		wbuf.WriteInt16(fd.FormatCode)
	}

	// execute
	wbuf.startMsg('E')
	wbuf.WriteByte(0)
	wbuf.WriteInt32(0)

	// sync
	wbuf.startMsg('S')
	wbuf.closeMsg()

	_, err := c.conn.Write(wbuf.buf)
	if err != nil {
		c.die(err)
	}

	return err
}
