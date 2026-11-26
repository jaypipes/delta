package convert_test

import (
	"fmt"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types/convert"
	"github.com/stretchr/testify/assert"
)

func TestToUint64(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		"not an integer",
		[]byte("not an integer"),
		false,
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		[]string{"42"},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToUint64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to uint64",
					subject, subject,
				),
			)
		}
	})

	convertible := []any{
		// All integer types, regardless of sign or size
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		// strings and byte arrays are parsed
		"42",
		[]byte("42"),
		// floats are type cast
		float32(42.0),
		float64(42.0),
	}
	t.Run("convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range convertible {
			got, ok := convert.ToUint64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be convertible to uint64",
					subject, subject,
				),
			)
			assert.Equal(uint64(42), got)
		}
	})
}

func TestToSliceUint64(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// scalars are not slices
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		"not an integer",
		false,
		// records are not slices
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// slices but not convertible
		[]string{"not an integer"},
		[][]byte{[]byte("not an integer")},
		// slices of bools are not parsed with convert
		[]bool{false},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToSliceUint64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to []uint64",
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
		// strings and byte arrays are parsed
		[]string{"42"},
		[][]byte{[]byte("42")},
		// floats are type cast
		[]float32{float32(42.0)},
		[]float64{float64(42.0)},
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

	// Even though this byte array is not probably not meant to be compared as
	// a slice of uint64, it *is* convertible directly to a []uint64.
	exp := []uint64{110, 111, 116, 32, 105, 110, 116}
	got, ok := convert.ToSliceUint64([]byte("not int"))
	assert.True(t, ok)
	assert.Equal(t, exp, got)

}

func TestToInt64(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		"not an integer",
		[]byte("not an integer"),
		false,
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		[]string{"42"},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToInt64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to int64",
					subject, subject,
				),
			)
		}
	})

	convertible := []any{
		// All integer types, regardless of sign or size
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		// strings and byte arrays are parsed
		"42",
		[]byte("42"),
		// floats are type cast
		float32(42.0),
		float64(42.0),
	}
	t.Run("convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range convertible {
			got, ok := convert.ToInt64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be convertible to int64",
					subject, subject,
				),
			)
			assert.Equal(int64(42), got)
		}
	})
}

func TestToSliceInt64(t *testing.T) {
	notConvertible := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// scalars are not slices
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		"not an integer",
		false,
		// records are not slices
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// slices but not convertible
		[]string{"not an integer"},
		[][]byte{[]byte("not an integer")},
		// slices of bools are not parsed with convert
		[]bool{false},
	}
	t.Run("not convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notConvertible {
			_, ok := convert.ToSliceInt64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be convertible to []int64",
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
		// strings and byte arrays are parsed
		[]string{"42"},
		[][]byte{[]byte("42")},
		// floats are type cast
		[]float32{float32(42.0)},
		[]float64{float64(42.0)},
	}
	t.Run("convertible", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range convertible {
			got, ok := convert.ToSliceInt64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be convertible to []int64",
					subject, subject,
				),
			)
			assert.Equal([]int64{int64(42)}, got)
		}
	})

	// Special cases

	// Even though this byte array is not probably not meant to be compared as
	// a slice of int64, it *is* convertible directly to a []int64.
	exp := []int64{110, 111, 116, 32, 105, 110, 116}
	got, ok := convert.ToSliceInt64([]byte("not int"))
	assert.True(t, ok)
	assert.Equal(t, exp, got)

}
