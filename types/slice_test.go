package types_test

import (
	"fmt"
	"slices"
	"sort"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestIsSliceScalar(t *testing.T) {
	notSliceScalars := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		map[string]string{"a": "a"},
	}
	t.Run("not slices of scalar", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notSliceScalars {
			assert.False(
				types.IsSliceScalar(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a slice of scalar",
					subject, subject,
				),
			)
		}
	})

	sliceScalars := slices.Concat(
		// All slices of primitives are slices of scalars
		testtypes.SlicePrimitiveZeroValues,
		// Slices of type definitions over built-in primitives are slices of
		// scalars.
		testtypes.SlicePrimitiveTypeDefZeroValues,
	)

	t.Run("slices of scalars", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range sliceScalars {
			assert.True(
				types.IsSliceScalar(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a slice of scalar",
					subject, subject,
				),
			)
		}
	})
}

func ExampleIsSliceScalar() {

	m := map[string]any{
		"nil":                   nil,
		"nil iface":             testtypes.NilInterface,
		"struct":                testtypes.MyStructType{},
		"func":                  func() {},
		"string":                "1",
		"[]string":              []string{"a", "b"},
		"[][]string":            [][]string{[]string{"a", "b"}},
		"map[string]string":     map[string]string{"a": "b"},
		"[]byte":                []byte("1"),
		"[][]byte":              [][]byte{[]byte("1")},
		"[]testtypes.MyIntType": []testtypes.MyIntType{testtypes.MyIntType(1)},
		"testtypes.MyIntType":   testtypes.MyIntType(1),
		"rune":                  rune('1'),
		"int32":                 int32('1'),
		"byte":                  byte('1'),
		"uint8":                 uint8('1'),
		"int":                   1,
		"uint":                  uint(1),
		"float32":               float32(1),
		"float64":               float64(1),
		"bool":                  true,
	}
	keys := lo.Keys(m)
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		fmt.Printf("Is %s a slice of scalar? %t\n", k, types.IsSliceScalar(v))
	}

	// Output:
	// Is [][]byte a slice of scalar? true
	// Is [][]string a slice of scalar? false
	// Is []byte a slice of scalar? true
	// Is []string a slice of scalar? true
	// Is []testtypes.MyIntType a slice of scalar? true
	// Is bool a slice of scalar? false
	// Is byte a slice of scalar? false
	// Is float32 a slice of scalar? false
	// Is float64 a slice of scalar? false
	// Is func a slice of scalar? false
	// Is int a slice of scalar? false
	// Is int32 a slice of scalar? false
	// Is map[string]string a slice of scalar? false
	// Is nil a slice of scalar? false
	// Is nil iface a slice of scalar? false
	// Is rune a slice of scalar? false
	// Is string a slice of scalar? false
	// Is struct a slice of scalar? false
	// Is testtypes.MyIntType a slice of scalar? false
	// Is uint a slice of scalar? false
	// Is uint8 a slice of scalar? false
}

func TestSliceUnderlyingPrimitive(t *testing.T) {
	notSliceTypeDefPrimitives := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		map[string]string{"a": "a"},
		[]testtypes.MyStructType{testtypes.MyStructType{}},
		// slice of primitives is not a slice of type def primitives
		[]string{"foo"},
		[]int{0},
		[]float64{0.0},
	}
	t.Run("not slices of type def primitives", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notSliceTypeDefPrimitives {
			got := types.SliceUnderlyingPrimitive(subject)
			assert.Nil(
				got,
				fmt.Sprintf(
					"expected %v (%T) not to be a slice "+
						"of type def primitives",
					subject, subject,
				),
			)
		}
	})

	sliceTypeDefPrimitives := [][2]any{
		{
			[]testtypes.MyIntType{testtypes.MyIntType(42)},
			[]int{42},
		},
		{
			[]testtypes.MyStringType{testtypes.MyStringType("foo")},
			[]string{"foo"},
		},
		{
			[]testtypes.MyFloat64Type{testtypes.MyFloat64Type(42.0)},
			[]float64{42.0},
		},
		{
			[]testtypes.MyFloat32Type{testtypes.MyFloat32Type(42.0)},
			[]float32{42.0},
		},
		// A type def of a byte array resolves to a []uint8
		{
			testtypes.MyByteArrayType("foo"),
			[]uint8{102, 111, 111},
		},
		// A type def of a slice of byte array resolves to a [][]uint8
		{
			[]testtypes.MyByteArrayType{testtypes.MyByteArrayType("foo")},
			[][]uint8{[]uint8{102, 111, 111}},
		},
		{
			testtypes.MySliceByteArrayType{[]byte("foo")},
			[][]uint8{[]uint8{102, 111, 111}},
		},
		{
			testtypes.MySliceTypeDefByteArrayType{testtypes.MyByteArrayType("foo")},
			[][]uint8{[]uint8{102, 111, 111}},
		},
		// A type def over a slice of primitives resolves to a slice of the
		// primitive
		{
			testtypes.MySliceIntType{42},
			[]int{42},
		},
		// A type def over a slice of type def primitives resolves to a slice
		// of the primitive.
		{
			testtypes.MySliceTypeDefIntType{testtypes.MyIntType(42)},
			[]int{42},
		},
	}

	t.Run("slices of scalars", func(t *testing.T) {
		assert := assert.New(t)
		for _, pair := range sliceTypeDefPrimitives {
			subject := pair[0]
			exp := pair[1]
			got := types.SliceUnderlyingPrimitive(subject)
			assert.Equal(exp, got)
		}
	})
}
