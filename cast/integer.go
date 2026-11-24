package cast

// ToUint64 takes an interface and returns the underlying value as a uint64
// along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied interface is not an Integer or Float type, returns
// uint64(0).
func ToUint64(v any) (uint64, bool) {
	switch v := v.(type) {
	case uint64:
		return v, true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int:
		return uint64(v), true
	case float32:
		return uint64(v), true
	case float64:
		return uint64(v), true
	}
	return 0, false
}

// ToSliceUint64 takes an interface and returns the underlying value as a
// []uint64 along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied inteface is not a slice of Integers, returns nil, false.
func ToSliceUint64(subject any) ([]uint64, bool) {
	switch subject := subject.(type) {
	case []uint64:
		return subject, true
	case []int:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []int8:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []int16:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []int32:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []int64:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []uint:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []uint8:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []uint16:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	case []uint32:
		out := make([]uint64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToUint64(el)
			out[x] = i
		}
		return out, true
	default:
		return nil, false
	}
}

// ToInt64 takes an interface and returns the underlying value as an int64
// along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied interface is not an Integer or Float type, returns
// int64(0).
func ToInt64(v any) (int64, bool) {
	switch v := v.(type) {
	case int64:
		return v, true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	}
	return 0, false
}

// ToSliceInt64 takes an interface and returns the underlying value as a
// []int64 along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied inteface is not a slice of Integers, returns nil, false.
func ToSliceInt64(subject any) ([]int64, bool) {
	switch subject := subject.(type) {
	case []int64:
		return subject, true
	case []int:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []int8:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []int16:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []int32:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []uint64:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to int64 or uint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []uint:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []uint8:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []uint16:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	case []uint32:
		out := make([]int64, len(subject))
		for x, el := range subject {
			// Integer types never error when type casting to uint64 or uuint64.
			i, _ := ToInt64(el)
			out[x] = i
		}
		return out, true
	default:
		return nil, false
	}
}
