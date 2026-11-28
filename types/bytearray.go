package types

import "reflect"

// IsByteArray returns true if the supplied subject is a `[]byte` or a type
// definition over `[]byte`
func IsByteArray(subject any) bool {
	if subject == nil {
		return false
	}
	if _, ok := subject.([]byte); ok {
		return true
	}
	return IsByteArrayTypeDef(subject)
}

// IsByteArrayTypeDef returns true if the supplied subject is a type
// definition over `[]byte`.
//
// Note that Both `string` and `MyByteArrayType` are scalars.  However, `string`
// is not a slice, but `MyByteArrayType` *is* a slice.
func IsByteArrayTypeDef(subject any) bool {
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	if k != reflect.Slice {
		return false
	}
	et := t.Elem()
	return et.Kind() == reflect.Uint8
}
