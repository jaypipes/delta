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
		return IsFloatTypeDef(subject)
	}
}

// IsFloatTypeDef returns true if the supplied subject is a type
// definition over `float32` or `float64.
func IsFloatTypeDef(subject any) bool {
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	if k == reflect.Float32 || k == reflect.Float64 {
		return true
	}
	return false
}
