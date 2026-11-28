package types_test

import (
	"fmt"
	"slices"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/stretchr/testify/assert"
)

func TestIsInteger(t *testing.T) {
	notIntegers := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		false,
		float32(1.0),
		map[string]string{"a": "a"},
		[]string{"a"},
		testtypes.MyStructType{A: 1},
		[]testtypes.MyStructType{testtypes.MyStructType{A: 1}},
		[]testtypes.MyIntType{1},
		testtypes.MySliceIntType{1},
	}
	t.Run("not integers", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notIntegers {
			assert.False(
				types.IsInteger(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a integer",
					subject, subject,
				),
			)
		}
	})

	integers := slices.Concat(
		testtypes.PrimitiveIntegerZeroValues,
		testtypes.PrimitiveIntegerTypeDefZeroValues,
	)

	t.Run("integers", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range integers {
			assert.True(
				types.IsInteger(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a integer",
					subject, subject,
				),
			)
		}
	})
}

func TestIsSigned(t *testing.T) {
	notSigned := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		false,
		float32(1.0),
		uint(1),
		uint64(42),
		map[string]string{"a": "a"},
		[]string{"a"},
		testtypes.MyStructType{A: 1},
		[]testtypes.MyStructType{testtypes.MyStructType{A: 1}},
		[]testtypes.MyIntType{1},
		testtypes.MySliceIntType{1},
	}
	t.Run("not signed integers", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notSigned {
			assert.False(
				types.IsSignedInteger(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a signed integer",
					subject, subject,
				),
			)
		}
	})

	signed := slices.Concat(
		testtypes.PrimitiveSignedIntegerZeroValues,
		testtypes.PrimitiveSignedIntegerTypeDefZeroValues,
	)

	t.Run("signed integers", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range signed {
			assert.True(
				types.IsSignedInteger(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a signed integer",
					subject, subject,
				),
			)
		}
	})
}
