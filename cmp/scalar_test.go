package cmp_test

import (
	"fmt"
	"testing"

	"github.com/jaypipes/delta/cmp"
	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/stretchr/testify/assert"
)

func TestComparable_Scalar_TypeDefinitions(t *testing.T) {
	assert := assert.New(t)

	itd := testtypes.MyIntType(42)
	i := int(42)

	uitd := testtypes.MyUintType(42)
	ui := uint(42)

	assert.True(cmp.Comparable(itd, i))
	assert.True(cmp.Comparable(itd, i, cmp.NoCast))
	assert.True(cmp.Comparable(itd, i, cmp.NoConvert))
	assert.False(cmp.Comparable(itd, i, cmp.NoResolveTypeDef))

	assert.True(cmp.Comparable(uitd, ui))
	assert.True(cmp.Comparable(uitd, ui, cmp.NoCast))
	assert.True(cmp.Comparable(uitd, ui, cmp.NoConvert))
	assert.False(cmp.Comparable(uitd, ui, cmp.NoResolveTypeDef))

	// Two type definitions are resolved
	assert.True(cmp.Comparable(itd, uitd))
	// If we can't resolve the type definitions, we cannot compare...
	assert.False(cmp.Comparable(itd, uitd, cmp.NoResolveTypeDef))
	// resolving type definitions don't change the logic of no cast, which
	// cares about size and sign of the integer.
	assert.False(cmp.Comparable(itd, uitd, cmp.NoCast))
	assert.True(cmp.Comparable(itd, uitd, cmp.NoConvert))

	// Byte arrays are special since they are both scalars and slices.
	batd := testtypes.MyByteArrayType("foo")
	ba := []byte("foo")
	std := testtypes.MyStringType("foo")
	s := "foo"
	si := []int64{101, 111, 111}
	sitd := testtypes.MySliceIntType{101, 111, 111}

	assert.True(cmp.Comparable(batd, ba))
	assert.True(cmp.Comparable(batd, std))
	assert.True(cmp.Comparable(batd, s))
	// If we can't resolve the type definitions, we canot compare...
	assert.False(cmp.Comparable(batd, s, cmp.NoResolveTypeDef))
	assert.False(cmp.Comparable(batd, std, cmp.NoResolveTypeDef))
	// And they are compared as strings even when using no cast.
	assert.True(cmp.Comparable(batd, s, cmp.NoCast))
	assert.True(cmp.Comparable(batd, std, cmp.NoCast))
	// type definitions over byte arrays are comparable to slices of integer
	// primitives
	assert.True(cmp.Comparable(batd, si))
	assert.True(cmp.Comparable(batd, sitd))
	// but only if type definitions are resolved.
	assert.False(cmp.Comparable(batd, si, cmp.NoResolveTypeDef))
	assert.False(cmp.Comparable(batd, sitd, cmp.NoResolveTypeDef))
}

