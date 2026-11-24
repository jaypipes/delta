package types

type MyIntType int
type MyInt8Type int8
type MyInt16Type int16
type MyInt32Type int32
type MyInt64Type int64

type MyUintType uint
type MyUint8Type uint8
type MyUint16Type uint16
type MyUint32Type uint32
type MyUint64Type uint64

type MyFloat32Type float32
type MyFloat64Type float64

type MyStringType string
type MyByteArrayType []byte

type MyBoolType bool

var PrimitiveTypeDefZeroValues = []any{
	MyIntType(0),
	MyInt8Type(0),
	MyInt16Type(0),
	MyInt32Type(0),
	MyInt64Type(0),
	MyUintType(0),
	MyUint8Type(0),
	MyUint16Type(0),
	MyUint32Type(0),
	MyUint64Type(0),
	MyFloat32Type(9),
	MyFloat64Type(0),
	MyStringType(""),
	MyByteArrayType(""),
	MyBoolType(false),
}
