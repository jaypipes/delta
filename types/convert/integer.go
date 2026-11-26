package convert

import (
	"strconv"

	"github.com/jaypipes/delta/types/cast"
)

// ToUint64 takes an interface and returns the underlying value as a uint64
// along with a bool indicating whether the type conversion was safe/possible.
func ToUint64(v any) (uint64, bool) {
	switch v := v.(type) {
	case
		uint, uint8, uint16, uint32, uint64,
		int, int8, int16, int32, int64,
		float32, float64:
		return cast.ToUint64(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return uint64(i), true
	case []byte:
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, false
		}
		return uint64(i), true
	}
	return 0, false
}

// ToSliceUint64 takes an interface and returns the underlying value as a
// []uint64 along with a bool indicating whether the type conversion was
// safe/possible.
func ToSliceUint64(subject any) ([]uint64, bool) {
	switch subject := subject.(type) {
	case []int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64:
		return cast.ToSliceUint64(subject)
	case []string:
		size := len(subject)
		out := make([]uint64, size)
		var ok bool
		for x, el := range subject {
			out[x], ok = ToUint64(el)
			if !ok {
				return nil, false
			}
		}
		return out, true
	case [][]byte:
		size := len(subject)
		out := make([]uint64, size)
		var ok bool
		for x, el := range subject {
			out[x], ok = ToUint64(el)
			if !ok {
				return nil, false
			}
		}
		return out, true
	default:
		return nil, false
	}
}

// ToInt64 takes an interface and returns the underlying value as an int64
// along with a bool indicating whether the type conversion was safe/possible.
func ToInt64(v any) (int64, bool) {
	switch v := v.(type) {
	case
		uint, uint8, uint16, uint32, uint64,
		int, int8, int16, int32, int64,
		float32, float64:
		return cast.ToInt64(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return int64(i), true
	case []byte:
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, false
		}
		return int64(i), true
	}
	return 0, false
}

// ToSliceInt64 takes an interface and returns the underlying value as a
// []int64 along with a bool indicating whether the type conversion was
// safe/possible.
func ToSliceInt64(subject any) ([]int64, bool) {
	switch subject := subject.(type) {
	case []int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64:
		return cast.ToSliceInt64(subject)
	case []string:
		size := len(subject)
		out := make([]int64, size)
		var ok bool
		for x, el := range subject {
			out[x], ok = ToInt64(el)
			if !ok {
				return nil, false
			}
		}
		return out, true
	case [][]byte:
		size := len(subject)
		out := make([]int64, size)
		var ok bool
		for x, el := range subject {
			out[x], ok = ToInt64(el)
			if !ok {
				return nil, false
			}
		}
		return out, true
	default:
		return nil, false
	}
}