// By default, Comparable resolves any type definitions to the underlying primitive
// types and compares the underlying types. Passing [cmp.NoResolveTypeDef]
// disables this functionality and means that [Comparable] will return false when
// comparing a variable with a primitive type to a variable of a the same
// underlying primitive type even if the underlying values are the same.
func ExampleComparable_scalar_type_definitions() {
	type myIntType int
	type myUintType uint

	itd := myIntType(42)
	i := int(42)

	uitd := myUintType(42)

	// type definitions over primitives are resolved to their underlying
	// primitive before comparison.
	fmt.Printf(
		"Is %v (%T) comparable to %v (%T) when passing no options? %t\n",
		itd, itd, i, i,
		cmp.Comparable(itd, i),
	)

	// But the same things are not comparable when passing the
	// [cmp.NoResolveTypeDef] option.
	fmt.Printf(
		"Is %v (%T) comparable to %v (%T) when passing cmp.NoResolveTypeDef? %t\n",
		itd, itd, i, i,
		cmp.Comparable(itd, i, cmp.NoResolveTypeDef),
	)

	// Resolving type definitions happens BEFORE any casting or conversion of
	// the subjects, so passing [cmp.NoCast] when resolving type definitions
	// results in the same comparison logic just on the underlying primitive.
	fmt.Printf(
		"Is %v (%T) comparable to %v (%T) when passing cmp.NoCast? %t\n",
		itd, itd, uitd, uitd,
		cmp.Comparable(itd, uitd, cmp.NoCast),
	)

	// Resolving type definitions happens before any casting or conversion of
	// the subjects, so passing [cmp.NoConvert] when resolving type definitions
	// results in the same comparison logic just on the underlying
	// primitive.
	fmt.Printf(
		"Is %v (%T) comparable to %v (%T) when passing cmp.NoConvert? %t\n",
		itd, itd, uitd, uitd,
		cmp.Comparable(itd, uitd, cmp.NoConvert),
	)

	// Output:
	// Is 42 (cmp_test.myIntType) comparable to 42 (int) when passing no options? true
	// Is 42 (cmp_test.myIntType) comparable to 42 (int) when passing cmp.NoResolveTypeDef? false
	// Is 42 (cmp_test.myIntType) comparable to 42 (cmp_test.myUintType) when passing cmp.NoCast? false
	// Is 42 (cmp_test.myIntType) comparable to 42 (cmp_test.myUintType) when passing cmp.NoConvert? true
}

