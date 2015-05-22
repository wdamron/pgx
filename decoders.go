package pgx

import (
	"time"
)

func (vr *ValueReader) DecodeBool() (bool, error) {
	v := decodeBool(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt2() (int16, error) {
	v := decodeInt2(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt4() (int32, error) {
	v := decodeInt4(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt8() (int64, error) {
	v := decodeInt8(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeFloat4() (float32, error) {
	v := decodeFloat4(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeFloat8() (float64, error) {
	v := decodeFloat8(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeBytea() ([]byte, error) {
	v := decodeBytea(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeText() (string, error) {
	v := decodeText(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeVarchar() (string, error) {
	x := decodeText(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeDate() (time.Time, error) {
	v := decodeDate(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeTimestampTz() (time.Time, error) {
	v := decodeTimestampTz(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeTimestamp() (time.Time, error) {
	v := decodeTimestamp(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeBoolArray() ([]bool, error) {
	v := decodeBoolArray(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt2Array() ([]int16, error) {
	v := decodeInt2Array(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt4Array() ([]int32, error) {
	v := decodeInt4Array(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeInt8Array() ([]int64, error) {
	v := decodeInt8Array(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeFloat4Array() ([]float32, error) {
	v := decodeFloat4Array(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeFloat8Array() ([]float64, error) {
	v := decodeFloat8Array(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeTextArray() ([]string, error) {
	v := decodeTextArray(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeVarcharArray() ([]string, error) {
	v := decodeTextArray(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeTimestampArray() ([]time.Time, error) {
	v := decodeTimestampArray(vr)
	return v, vr.Err()
}

func (vr *ValueReader) DecodeOid() (Oid, error) {
	v := decodeOid(vr)
	return v, vr.Err()
}
