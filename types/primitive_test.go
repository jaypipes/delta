package types_test

import (
	"fmt"
	"slices"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/stretchr/testify/require"
)

func TestIsPrimitive(t *testing.T) {
	notPrims := slices.Concat(
		[]any{
			testtypes.NilInterface,
			testtypes.NotNilInterface,
			struct{ a string }{a: "1"},
			func() {},
			nil,
			map[string]string{"a": "a"},
			[]string{"a"},
		},
		// Type definitions over built-in primitives are not primitives.
		testtypes.PrimitiveTypeDefZeroValues,
	)
	t.Run("not primitives", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range notPrims {
			require.False(
				types.IsPrimitive(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a primitive",
					subject, subject,
				),
			)
		}
	})

	t.Run("primitives", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range testtypes.PrimitiveZeroValues {
			require.True(
				types.IsPrimitive(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a primitive",
					subject, subject,
				),
			)
		}
	})
}

func TestIsPrimitiveTypeDef(t *testing.T) {
	notPrimTypeDefs := slices.Concat(
		[]any{
			testtypes.NilInterface,
			testtypes.NotNilInterface,
			struct{ a string }{a: "1"},
			func() {},
			nil,
			map[string]string{"a": "a"},
			[]string{"a"},
		},
		// Primitives are not primitive type definitions.
		testtypes.PrimitiveZeroValues,
	)
	t.Run("not primitive type defs", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range notPrimTypeDefs {
			require.False(
				types.IsPrimitiveTypeDef(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a primitive type definition",
					subject, subject,
				),
			)
		}
	})

	t.Run("primitive type defs", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range testtypes.PrimitiveTypeDefZeroValues {
			require.True(
				types.IsPrimitiveTypeDef(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a primitive type definition",
					subject, subject,
				),
			)
		}
	})
}

func TestUnderlyingPrimitive(t *testing.T) {
	notPrimTypeDefs := slices.Concat(
		[]any{
			testtypes.NilInterface,
			testtypes.NotNilInterface,
			struct{ a string }{a: "1"},
			func() {},
			nil,
			map[string]string{"a": "a"},
			[]string{"a"},
		},
		// Primitives are not primitive type definitions.
		testtypes.PrimitiveZeroValues,
	)
	t.Run("expect nil", func(t *testing.T) {
		require := require.New(t)
		for _, subject := range notPrimTypeDefs {
			require.Nil(
				types.UnderlyingPrimitive(subject),
				fmt.Sprintf(
					"expected %v (%T) to return nil",
					subject, subject,
				),
			)
		}
	})

	primTypeDefPairs := [][2]any{
		{testtypes.MyIntType(0), int(0)},
		{testtypes.MyInt8Type(0), int8(0)},
		{testtypes.MyInt16Type(0), int16(0)},
		{testtypes.MyInt32Type(0), int32(0)},
		{testtypes.MyInt64Type(0), int64(0)},
		{testtypes.MyUintType(0), uint(0)},
		{testtypes.MyUint8Type(0), uint8(0)},
		{testtypes.MyUint16Type(0), uint16(0)},
		{testtypes.MyUint32Type(0), uint32(0)},
		{testtypes.MyUint64Type(0), uint64(0)},
		{testtypes.MyFloat32Type(0), float32(0)},
		{testtypes.MyFloat64Type(0), float64(0)},
		{testtypes.MyStringType("b"), string("b")},
		{testtypes.MyByteArrayType("b"), []byte("b")},
		{testtypes.MyBoolType(false), bool(false)},
	}
	t.Run("exp underlying primitive", func(t *testing.T) {
		require := require.New(t)
		for _, pair := range primTypeDefPairs {
			subject := pair[0]
			exp := pair[1]
			got := types.UnderlyingPrimitive(subject)
			require.Equal(exp, got)
		}
	})
}