func TestComparable_Scalar_Integer(t *testing.T) {
	type tc struct {
		a    any
		b    any
		opts cmp.OptionList
		exp  bool
	}
	cases := []tc{}
	for _, pair := range testtypes.UniquePairsPrimitiveIntegerZeroValues {
		// integer primitives are comparable to all integer size and sign
		// when no options provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				exp: true,
			},
		)
		// Or when no convert option is provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoConvert},
				exp:  true,
			},
		)
		// But not when no cast option is provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoCast},
				exp:  false,
			},
		)
	}

	// String primitives can be compared to integer primitives as long as the
	// no cast or no convert options are passed.
	cases = append(
		cases,
		tc{
			a: "1", b: int(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: int(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: int(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: int(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: int(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: int(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)

	// Float primitives can be compared to integer primitives as long as the
	// no cast option is not passed.
	cases = append(
		cases,
		tc{
			a: float32(1), b: int(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: float32(1), b: int(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: float32(1), b: int(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  true,
		},
	)
	cases = append(
		cases,
		tc{
			a: float64(1), b: int(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: float64(1), b: int(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: float64(1), b: int(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  true,
		},
	)

	// Bool primitives can be compared to integer primitives as long as neither
	// the no cast or no convert options are provided.
	cases = append(
		cases,
		tc{
			a: false, b: int(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: int(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: int(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)

	assert := assert.New(t)
	for _, c := range cases {
		got := cmp.Comparable(c.a, c.b, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.a, c.a, c.b, c.b, c.opts.String(),
			),
		)
		// Comparable is transitive, so reverse the argument order and verify
		// same result.
		got = cmp.Comparable(c.b, c.a, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.b, c.b, c.a, c.a, c.opts.String(),
			),
		)
	}
}

func TestComparable_Scalar_Float(t *testing.T) {
	type tc struct {
		a    any
		b    any
		opts cmp.OptionList
		exp  bool
	}
	cases := []tc{}
	for _, pair := range testtypes.UniquePairsPrimitiveFloatZeroValues {
		// float primitives are comparable to all float sizes when no options
		// provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				exp: true,
			},
		)
		// Or when no convert option is provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoConvert},
				exp:  true,
			},
		)
		// But not when no cast option is provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoCast},
				exp:  false,
			},
		)
	}

	// String primitives can be compared to float primitives as long as the
	// no cast or no convert options are passed.
	cases = append(
		cases,
		tc{
			a: "1", b: float32(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: float32(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: float32(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float32(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float32(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float32(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: float64(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: float64(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: "1", b: float64(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float64(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float64(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: []byte("1"), b: float64(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)

	// Integer primitives can be compared to float primitives as long as the no
	// cast option is not passed.
	cases = append(
		cases,
		tc{
			a: int(0), b: float32(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: int(0), b: float32(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: int(0), b: float32(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  true,
		},
	)
	cases = append(
		cases,
		tc{
			a: int(0), b: float64(1),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: int(0), b: float64(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: int(0), b: float64(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  true,
		},
	)

	// Bool primitives can never be compared to float primitives.
	cases = append(
		cases,
		tc{
			a: false, b: float32(1),
			exp: false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: float32(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: float32(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: float64(1),
			exp: false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: float64(1),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: float64(1),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)

	assert := assert.New(t)
	for _, c := range cases {
		got := cmp.Comparable(c.a, c.b, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.a, c.a, c.b, c.b, c.opts.String(),
			),
		)
		// Comparable is transitive, so reverse the argument order and verify
		// same result.
		got = cmp.Comparable(c.b, c.a, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.b, c.b, c.a, c.a, c.opts.String(),
			),
		)
	}
}

func TestComparable_Scalar_String(t *testing.T) {
	type tc struct {
		a    any
		b    any
		opts cmp.OptionList
		exp  bool
	}
	cases := []tc{}
	for _, pair := range testtypes.UniquePairsPrimitiveStringZeroValues {
		// string primitives are comparable to other string primitives
		// regardless of options provided.
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				exp: true,
			},
		)
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoConvert},
				exp:  true,
			},
		)
		cases = append(
			cases,
			tc{
				a: pair[0], b: pair[1],
				opts: []cmp.Option{cmp.NoCast},
				exp:  true,
			},
		)
	}

	// Integer primitives can be compared to string primitives as long as
	// neither the no cast option or no convert options are passed.
	for _, izv := range testtypes.PrimitiveIntegerZeroValues {
		cases = append(
			cases,
			tc{
				a: izv, b: "0",
				exp: true,
			},
		)
		cases = append(
			cases,
			tc{
				a: izv, b: "0",
				opts: []cmp.Option{cmp.NoCast},
				exp:  false,
			},
		)
		cases = append(
			cases,
			tc{
				a: izv, b: "0",
				opts: []cmp.Option{cmp.NoConvert},
				exp:  false,
			},
		)
		cases = append(
			cases,
			tc{
				a: izv, b: []byte("0"),
				exp: true,
			},
		)
		cases = append(
			cases,
			tc{
				a: izv, b: []byte("0"),
				opts: []cmp.Option{cmp.NoCast},
				exp:  false,
			},
		)
		cases = append(
			cases,
			tc{
				a: izv, b: []byte("0"),
				opts: []cmp.Option{cmp.NoConvert},
				exp:  false,
			},
		)
	}

	// Bool primitives can be compared to string primitives only when neither
	// the no cast or no convert options are provided.
	cases = append(
		cases,
		tc{
			a: false, b: "false",
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: "false",
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: "false",
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: []byte("false"),
			exp: true,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: []byte("false"),
			opts: []cmp.Option{cmp.NoCast},
			exp:  false,
		},
	)
	cases = append(
		cases,
		tc{
			a: false, b: []byte("false"),
			opts: []cmp.Option{cmp.NoConvert},
			exp:  false,
		},
	)

	assert := assert.New(t)
	for _, c := range cases {
		got := cmp.Comparable(c.a, c.b, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.a, c.a, c.b, c.b, c.opts.String(),
			),
		)
		// Comparable is transitive, so reverse the argument order and verify
		// same result.
		got = cmp.Comparable(c.b, c.a, c.opts...)
		assert.Equal(
			c.exp, got,
			fmt.Sprintf(
				"expected %t for %v (%T) and %v (%T) "+
					"comparable using options %s",
				c.exp, c.b, c.b, c.a, c.a, c.opts.String(),
			),
		)
	}
}
