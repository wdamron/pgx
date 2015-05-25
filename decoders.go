package pgx

import (
	"time"
)

func (vr *ValueReader) DecodeBool() bool {
	return decodeBool(vr)
}

func (vr *ValueReader) DecodeInt2() int16 {
	return decodeInt2(vr)
}

func (vr *ValueReader) DecodeInt4() int32 {
	return decodeInt4(vr)
}

func (vr *ValueReader) DecodeInt8() int64 {
	return decodeInt8(vr)
}

func (vr *ValueReader) DecodeFloat4() float32 {
	return decodeFloat4(vr)
}

func (vr *ValueReader) DecodeFloat8() float64 {
	return decodeFloat8(vr)
}

func (vr *ValueReader) DecodeBytea() []byte {
	return decodeBytea(vr)
}

func (vr *ValueReader) DecodeText() string {
	return decodeText(vr)
}

func (vr *ValueReader) DecodeVarchar() string {
	return decodeText(vr)
}

func (vr *ValueReader) DecodeDate() time.Time {
	return decodeDate(vr)
}

func (vr *ValueReader) DecodeTimestampTz() time.Time {
	return decodeTimestampTz(vr)
}

func (vr *ValueReader) DecodeTimestamp() time.Time {
	return decodeTimestamp(vr)
}

func (vr *ValueReader) DecodeBoolArray() []bool {
	return decodeBoolArray(vr)
}

func (vr *ValueReader) DecodeInt2Array() []int16 {
	return decodeInt2Array(vr)
}

func (vr *ValueReader) DecodeInt4Array() []int32 {
	return decodeInt4Array(vr)
}

func (vr *ValueReader) DecodeInt8Array() []int64 {
	return decodeInt8Array(vr)
}

func (vr *ValueReader) DecodeFloat4Array() []float32 {
	return decodeFloat4Array(vr)
}

func (vr *ValueReader) DecodeFloat8Array() []float64 {
	return decodeFloat8Array(vr)
}

func (vr *ValueReader) DecodeTextArray() []string {
	return decodeTextArray(vr)
}

func (vr *ValueReader) DecodeVarcharArray() []string {
	return decodeTextArray(vr)
}

func (vr *ValueReader) DecodeTimestampArray() []time.Time {
	return decodeTimestampArray(vr)
}

func (vr *ValueReader) DecodeOid() Oid {
	return decodeOid(vr)
}
