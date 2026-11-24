package cast

// ToString takes an interface and returns the underlying value as a string
// along with a bool indicating whether the type cast conversion was
// safe/possible.
//
// If the supplied interface is not string or []byte, returns "", false.
func ToString(v any) (string, bool) {
	switch v := v.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	}
	return "", false
}
