package cast

// ToFloat64 takes an interface and returns the underlying value as a float64
// along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied interface is not an Integer or Float type, returns
// float64(0).
func ToFloat64(v any) (float64, bool) {
	switch v := v.(type) {
	case float64:
		return v, true
	case uint64:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint:
		return float64(v), true
	case int64:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	}
	return 0, false
}
