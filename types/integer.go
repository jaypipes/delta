package types

import "reflect"

// IsInteger returns true if the supplied subject is one of the built-in Go
// integer types or a type definition over a built-in Go integer type.
func IsInteger(subject any) bool {
	switch {
	case subject == nil:
		return false
	case IsPrimitiveInteger(subject):
		return true
	default:
		return IsIntegerTypeDef(subject)
	}
}

// IsSigned returns true if the supplied subject is one of the built-in signed
// Go integer types or a type definition over a built-in signed Go integer
// type.
func IsSignedInteger(subject any) bool {
	return IsPrimitiveSignedInteger(subject) || IsSignedIntegerTypeDef(subject)
}

// IsIntegerTypeDef returns true if the supplied subject is a type definition
// over one of the Go built-in integer types.
func IsIntegerTypeDef(subject any) bool {
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	switch k {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

// IsSignedIntegerTypeDef returns true if the supplied subject is a type
// definition over one of the Go built-in signed integer types.
func IsSignedIntegerTypeDef(subject any) bool {
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	switch k {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}
