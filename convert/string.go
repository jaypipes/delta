package convert

import (
	"strconv"

	"github.com/jaypipes/delta/cast"
)

// ToString takes an interface and returns the underlying value as a string
// along with a bool indicating whether the type conversion was safe/possible.
func ToString(v any) (string, bool) {
	switch v := v.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	case uint, uint8, uint16, uint32, uint64:
		u, _ := cast.ToUint64(v)
		return strconv.FormatUint(u, 10), true
	case int, int8, int16, int32, int64:
		i, _ := cast.ToInt64(v)
		return strconv.FormatInt(i, 10), true
	case float32, float64:
		// NOTE(jaypipes): You must use [ToFloatString] when comparing strings
		// and floats.
		return "", false
	case bool:
		return strconv.FormatBool(v), true
	}
	return "", false
}

// ToSliceString takes an interface and returns the underlying value as a
// []string along with a bool indicating whether the type conversion was
// safe/possible.
func ToSliceString(subject any) ([]string, bool) {
	switch subject := subject.(type) {
	case []string:
		return subject, true
	case [][]byte:
		out := make([]string, len(subject))
		for x, el := range subject {
			out[x] = string(el)
		}
		return out, true
	case []int:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []int8:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []int16:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []int32:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []int64:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []uint:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []uint8:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []uint16:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []uint32:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []uint64:
		out := make([]string, len(subject))
		for x, el := range subject {
			s, _ := ToString(el)
			out[x] = s
		}
		return out, true
	case []float32, []float64:
		// NOTE(jaypipes): You must use [ToSliceFloatString] when comparing
		// slices of strings and floats.
		return nil, false
	case []bool:
		out := make([]string, len(subject))
		for x, el := range subject {
			out[x] = strconv.FormatBool(el)
		}
		return out, true
	default:
		return nil, false
	}
}
