package convert_test

import (
	"fmt"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types/convert"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// floats are not convertible to strings. Must use
		// [convert.ToFloatString]
		float32(42.0),
		float64(42.0),
		// Records are not convertible to string
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// Neither are slices
		[]string{"42"},
		[][]byte{[]byte("not an integer")},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToString(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to string",
					subject, subject,
				),
			)
		}
	})

	convertibleIntish := []any{
		// All integer types, regardless of sign or size are convertible to
		// string.
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		// strings and byte arrays are type cast
		"42",
		[]byte("42"),
	}
	t.Run("convertible intish", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range convertibleIntish {
			got, ok := convert.ToString(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be convertible to string",
					subject, subject,
				),
			)
			assert.Equal("42", got)
		}
	})

	// Special cases

	// Bools are converted to strings using [strconv.FormatBool]
	got, ok := convert.ToString(false)
	assert.True(t, ok)
	assert.Equal(t, "false", got)

	got, ok = convert.ToString(true)
	assert.True(t, ok)
	assert.Equal(t, "true", got)
}

func TestToSliceString(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// scalars are not slices
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		"foo",
		false,
		// records are not slices
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// slices but not convertible
		[]float32{float32(42.0)},
		[]float64{float64(42.0)},
		[]testtypes.MyStructType{testtypes.MyStructType{}},
		// slice of slice not convertible
		[][]string{[]string{"42"}},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToSliceString(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to []string",
					subject, subject,
				),
			)
		}
	})

	convertible := []any{
		// All integer types, regardless of sign or size
		[]int{int(42)},
		[]int16{int16(42)},
		[]uint{uint(42)},
		[]uint8{uint8(42)},
		// strings and byte arrays are type cast
		[]string{"42"},
		[][]byte{[]byte("42")},
	}
	t.Run("convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range convertible {
			got, ok := convert.ToSliceUint64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be convertible to []uint64",
					subject, subject,
				),
			)
			assert.Equal([]uint64{uint64(42)}, got)
		}
	})

	// Special cases

	// Bools are converted to strings using [strconv.FormatBool]
	got, ok := convert.ToSliceString([]bool{false, true})
	assert.True(t, ok)
	assert.Equal(t, []string{"false", "true"}, got)

	// Even though this byte array is probably not meant to be compared as a
	// slice of string, it *is* convertible directly to a []string since every
	// byte in the byte array is convertible to a string representation of the
	// uint8 value.
	exp := []string{"102", "111", "111"}
	got, ok = convert.ToSliceString([]byte("foo"))
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}
