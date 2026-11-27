package cmp_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jaypipes/delta/cmp"
	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/stretchr/testify/assert"
)

func TestEqual_Invariants(t *testing.T) {

	f := func() {}
	c := make(chan struct{}, 1)
	notEquals := [][2]any{
		// if only a or b is `nil`, Equal always returns false
		{nil, 1},
		// nil is not the same as a nil interface.
		{testtypes.NilInterface, testtypes.NotNilInterface},
		// func and channel types are not comparable and supplying even the
		// exact same func or channel as both arguments to Equal will return
		// false.
		{f, 1},
		{f, f},
		{c, 1},
		{c, c},
	}
	t.Run("not equal", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range notEquals {
			a := pair[0]
			b := pair[1]
			assert.False(
				cmp.Equal(a, b),
				fmt.Sprintf(
					"[default options] expected %v (%T) "+
						"and %v (%T) not to be equal",
					a, a, b, b,
				),
			)
			// Equal is transitive, so reverse the argument order and verify
			// same result.
			assert.False(
				cmp.Equal(b, a),
				fmt.Sprintf(
					"[default options] expected %v (%T) "+
						"and %v (%T) not to be equal",
					b, b, a, a,
				),
			)
			assert.False(
				cmp.Equal(a, b, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"[no resolve type def] expected %v (%T) "+
						"and %v (%T) not to be equal",
					a, a, b, b,
				),
			)
			assert.False(
				cmp.Equal(b, a, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"[no resolve type def] expected %v (%T) and "+
						"%v (%T) not to be equal",
					b, b, a, a,
				),
			)
			assert.False(
				cmp.Equal(a, b, cmp.NoCast),
				fmt.Sprintf(
					"[no cast] expected %v (%T) "+
						"and %v (%T) not to be equal",
					a, a, b, b,
				),
			)
			assert.False(
				cmp.Equal(b, a, cmp.NoCast),
				fmt.Sprintf(
					"[no cast] expected %v (%T) and "+
						"%v (%T) not to be equal",
					b, b, a, a,
				),
			)
			assert.False(
				cmp.Equal(a, b, cmp.NoConvert),
				fmt.Sprintf(
					"[no convert] expected %v (%T) "+
						"and %v (%T) not to be equal",
					a, a, b, b,
				),
			)
			assert.False(
				cmp.Equal(b, a, cmp.NoConvert),
				fmt.Sprintf(
					"[no convert] expected %v (%T) and "+
						"%v (%T) not to be equal",
					b, b, a, a,
				),
			)
		}
	})

	equals := [][2]any{
		// if both a and b are `nil`, Equal always returns true
		{nil, nil},
	}

	t.Run("equal", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range equals {
			a := pair[0]
			b := pair[1]
			assert.True(
				cmp.Equal(a, b),
				fmt.Sprintf(
					"[default options] expected %v (%T) and "+
						"%v (%T) to be equal",
					a, a, b, b,
				),
			)
			assert.True(
				cmp.Equal(b, a),
				fmt.Sprintf(
					"[default options] expected %v (%T) and "+
						"%v (%T) to be equal",
					b, b, a, a,
				),
			)
			assert.True(
				cmp.Equal(a, b, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"[no resolve type def] expected %v (%T) and "+
						"%v (%T) to be equal",
					a, a, b, b,
				),
			)
			assert.True(
				cmp.Equal(b, a, cmp.NoResolveTypeDef),
				fmt.Sprintf(
					"[no resolve type def] expected %v (%T) and "+
						"%v (%T) to be equal",
					b, b, a, a,
				),
			)
			assert.True(
				cmp.Equal(a, b, cmp.NoCast),
				fmt.Sprintf(
					"[no cast] expected %v (%T) and "+
						"%v (%T) to be equal",
					a, a, b, b,
				),
			)
			assert.True(
				cmp.Equal(b, a, cmp.NoCast),
				fmt.Sprintf(
					"[no cast] expected %v (%T) and "+
						"%v (%T) to be equal",
					b, b, a, a,
				),
			)
			assert.True(
				cmp.Equal(a, b, cmp.NoConvert),
				fmt.Sprintf(
					"[no convert] expected %v (%T) and "+
						"%v (%T) to be equal",
					a, a, b, b,
				),
			)
			assert.True(
				cmp.Equal(b, a, cmp.NoConvert),
				fmt.Sprintf(
					"[no convert] expected %v (%T) and "+
						"%v (%T) to be equal",
					b, b, a, a,
				),
			)
		}
	})
}

// There are certain invariants when calling Equal regardless of what options
// you supply to Equal that these examples illustrate.
func ExampleEqual_invariants() {
	type myInterfaceType interface {
		F()
	}

	type myStructType struct {
		myInterfaceType
		A int
	}

	var nilInterface myInterfaceType = (*myStructType)(nil)

	// nil is always equal to nil
	fmt.Printf(
		"Is <nil> equal to <nil>? %t\n",
		cmp.Equal(nil, nil),
	)
	// but if only one of the subjects is nil, Equal returns false.
	fmt.Printf(
		"Is <nil> equal to 1? %t\n",
		cmp.Equal(nil, 1),
	)
	// An interface that references a nil is NOT the same as nil.
	fmt.Printf(
		"Is <nil> equal to an interface that references a nil pointer? %t\n",
		cmp.Equal(nil, nilInterface),
	)

	f := func() {}

	// Functions and channels are never comparable.
	fmt.Printf(
		"Is a func equal to the exact same func? %t\n",
		cmp.Equal(f, f),
	)

	c := make(chan struct{}, 1)

	// Functions and channels are never comparable.
	fmt.Printf(
		"Is a channel equal to the exact same channel? %t\n",
		cmp.Equal(c, c),
	)

	// Output:
	// Is <nil> equal to <nil>? true
	// Is <nil> equal to 1? false
	// Is <nil> equal to an interface that references a nil pointer? false
	// Is a func equal to the exact same func? false
	// Is a channel equal to the exact same channel? false
}

