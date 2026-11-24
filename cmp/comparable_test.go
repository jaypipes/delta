package cmp_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jaypipes/delta/cmp"
	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/stretchr/testify/assert"
)

func TestComparable_DefaultOptions(t *testing.T) {
	notComps := [][2]any{
		// nil is not the same as a nil interface.
		{testtypes.NilInterface, testtypes.NotNilInterface},
		// func and channel types are not comparable to each other.
		{func() {}, 1},
		{make(chan int, 0), 1},
		// slices of primitives or primitive type definitions are
		// comparable only if the element type is comparable.
		{[]string{"1"}, []struct{ a int }{{1}}},
	}
	t.Run("not comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notComps {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Comparable(a, b),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Comparable(b, a),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					b, b, a, a,
				),
			)
		}
	})

	comps := slices.Concat(
		testtypes.UniquePairsPrimitiveIntegerZeroValues,
		testtypes.UniquePairsPrimitiveFloatZeroValues,
		testtypes.UniquePairsPrimitiveStringZeroValues,
		[][2]any{
			// when either argument to Compare is nil, returns true.
			{nil, 1},
			// default compare options treat Go type definitions as their
			// underlying primitive type.
			{testtypes.MyIntType(0), 0},
			// Sign doesn't matter by default.
			{testtypes.MyIntType(0), uint(0)},
			// Neither does size.
			{testtypes.MyInt16Type(0), int64(0)},
			// slices of primitives or primitive type definitions are
			// comparable if the element type is comparable.
			{[]string{"1"}, []int{1}},
			{[]string{"1"}, [][]byte{[]byte("1")}},
		},
	)

	t.Run("comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range comps {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Comparable(a, b),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.True(
				cmp.Comparable(b, a),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					b, b, a, a,
				),
			)
		}
	})
}

func TestComparable_NoResolveTypeDef(t *testing.T) {
	notComps := [][2]any{
		// nil is not the same as a nil interface.
		{testtypes.NilInterface, testtypes.NotNilInterface},
		// func and channel types are not comparable to each other.
		{func() {}, 1},
		{make(chan int, 0), 1},
		// NoResolveTypeDef causes Go type definitions to not be comparable as
		// their underlying primitive type.
		{testtypes.MyIntType(0), 0},
		{testtypes.MyIntType(0), uint(0)},
		{testtypes.MyInt16Type(0), int64(0)},
	}
	t.Run("not comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notComps {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Comparable(a, b, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Comparable(b, a, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					b, b, a, a,
				),
			)
		}
	})

	comps := slices.Concat(
		testtypes.UniquePairsPrimitiveIntegerZeroValues,
		testtypes.UniquePairsPrimitiveFloatZeroValues,
		testtypes.UniquePairsPrimitiveStringZeroValues,
		[][2]any{
			// when either argument to Compare is nil, returns true.
			{nil, 1},
		},
	)
	t.Run("comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range comps {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Comparable(a, b, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.True(
				cmp.Comparable(b, a, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					b, b, a, a,
				),
			)
		}
	})
}

