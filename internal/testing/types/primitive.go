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
var UniquePairsPrimitiveIntegerZeroValues = genutil.UniquePairs(PrimitiveZeroValues[0:9])
var UniquePairsSlicePrimitiveIntegerZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[0:9])
var UniquePairsPrimitiveFloatZeroValues = genutil.UniquePairs(PrimitiveZeroValues[10:11])
var UniquePairsSlicePrimitiveFloatZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[10:11])
var UniquePairsPrimitiveStringZeroValues = genutil.UniquePairs(PrimitiveZeroValues[12:13])
var UniquePairsSlicePrimitiveStringZeroValues = genutil.UniquePairs(SlicePrimitiveZeroValues[12:13])