func TestEqual_TypeDefinitions(t *testing.T) {
	assert := assert.New(t)

	itd := testtypes.MyIntType(42)
	i := int(42)

	uitd := testtypes.MyUintType(42)
	ui := uint(42)

	assert.True(cmp.Equal(itd, i))
	assert.True(cmp.Equal(itd, i, cmp.NoCast))
	assert.True(cmp.Equal(itd, i, cmp.NoConvert))
	assert.False(cmp.Equal(itd, i, cmp.NoResolveTypeDef))

	assert.True(cmp.Equal(uitd, ui))
	assert.True(cmp.Equal(uitd, ui, cmp.NoCast))
	assert.True(cmp.Equal(uitd, ui, cmp.NoConvert))
	assert.False(cmp.Equal(uitd, ui, cmp.NoResolveTypeDef))

	// Two type definitions are resolved
	assert.True(cmp.Equal(itd, uitd))
	// resolving type definitions don't change the logic of no cast, which
	// cares about size and sign of the integer.
	assert.False(cmp.Equal(itd, uitd, cmp.NoCast))
	// If we can't resolve the type definitions, we cannot compare...
	assert.False(cmp.Equal(itd, uitd, cmp.NoResolveTypeDef))
}

// By default, Equal resolves any type definitions to the underlying primitive
// types and compares the underlying types. Passing [cmp.NoResolveTypeDef]
// disables this functionality and means that [Equal] will return false when
// comparing a variable with a primitive type to a variable of a the same
// underlying primitive type even if the underlying values are the same.
func ExampleEqual_type_definitions() {

	type myIntType int
	type myUintType uint

	itd := myIntType(42)
	i := int(42)

	uitd := myUintType(42)

	// type definitions over primitives are resolved to their underlying
	// primitive before comparison.
	fmt.Printf(
		"Is %v (%T) equal to %v (%T) when passing no options? %t\n",
		itd, itd, i, i,
		cmp.Equal(itd, i),
	)

	// But the same things are not equal when passing the
	// [cmp.NoResolveTypeDef] option.
	fmt.Printf(
		"Is %v (%T) equal to %v (%T) when passing cmp.NoResolveTypeDef? %t\n",
		itd, itd, i, i,
		cmp.Equal(itd, i, cmp.NoResolveTypeDef),
	)

	// Resolving type definitions happens before any casting or conversion of
	// the subjects, so passing [cmp.NoCast] when resolving type definitions
	// still results in the same comparison logic just on the underlying
	// primitive.
	fmt.Printf(
		"Is %v (%T) equal to %v (%T) when passing cmp.NoCast? %t\n",
		itd, itd, uitd, uitd,
		cmp.Equal(itd, uitd, cmp.NoCast),
	)

	// Output:
	// Is 42 (cmp_test.myIntType) equal to 42 (int) when passing no options? true
	// Is 42 (cmp_test.myIntType) equal to 42 (int) when passing cmp.NoResolveTypeDef? false
	// Is 42 (cmp_test.myIntType) equal to 42 (cmp_test.myUintType) when passing cmp.NoCast? false
}

func TestEqual_DefaultOptions(t *testing.T) {
	notEquals := [][2]any{
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
			{[]string{"1"}, []testtypes.MyIntType{testtypes.MyIntType(1)}},
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

// Calling Equal with no options causes the two subject types to be
// evaluated for comparability using conversion rules after the subjects have
// had any type definitions resolved to their underlying primitive types.
func ExampleEqual_default_options() {
	type myStructType struct {
		A int
	}

	type myIntType int
	type myUint64Type uint64

	comps := [][2]any{
		// Comparing two integer types with the default options ignores any
		// size and sign of the underlying integer primitive, so for all of the
		// following, Equal will return true..
		{0, int8(0)},
		{int(0), uint(0)},
		{myIntType(0), int(0)},
		{uint32(42), uint64(42)},
		{myUint64Type(42), uint64(42)},
		// A string is directly comparable to a []byte.
		{"1", []byte("1")},
		// slices are resolved to their underlying primitive
		{[]string{"a", "b"}, [][]byte{[]byte("a"), []byte("b")}},
		{map[string]string{"a": "b"}, 1},
		{[]byte("1"), string("1")},
		{string("1"), 1},
		{string("1"), uint(1)},
		{"0", myIntType(0)},
	}
	for _, pair := range comps {
		a := pair[0]
		b := pair[1]
		fmt.Printf(
			"Is %v (%T) equal to %v (%T)? %t\n",
			a, a, b, b,
			cmp.Equal(a, b),
		)
	}

	// Output:
	// Is 0 (int) equal to 0 (int8)? true
	// Is 0 (int) equal to 0 (uint)? true
	// Is 0 (cmp_test.myIntType) equal to 0 (int)? true
	// Is 42 (uint32) equal to 42 (uint64)? true
	// Is 42 (cmp_test.myUint64Type) equal to 42 (uint64)? true
	// Is 1 (string) equal to [49] ([]uint8)? true
	// Is [a b] ([]string) equal to [[97] [98]] ([][]uint8)? true
	// Is map[a:b] (map[string]string) equal to 1 (int)? false
	// Is [49] ([]uint8) equal to 1 (string)? true
	// Is 1 (string) equal to 1 (int)? true
	// Is 1 (string) equal to 1 (uint)? true
	// Is 0 (string) equal to 0 (cmp_test.myIntType)? true
}
