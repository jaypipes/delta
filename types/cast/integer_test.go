package cast_test

import (
	"fmt"
	"testing"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types/cast"
	"github.com/stretchr/testify/assert"
)

func TestToUint64(t *testing.T) {
	notCastable := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		false,
		// strings and byte arrays are not parsed with cast
		"42",
		[]byte("42"),
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		[]string{"42"},
	}
	t.Run("not castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notCastable {
			_, ok := cast.ToUint64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be castable to uint64",
					subject, subject,
				),
			)
		}
	})

	castable := []any{
		// All integer types, regardless of sign or size
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		// floats are type cast
		float32(42.0),
		float64(42.0),
	}
	t.Run("castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range castable {
			got, ok := cast.ToUint64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be castable to uint64",
					subject, subject,
				),
			)
			assert.Equal(uint64(42), got)
		}
	})
}

func TestToSliceUint64(t *testing.T) {
	notCastable := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// scalars are not slices
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		false,
		// records are not slices
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// strings and byte arrays are not parsed with cast
		[]string{"42"},
		[][]byte{[]byte("42")},
		// slices of bools are not parsed with cast
		[]bool{false},
	}
	t.Run("not castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notCastable {
			_, ok := cast.ToSliceUint64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be castable to []uint64",
					subject, subject,
				),
			)
		}
	})

	castable := []any{
		// All integer types, regardless of sign or size
		[]int{int(42)},
		[]int16{int16(42)},
		[]uint{uint(42)},
		[]uint8{uint8(42)},
		// floats are type cast
		[]float32{float32(42.0)},
		[]float64{float64(42.0)},
	}
	t.Run("castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range castable {
			got, ok := cast.ToSliceUint64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be castable to []uint64",
					subject, subject,
				),
			)
			assert.Equal([]uint64{uint64(42)}, got)
		}
	})

	// Special cases

	// Even though this byte array is not probably not meant to be compared as
	// a slice of uint64, it *is* castable directly to a []uint64.
	exp := []uint64{110, 111, 116, 32, 105, 110, 116}
	got, ok := cast.ToSliceUint64([]byte("not int"))
	assert.True(t, ok)
	assert.Equal(t, exp, got)

}

func TestToInt64(t *testing.T) {
	notCastable := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		false,
		// strings and byte arrays are not parsed with cast
		"42",
		[]byte("42"),
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		[]string{"42"},
	}
	t.Run("not castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notCastable {
			_, ok := cast.ToInt64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be castable to int64",
					subject, subject,
				),
			)
		}
	})

	castable := []any{
		// All integer types, regardless of sign or size
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		// floats are type cast
		float32(42.0),
		float64(42.0),
	}
	t.Run("castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range castable {
			got, ok := cast.ToInt64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be castable to int64",
					subject, subject,
				),
			)
			assert.Equal(int64(42), got)
		}
	})
}

func TestToSliceInt64(t *testing.T) {
	notCastable := []any{
		nil,
		testtypes.NilInterface,
		testtypes.NotNilInterface,
		// scalars are not slices
		int(42),
		int16(42),
		uint(42),
		uint8(42),
		false,
		// records are not slices
		testtypes.MyStructType{},
		map[string]string{"foo": "bar"},
		// strings and byte arrays are not parsed with cast
		[]string{"42"},
		[][]byte{[]byte("42")},
		// slices of bools are not parsed with cast
		[]bool{false},
	}
	t.Run("not castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range notCastable {
			_, ok := cast.ToSliceInt64(subject)
			assert.False(
				ok, fmt.Sprintf(
					"%v (%T) should not be castable to []int64",
					subject, subject,
				),
			)
		}
	})

	castable := []any{
		// All integer types, regardless of sign or size
		[]int{int(42)},
		[]int16{int16(42)},
		[]uint{uint(42)},
		[]uint8{uint8(42)},
		// floats are type cast
		[]float32{float32(42.0)},
		[]float64{float64(42.0)},
	}
	t.Run("castable", func(t *testing.T) {
		assert := assert.New(t)
		for _, subject := range castable {
			got, ok := cast.ToSliceInt64(subject)
			assert.True(
				ok, fmt.Sprintf(
					"%v (%T) should be castable to []int64",
					subject, subject,
				),
			)
			assert.Equal([]int64{int64(42)}, got)
		}
	})

	// Special cases

	// Even though this byte array is not probably not meant to be compared as
	// a slice of int64, it *is* castable directly to a []int64.
	exp := []int64{110, 111, 116, 32, 105, 110, 116}
	got, ok := cast.ToSliceInt64([]byte("not int"))
	assert.True(t, ok)
	assert.Equal(t, exp, got)

}