func TestComparable_NoCast(t *testing.T) {
	notComps := slices.Concat(
		// size and sign matter
		testtypes.UniquePairsPrimitiveIntegerZeroValues,
		testtypes.UniquePairsPrimitiveFloatZeroValues,
		[][2]any{
			// nil is not the same as a nil interface.
			{testtypes.NilInterface, testtypes.NotNilInterface},
			// func and channel types are not comparable to each other.
			{func() {}, 1},
			{make(chan int, 0), 1},
			// NoCast triggers exact integer sign and size matching, regardless
			// of whether the argument is a type definition.
			{int(0), int8(0)},
			{int(9), uint(0)},
			{int32(0), int64(0)},
			{testtypes.MyIntType(0), int8(0)},
			{testtypes.MyIntType(0), uint(0)},
			{testtypes.MyInt32Type(0), int64(0)},
		},
	)
	t.Run("not comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notComps {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Comparable(a, b, cmp.NoCast),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Comparable(b, a, cmp.NoCast),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					b, b, a, a,
				),
			)
		}
	})

	comps := slices.Concat(
		// string and []byte are considered the same even when NoCast.
		testtypes.UniquePairsPrimitiveStringZeroValues,
		[][2]any{
			// when either argument to Compare is nil, returns true.
			{nil, 1},
			// NoCast triggers exact integer sign and size matching, but
			// type definitions are first resolved to their underlying type.
			// of whether the argument is a type definition.
			{testtypes.MyIntType(0), 0},
			{testtypes.MyUintType(0), uint(0)},
			{testtypes.MyInt64Type(0), int64(0)},
		},
	)
	t.Run("comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range comps {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Comparable(a, b, cmp.NoCast),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.True(
				cmp.Comparable(b, a, cmp.NoCast),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					b, b, a, a,
				),
			)
		}
	})
}

func TestComparable_NoConvert(t *testing.T) {
	notComps := slices.Concat(
		[][2]any{
			// nil is not the same as a nil interface.
			{testtypes.NilInterface, testtypes.NotNilInterface},
			// func and channel types are not comparable to each other.
			{func() {}, 1},
			{make(chan int, 0), 1},
		},
	)
	t.Run("not comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notComps {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Comparable(a, b, cmp.NoConvert),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Comparable(b, a, cmp.NoConvert),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be comparable",
					b, b, a, a,
				),
			)
		}
	})

	comps := slices.Concat(
		// size and sign don't matter for NoConvert.
		testtypes.UniquePairsPrimitiveIntegerZeroValues,
		testtypes.UniquePairsPrimitiveFloatZeroValues,
		// string and []byte are considered the same.
		testtypes.UniquePairsPrimitiveStringZeroValues,
		[][2]any{
			// when either argument to Compare is nil, returns true.
			{nil, 1},
			// NoConvert uses type casting of different integer signs and sizes.
			{int(0), int8(0)},
			{int(9), uint(0)},
			{int32(0), int64(0)},
			// NoConvert triggers exact integer sign and size matching, but
			// type definitions are first resolved to their underlying type.
			// of whether the argument is a type definition.
			{testtypes.MyIntType(0), 0},
			{testtypes.MyUintType(0), uint(0)},
			{testtypes.MyInt64Type(0), int64(0)},
			{testtypes.MyIntType(0), testtypes.MyInt8Type(0)},
			{testtypes.MyIntType(0), testtypes.MyUintType(0)},
			{testtypes.MyInt32Type(0), testtypes.MyInt64Type(0)},
		},
	)
	t.Run("comparable", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range comps {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Comparable(a, b, cmp.NoConvert),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					a, a, b, b,
				),
			)
			// Compare is transitive, so reverse the argument order and verify
			// same result.
			assert.True(
				cmp.Comparable(b, a, cmp.NoConvert),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be comparable",
					b, b, a, a,
				),
			)
		}
	})
}

