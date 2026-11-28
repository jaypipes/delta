package types_test

import (
	"fmt"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
	"github.com/stretchr/testify/assert"
)

func TestIsByteArray(t *testing.T) {
	notByteArrays := []any{
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		struct{ a string }{a: "1"},
		func() {},
		nil,
		int(42),
		float32(42.0),
		false,
		"foo", // string is not byte array...
		map[string]string{"a": "a"},
		[]string{"a"},
		testtypes.MyStructType{A: 1},
		[]testtypes.MyStructType{testtypes.MyStructType{A: 1}},
		[]testtypes.MyIntType{1},
		testtypes.MySliceIntType{1},
	}
	t.Run("not byte arrays", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notByteArrays {
			assert.False(
				types.IsByteArray(subject),
				fmt.Sprintf(
					"expected %v (%T) not to be a byte array",
					subject, subject,
				),
			)
		}
	})

	byteArrays := []any{
		[]byte("foo"),
		testtypes.MyByteArrayType("foo"),
	}

	t.Run("byte arrays", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range byteArrays {
			assert.True(
				types.IsByteArray(subject),
				fmt.Sprintf(
					"expected %v (%T) to be a byte array",
					subject, subject,
				),
			)
		}
	})
}
