package types

import "github.com/jaypipes/delta/internal/genutil"

var PrimitiveZeroValues = []any{
	// index 0 through 9 are integer types.
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	// index 10 through 11 are float types.
	float32(0),
	float64(0),
	// index 12 through 13 are string types.
	string(""),
	[]byte(""),
	bool(false),
}

var SlicePrimitiveZeroValues = []any{
	// index 0 through 9 are integer types.
	[]int{0},
	[]int8{0},
	[]int16{0},
	[]int32{0},
	[]int64{0},
	[]uint{0},
	[]uint8{0},
	[]uint16{0},
	[]uint32{0},
	[]uint64{0},
	// index 10 through 11 are float types.
	[]float32{0},
	[]float64{0},
	// index 12 through 13 are string types.
	[]string{""},
	[][]byte{[]byte("")},
	[]bool{false},
}

var UniquePairsPrimitiveZeroValues = genutil.UniquePairs(PrimitiveZeroValues)
var UniquePairsSlicePrimitiveZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues)
var PrimitiveIntegerZeroValues = PrimitiveZeroValues[0:9]
var PrimitiveSignedIntegerZeroValues = PrimitiveZeroValues[0:4]
var UniquePairsPrimitiveIntegerZeroValues = genutil.UniquePairs(PrimitiveIntegerZeroValues)
var UniquePairsPrimitiveSignedIntegerZeroValues = genutil.UniquePairs(PrimitiveSignedIntegerZeroValues)
var UniquePairsSlicePrimitiveIntegerZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[0:9])
var PrimitiveFloatZeroValues = PrimitiveZeroValues[10:11]
var UniquePairsPrimitiveFloatZeroValues = genutil.UniquePairs(PrimitiveFloatZeroValues)
var UniquePairsSlicePrimitiveFloatZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[10:11])
var PrimitiveStringZeroValues = PrimitiveZeroValues[12:13]
var UniquePairsPrimitiveStringZeroValues = genutil.UniquePairs(PrimitiveStringZeroValues)
var UniquePairsSlicePrimitiveStringZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[12:13])
