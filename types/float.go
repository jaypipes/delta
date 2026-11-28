package types

import "reflect"

// IsFloat returns true if the supplied subject is one of the built-in Go float
// types or a type definition over a built-in Go float type.
func IsFloat(subject any) bool {
	switch {
	case subject == nil:
		return false
	case IsPrimitiveFloat(subject):
		return true
	default:
		t := reflect.TypeOf(subject)
		k := t.Kind()
		return k == reflect.Float32 || k == reflect.Float64
	}
}
