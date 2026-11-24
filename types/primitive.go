package types

import (
	"reflect"
)

// IsPrimitive returns true if the supplied subject is one of the
// built-in primitive types:
//
// - `int`, `int8`, `int16`, `int32`, `int64`
// - `uint`, `uint8`, `uint16`, `uint32`, `uint64`
// - `string`, `[]byte`
// - `float32`, `float64`
// - `bool
//
// Note `nil` is not the same as an interface that references a nil pointer.
//
// Note `rune` and `int32` are the exact same type since `rune` is a type alias
// for `int32`.
//
// Note `byte` and `uint8` are the exact same type since `byte` is a type alias
// for `uint8`.
//
// Note `string` and `[]byte` are considered the exact same type since `string`
// is an immutable `[]byte`.
func IsPrimitive(subject any) bool {
	switch subject.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		string, []byte,
		float32, float64,
		bool:
		return true
	default:
		return false
	}
}

// IsPrimitiveInteger returns true if the supplied subject is one of the
// built-in Go integer types of any size and sign.
func IsPrimitiveInteger(subject any) bool {
	switch subject.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsPrimitiveSigned returns true if the supplied subject is one of the
// built-in signed Go integer types of any size.
func IsPrimitiveSigned(subject any) bool {
	switch subject.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

// IsPrimitiveFloat returns true if the supplied subject is one of the
// built-in Go float types of any size.
func IsPrimitiveFloat(subject any) bool {
	switch subject.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// IsPrimitiveString returns true if the supplied subject is either a `string`
// or `[]byte`.
func IsPrimitiveString(subject any) bool {
	switch subject.(type) {
	case string, []byte:
		return true
	default:
		return false
	}
}

// IsPrimitiveBool returns true if the supplied subject is a `bool`.
func IsPrimitiveBool(subject any) bool {
	switch subject.(type) {
	case bool:
		return true
	default:
		return false
	}
}

// IsPrimitiveTypeDef returns true if the supplied subject is a Go type
// definition that defines a built-in primitive as a custom type.
func IsPrimitiveTypeDef(subject any) bool {
	if subject == nil || IsPrimitive(subject) {
		return false
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	switch k {
	case
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	// string is just an immutable []byte
	case reflect.Slice:
		et := t.Elem()
		return et.Kind() == reflect.Uint8
	default:
		return false
	}
}

// PrimitiveTypeDefName returns the string name of the supplied subject's Go
// type definition that defines a built-in primitive as a custom type.
//
// If subject is not a type definition over a built-in primitive, returns "".
func PrimitiveTypeDefName(subject any) string {
	if !IsPrimitiveTypeDef(subject) {
		return ""
	}
	return reflect.TypeOf(subject).Name()
}

// UnderlyingPrimitive returns the underlying primitive of the supplied Go type
// type definition that defines a built-in primitive as a custom type.
//
// If the subject is not a primitive type definition or is itself nil, returns
// nil.
func UnderlyingPrimitive(subject any) any {
	if subject == nil || !IsPrimitiveTypeDef(subject) {
		return nil
	}
	t := reflect.TypeOf(subject)
	k := t.Kind()
	v := reflect.ValueOf(subject)
	switch k {
	case reflect.Bool:
		return v.Convert(reflect.TypeOf(false)).Interface()
	case reflect.Int:
		return v.Convert(reflect.TypeOf(int(0))).Interface()
	case reflect.Int8:
		return v.Convert(reflect.TypeOf(int8(0))).Interface()
	case reflect.Int16:
		return v.Convert(reflect.TypeOf(int16(0))).Interface()
	case reflect.Int32:
		return v.Convert(reflect.TypeOf(int32(0))).Interface()
	case reflect.Int64:
		return v.Convert(reflect.TypeOf(int64(0))).Interface()
	case reflect.Uint:
		return v.Convert(reflect.TypeOf(uint(0))).Interface()
	case reflect.Uint8:
		return v.Convert(reflect.TypeOf(uint8(0))).Interface()
	case reflect.Uint16:
		return v.Convert(reflect.TypeOf(uint16(0))).Interface()
	case reflect.Uint32:
		return v.Convert(reflect.TypeOf(uint32(0))).Interface()
	case reflect.Uint64:
		return v.Convert(reflect.TypeOf(uint64(0))).Interface()
	case reflect.Float32:
		return v.Convert(reflect.TypeOf(float32(0))).Interface()
	case reflect.Float64:
		return v.Convert(reflect.TypeOf(float64(0))).Interface()
	case reflect.String:
		return v.Convert(reflect.TypeOf("")).Interface()
	case reflect.Slice:
		et := t.Elem()
		if et.Kind() == reflect.Uint8 {
			return v.Convert(reflect.TypeOf([]byte(""))).Interface()
		}
		return nil
	default:
		return nil
	}
}
