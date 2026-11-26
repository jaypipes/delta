package types

import "github.com/jaypipes/delta/internal/genutil"

type MyIntType int
type MySliceIntType []int
type MySliceTypeDefIntType []MyIntType
type MyInt8Type int8
type MyInt16Type int16
type MyInt32Type int32
type MyInt64Type int64

type MyUintType uint
type MySliceUintType []uint
type MyUint8Type uint8
type MyUint16Type uint16
type MyUint32Type uint32
type MyUint64Type uint64

type MyFloat32Type float32
type MySliceFloat32Type []float32
type MyFloat64Type float64
type MySliceFloat64Type []float64

type MyStringType string
type MySliceStringType []string
type MyByteArrayType []byte
type MySliceByteArrayType [][]byte
type MySliceTypeDefByteArrayType []MyByteArrayType

type MyBoolType bool
type MySliceBoolType []bool

var PrimitiveTypeDefZeroValues = []any{
	// index 0 through 9 are integer types.
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
	// index 10 through 11 are float types.
	MyFloat32Type(9),
	MyFloat64Type(0),
	// index 12 through 13 are string types.
	MyStringType(""),
	MyByteArrayType(""),
	MyBoolType(false),
}

var SlicePrimitiveTypeDefZeroValues = []any{
	// index 0 through 9 are integer types.
	[]MyIntType{0},
	[]MyInt8Type{0},
	[]MyInt16Type{0},
	[]MyInt32Type{0},
	[]MyInt64Type{0},
	[]MyUintType{0},
	[]MyUint8Type{0},
	[]MyUint16Type{0},
	[]MyUint32Type{0},
	[]MyUint64Type{0},
	// index 10 through 11 are float types.
	[]MyFloat32Type{9},
	[]MyFloat64Type{0},
	// index 12 through 13 are string types.
	[]MyStringType{""},
	[]MyByteArrayType{MyByteArrayType("")},
	[]MyBoolType{false},
}

var UniquePairsPrimitiveTypeDefZeroValues = genutil.UniquePairs(PrimitiveTypeDefZeroValues)
var UniquePairsSlicePrimitiveTypeDefZeroValues = genutil.UniquePairs(SlicePrimitiveTypeDefZeroValues)
var UniquePairsPrimitiveTypeDefIntegerZeroValues = genutil.UniquePairs(PrimitiveTypeDefZeroValues[0:9])
var UniquePairsSlicePrimitiveTypeDefIntegerZeroValues = genutil.UniquePairs(SlicePrimitiveTypeDefZeroValues[0:9])
var UniquePairsPrimitiveTypeDefFloatZeroValues = genutil.UniquePairs(PrimitiveTypeDefZeroValues[10:11])
var UniquePairsSlicePrimitiveTypeDefFloatZeroValues = genutil.UniquePairs(SlicePrimitiveTypeDefZeroValues[10:11])
var UniquePairsPrimitiveTypeDefStringZeroValues = genutil.UniquePairs(PrimitiveTypeDefZeroValues[12:13])
var UniquePairsSlicePrimitiveTypeDefStringZeroValues = genutil.UniquePairs(SlicePrimitiveTypeDefZeroValues[12:13])