func ExampleComparable() {
	comps := [][2]any{
		{nil, 1},
		{testtypes.NilInterface, nil},
		{testtypes.MyStructType{}, nil},
		{func() {}, 1},
		{0, int8(0)},
		{int(0), uint(0)},
		{"1", []byte("1")},
		{testtypes.MyIntType(0), int(0)},
		{[]string{"a", "b"}, [][]byte{[]byte("a"), []byte("b")}},
		{map[string]string{"a": "b"}, 1},
		{[]byte("1"), string("1")},
		{string("1"), 1},
		{string("1"), uint(1)},
		{"0", testtypes.MyIntType(0)},
	}

	fmt.Println("=================================")
	fmt.Println("Using default compare options....")
	fmt.Println("=================================")
	for _, pair := range comps {
		a := pair[0]
		b := pair[1]
		fmt.Printf(
			"Is %T comparable to %T? %t\n",
			a, b,
			cmp.Comparable(a, b),
		)
	}

	fmt.Println("=================================")
	fmt.Println("Using NoResolveTypeDef ....")
	fmt.Println("=================================")
	for _, pair := range comps {
		a := pair[0]
		b := pair[1]
		fmt.Printf(
			"Is %T comparable to %T? %t\n",
			a, b,
			cmp.Comparable(a, b, cmp.NoResolveTypeDef),
		)
	}

	fmt.Println("=================================")
	fmt.Println("Using NoCast ....")
	fmt.Println("=================================")
	for _, pair := range comps {
		a := pair[0]
		b := pair[1]
		fmt.Printf(
			"Is %T comparable to %T? %t\n",
			a, b,
			cmp.Comparable(a, b, cmp.NoCast),
		)
	}

	fmt.Println("=================================")
	fmt.Println("Using NoConvert ....")
	fmt.Println("=================================")
	for _, pair := range comps {
		a := pair[0]
		b := pair[1]
		fmt.Printf(
			"Is %T comparable to %T? %t\n",
			a, b,
			cmp.Comparable(a, b, cmp.NoConvert),
		)
	}

	// Output:
	// =================================
	// Using default compare options....
	// =================================
	// Is <nil> comparable to int? true
	// Is *types.MyStructType comparable to <nil>? true
	// Is types.MyStructType comparable to <nil>? true
	// Is func() comparable to int? false
	// Is int comparable to int8? true
	// Is int comparable to uint? true
	// Is string comparable to []uint8? true
	// Is types.MyIntType comparable to int? true
	// Is []string comparable to [][]uint8? true
	// Is map[string]string comparable to int? false
	// Is []uint8 comparable to string? true
	// Is string comparable to int? true
	// Is string comparable to uint? true
	// Is string comparable to types.MyIntType? true
	// =================================
	// Using NoResolveTypeDef ....
	// =================================
	// Is <nil> comparable to int? true
	// Is *types.MyStructType comparable to <nil>? true
	// Is types.MyStructType comparable to <nil>? true
	// Is func() comparable to int? false
	// Is int comparable to int8? true
	// Is int comparable to uint? true
	// Is string comparable to []uint8? true
	// Is types.MyIntType comparable to int? false
	// Is []string comparable to [][]uint8? true
	// Is map[string]string comparable to int? false
	// Is []uint8 comparable to string? true
	// Is string comparable to int? true
	// Is string comparable to uint? true
	// Is string comparable to types.MyIntType? false
	// =================================
	// Using NoCast ....
	// =================================
	// Is <nil> comparable to int? true
	// Is *types.MyStructType comparable to <nil>? true
	// Is types.MyStructType comparable to <nil>? true
	// Is func() comparable to int? false
	// Is int comparable to int8? false
	// Is int comparable to uint? false
	// Is string comparable to []uint8? true
	// Is types.MyIntType comparable to int? true
	// Is []string comparable to [][]uint8? true
	// Is map[string]string comparable to int? false
	// Is []uint8 comparable to string? true
	// Is string comparable to int? false
	// Is string comparable to uint? false
	// Is string comparable to types.MyIntType? false
	// =================================
	// Using NoConvert ....
	// =================================
	// Is <nil> comparable to int? true
	// Is *types.MyStructType comparable to <nil>? true
	// Is types.MyStructType comparable to <nil>? true
	// Is func() comparable to int? false
	// Is int comparable to int8? true
	// Is int comparable to uint? true
	// Is string comparable to []uint8? true
	// Is types.MyIntType comparable to int? true
	// Is []string comparable to [][]uint8? true
	// Is map[string]string comparable to int? false
	// Is []uint8 comparable to string? true
	// Is string comparable to int? false
	// Is string comparable to uint? false
	// Is string comparable to types.MyIntType? false
}
