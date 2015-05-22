package pgx

import (
	"encoding/binary"
	"math"
	"time"
)

func (wb *WriteBuf) EncodeBool(v bool) {
	cast := 0
	if v {
		cast = 1
	}
	b := []byte{0, 0, 0, 1, cast}
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeInt2(v int16) {
	b := []byte{0, 0, 0, 2, byte(v >> 8), byte(v)}
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeInt4(v int32) {
	b := []byte{0, 0, 0, 4, byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeInt8(v int64) {
	b := []byte{0, 0, 0, 8}
	b = append(b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32))
	b = append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeFloat4(v float32) {
	b := []byte{0, 0, 0, 4}
	cast := int32(math.Float32bits(v))
	b = append(b, byte(cast>>24), byte(cast>>16), byte(cast>>8), byte(cast))
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeFloat8(v float64) {
	b := []byte{0, 0, 0, 8}
	cast := int32(math.Float64bits(v))
	b = append(b, byte(cast>>56), byte(cast>>48), byte(cast>>40), byte(cast>>32))
	b = append(b, byte(cast>>24), byte(cast>>16), byte(cast>>8), byte(cast))
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeBytea(v []byte) {
	totalLen := 4 + len(v)
	b := make([]byte, totalLen)
	binary.BigEndian.PutUint32(b, uint32(len(v)))
	if len(v) != 0 {
		copy(b[4:totalLen], v)
	}
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeText(v string) {
	totalLen := 4 + len(v)
	b := make([]byte, totalLen)
	binary.BigEndian.PutUint32(b, uint32(len(v)))
	if len(v) != 0 {
		copy(b[4:totalLen], v)
	}
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeVarchar(v string) {
	wb.EncodeText(v)
}

func (wb *WriteBuf) EncodeDate(v time.Time) {
	b := []byte{0, 0, 0, 10}
	wb.buf = append(wb.buf, append(b, v.Format("2006-01-02"))...)
}

func (wb *WriteBuf) EncodeTimestampTz(v time.Time) {
	microsecSinceUnixEpoch := v.Unix()*1000000 + int64(v.Nanosecond())/1000
	microsecSinceY2K := microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	x := microsecSinceY2K
	b := []byte{0, 0, 0, 8}
	b = append(b, byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32))
	b = append(b, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeTimestamp(v time.Time) {
	wb.EncodeTimestampTz(v)
}

func encodeArrayHeaderBytes(oid Oid, length, sizePerItem int) []byte {
	b := make([]byte, 24)
	binary.BigEndian.PutUint32(b[:4], int32(20+length*sizePerItem))
	binary.BigEndian.PutUint32(b[4:8], 1)               // number of dimensions
	binary.BigEndian.PutUint32(b[8:12], 0)              // no nulls
	binary.BigEndian.PutUint32(b[12:16], int32(oid))    // type of elements
	binary.BigEndian.PutUint32(b[16:20], int32(length)) // number of elements
	binary.BigEndian.PutUint32(b[20:24], 1)             // index of first element
	return b
}

func (wb *WriteBuf) EncodeBoolArray(vs []bool) {
	h := encodeArrayHeaderBytes(BoolOid, len(vs), 5)
	b := make([]byte, 0, 5*len(vs))
	for _, v := range vs {
		cast := 0
		if v {
			cast = 1
		}
		b = append(b, 0, 0, 0, 1, cast)
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeInt2Array(vs []int16) {
	h := encodeArrayHeaderBytes(Int2Oid, len(vs), 6)
	b := make([]byte, 0, 6*len(vs))
	for _, v := range vs {
		b = append(b, 0, 0, 0, 2, byte(v>>8), byte(v))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeInt4Array(vs []int32) {
	h := encodeArrayHeaderBytes(Int4Oid, len(vs), 8)
	b := make([]byte, 0, 8*len(vs))
	for _, v := range vs {
		b = append(b, 0, 0, 0, 4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeInt8Array(vs []int64) {
	h := encodeArrayHeaderBytes(Int8Oid, len(vs), 12)
	b := make([]byte, 0, 12*len(vs))
	for _, v := range vs {
		b = append(b, 0, 0, 0, 8)
		b = append(b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32))
		b = append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeFloat4Array(vs []float32) {
	h := encodeArrayHeaderBytes(Int4Oid, len(vs), 8)
	b := make([]byte, 0, 8*len(vs))
	for _, v := range vs {
		cast := int32(math.Float32bits(v))
		b = append(b, 0, 0, 0, 4, byte(cast>>24), byte(cast>>16), byte(cast>>8), byte(cast))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeFloat8Array(vs []float64) {
	h := encodeArrayHeaderBytes(Int8Oid, len(vs), 12)
	b := make([]byte, 0, 12*len(vs))
	for _, v := range vs {
		cast := int32(math.Float64bits(v))
		b = append(b, 0, 0, 0, 8)
		b = append(b, byte(cast>>56), byte(cast>>48), byte(cast>>40), byte(cast>>32))
		b = append(b, byte(cast>>24), byte(cast>>16), byte(cast>>8), byte(cast))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeTextArray(vs []string) {
	var totalStringSize int
	for _, v := range vs {
		totalStringSize += len(v)
	}
	size := 20 + len(vs)*4 + totalStringSize
	b := make([]byte, 24)
	binary.BigEndian.PutUint32(b[:4], size)
	binary.BigEndian.PutUint32(b[4:8], 1)                // number of dimensions
	binary.BigEndian.PutUint32(b[8:12], 0)               // no nulls
	binary.BigEndian.PutUint32(b[12:16], int32(TextOid)) // type of elements
	binary.BigEndian.PutUint32(b[16:20], int32(len(vs))) // number of elements
	binary.BigEndian.PutUint32(b[20:24], 1)              // index of first element

	for _, v := range slice {
		tmp := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp, len(v))
		b = append(b, append(tmp, v...)...)
	}

	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeVarcharArray(v []string) {
	var totalStringSize int
	for _, v := range vs {
		totalStringSize += len(v)
	}
	size := 20 + len(vs)*4 + totalStringSize
	b := make([]byte, 24)
	binary.BigEndian.PutUint32(b[:4], size)
	binary.BigEndian.PutUint32(b[4:8], 1)                   // number of dimensions
	binary.BigEndian.PutUint32(b[8:12], 0)                  // no nulls
	binary.BigEndian.PutUint32(b[12:16], int32(VarcharOid)) // type of elements
	binary.BigEndian.PutUint32(b[16:20], int32(len(vs)))    // number of elements
	binary.BigEndian.PutUint32(b[20:24], 1)                 // index of first element

	for _, v := range slice {
		tmp := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp, len(v))
		b = append(b, append(tmp, v...)...)
	}

	wb.buf = append(wb.buf, b)
}

func (wb *WriteBuf) EncodeTimestampArray(vs []time.Time) {
	h := encodeArrayHeaderBytes(TimestampOid, len(vs), 12)
	b := make([]byte, 0, 12*len(vs))
	for _, v := range vs {
		microsecSinceUnixEpoch := v.Unix()*1000000 + int64(v.Nanosecond())/1000
		microsecSinceY2K := microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
		x := microsecSinceY2K
		b = append(b, 0, 0, 0, 8)
		b = append(b, byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32))
		b = append(b, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
	}
	wb.buf = append(wb.buf, append(h, b...)...)
}

func (wb *WriteBuf) EncodeOid(v Oid) {
	b := []byte{0, 0, 0, 4, byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	wb.buf = append(wb.buf, b)
}
