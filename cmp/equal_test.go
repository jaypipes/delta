package cmp_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jaypipes/delta/cmp"
	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/stretchr/testify/assert"
)

func TestEqual_DefaultOptions(t *testing.T) {
	notEquals := [][2]any{
		// if only a or b is `nil`, Equal always returns false
		{nil, 1},
		// nil is not the same as a nil interface.
		{testtypes.NilInterface, testtypes.NotNilInterface},
		// func and channel types are not comparable to each other.
		{func() {}, 1},
		{make(chan int, 0), 1},
		// slices of primitives or primitive type definitions are
		// comparable only if the element type is comparable.
		{[]string{"1"}, []struct{ a int }{{1}}},
	}
	t.Run("not equal", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notEquals {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Equal(a, b),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be equal",
					a, a, b, b,
				),
			)
			// Equal is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Equal(b, a),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) not to be equal",
					b, b, a, a,
				),
			)
		}
	})

	equals := slices.Concat(
		testtypes.UniquePairsPrimitiveIntegerZeroValues,
		testtypes.UniquePairsPrimitiveFloatZeroValues,
		testtypes.UniquePairsPrimitiveStringZeroValues,
		testtypes.UniquePairsSlicePrimitiveIntegerZeroValues,
		testtypes.UniquePairsSlicePrimitiveFloatZeroValues,
		testtypes.UniquePairsSlicePrimitiveStringZeroValues,
		// default compare options treat Go type definitions as their
		// underlying primitive type.
		testtypes.UniquePairsPrimitiveTypeDefIntegerZeroValues,
		testtypes.UniquePairsPrimitiveTypeDefFloatZeroValues,
		testtypes.UniquePairsPrimitiveTypeDefStringZeroValues,
		testtypes.UniquePairsSlicePrimitiveTypeDefIntegerZeroValues,
		testtypes.UniquePairsSlicePrimitiveTypeDefFloatZeroValues,
		testtypes.UniquePairsSlicePrimitiveTypeDefStringZeroValues,
		[][2]any{
			// if both a and b are `nil`, Equal always returns true
			{nil, nil},
			// slices of primitives or primitive type definitions are
			// comparable if the underlying element type is comparable.
			{[]string{"1"}, []int{1}},
			{[]string{"1"}, [][]byte{[]byte("1")}},
			// Type definitions are first resolved to their underlying type and
			// then compared.
			{testtypes.MyIntType(0), 0},
			{testtypes.MyUintType(0), uint(0)},
			{testtypes.MyInt64Type(0), int64(0)},
			{testtypes.MyStringType("0"), int64(0)},
			{testtypes.MyByteArrayType("0"), int64(0)},
		},
	)

	t.Run("equal", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range equals {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Equal(a, b),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be equal",
					a, a, b, b,
				),
			)
			// Equal is transitive, so reverse the argument order and verify
			// same result.
			assert.True(
				cmp.Equal(b, a),
				fmt.Sprintf(
					"expected %v (%T) and %v (%T) to be equal",
					b, b, a, a,
				),
			)
		}
	})
}

func ExampleEqual() {
	s1 := "1"
	bs1 := []byte("1")

	fmt.Printf("Are s1 and bs1 equal? %t\n", cmp.Equal(s1, bs1))

	r1 := rune('1')
	i32 := int32('1')

	fmt.Printf("Are r1 and i32 equal? %t\n", cmp.Equal(r1, i32))

	b1 := byte('1')
	u8 := uint8('1')

	fmt.Printf("Are b1 and u8 equal? %t\n", cmp.Equal(b1, u8))

	i1 := 1
	ui1 := uint(1)

	fmt.Printf("Are i1 and ui1 equal? %t\n", cmp.Equal(i1, ui1))

	f32 := float32(1)
	f64 := float64(1)

	fmt.Printf("Are f32 and f64 equal? %t\n", cmp.Equal(f32, f64))
	// Output:
	// Are s1 and bs1 equal? true
	// Are r1 and i32 equal? true
	// Are b1 and u8 equal? true
	// Are i1 and ui1 equal? true
	// Are f32 and f64 equal? true
}
