package types_test

import (
	"fmt"
	"slices"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/stretchr/testify/assert"
)

func TestIsFloat(t *testing.T) {
	notFloats := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		int(42),
		false,
		"foo", // string is not byte array...
		map[string]string{"a": "a"},
		[]string{"a"},
		testtypes.MyStructType{A: 1},
		[]testtypes.MyStructType{testtypes.MyStructType{A: 1}},
		[]testtypes.MyIntType{1},
		testtypes.MySliceIntType{1},
	}
	t.Run("not floats", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notFloats {
			assert.False(
				types.IsFloat(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a float",
					subject, subject,
				),
			)
		}
	})

	floats := slices.Concat(
		testtypes.PrimitiveFloatZeroValues,
		testtypes.PrimitiveFloatTypeDefZeroValues,
	)

	t.Run("floats", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range floats {
			assert.True(
				types.IsFloat(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a float",
					subject, subject,
				),
			)
		}
	})
}
