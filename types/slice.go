package types

import (
	"reflect"
)

// IsSlice returns true if the supplied subject is a slice of something.
func IsSlice(subject any) bool {
	if subject == nil {
		return false
	}
	return reflect.TypeOf(subject).Kind() == reflect.Slice
}

// IsSliceScalar returns true if the supplied subject is a slice of scalars.
func IsSliceScalar(subject any) bool {
	return IsSlicePrimitive(subject) || IsSlicePrimitiveTypeDef(subject)
}

// IsSliceSlice returns true if the supplied subject is a slice of slice of
// something. The one exception is a [][]byte, which is considered simply a
// slice of strings.
func IsSliceSlice(subject any) bool {
	if subject == nil || IsSliceScalar(subject) {
		return false
	}
	if IsSlice(subject) {
		return reflect.TypeOf(subject).Elem().Kind() == reflect.Slice
	}
	return false
}

// SliceScalarElementZeroValue returns the zero value of the supplied slice of
// primitive or slice of type definition of primitive's element.
//
// For example, if subject is `[]string{"a", "b"}`, returns `[]string{}`. If
// subject is `[]MyIntType{MyIntType(1), MyIntType(2)}`, returns
// `[]MyIntType{}`.
func SliceScalarElementZeroValue(subject any) any {
	if subject == nil || !IsSlice(subject) {
		return nil
	}
	if IsSlicePrimitive(subject) {
		return SliceScalarPrimitiveElementZeroValue(subject)
	}
	up := SliceUnderlyingPrimitive(subject)
	return SliceScalarPrimitiveElementZeroValue(up)
}

// SliceScalarPrimitiveElementZeroValue returns the zero value of the supplied slice of
// primitive.
//
// For example, if subject is `[]string{"a", "b"}`, returns `[]string{}`.
func SliceScalarPrimitiveElementZeroValue(subject any) any {
	switch subject.(type) {
	case []int:
		return []int{}
	case []int8:
		return []int8{}
	case []int16:
		return []int16{}
	case []int32:
		return []int32{}
	case []int64:
		return []int64{}
	case []uint:
		return []uint{}
	case []uint8:
		return []uint8{}
	case []uint16:
		return []uint16{}
	case []uint32:
		return []uint32{}
	case []uint64:
		return []uint64{}
	case []float32:
		return []float32{}
	case []float64:
		return []float64{}
	case [][]byte, []string:
		return []string{}
	case []bool:
		return []bool{}
	default:
		return nil
	}
}

// IsSlicePrimitive returns true if the supplied subject is a slice of
// one of the build-in scalar primitive types.
func IsSlicePrimitive(subject any) bool {
	switch subject.(type) {
	case
		[]int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64,
		[][]byte, []string,
		[]bool:
		return true
	default:
		return false
	}
}

// IsSlicePrimitiveTypeDef returns true if the supplied subject is a slice of
// Go type definition that defines a built-in primitive as a custom type.
func IsSlicePrimitiveTypeDef(subject any) bool {
	if !IsSlice(subject) || IsSlicePrimitive(subject) {
		return false
	}
	et := reflect.TypeOf(subject).Elem()
	ek := et.Kind()
	switch ek {
	case
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	case reflect.Uint8:
		// There isn't a way to differentiate between a `string` and a type
		// definition of a `[]byte`. So, we return true for a type definition
		// of `[]byte` here, even though we treat `[]byte` as a scalar.
		return true
	// []string is just a slice of immutable []byte
	case reflect.Slice:
		et := et.Elem()
		return et.Kind() == reflect.Uint8
	default:
		return false
	}
}

// SlicePrimitiveTypeDefZeroValue returns the zero value of the supplied slice
// of Go type definition of a built-in primitive.
func SlicePrimitiveTypeDefZeroValue(subject any) any {
	if subject == nil || !IsSlice(subject) {
		return nil
	}
	et := reflect.TypeOf(subject).Elem()
	ek := et.Kind()
	switch ek {
	case reflect.Bool:
		return []bool{}
	case reflect.Int:
		return []int{}
	case reflect.Int8:
		return []int8{}
	case reflect.Int16:
		return []int16{}
	case reflect.Int32:
		return []int32{}
	case reflect.Int64:
		return []int64{}
	case reflect.Uint:
		return []uint{}
	case reflect.Uint8:
		return []uint8{}
	case reflect.Uint16:
		return []uint16{}
	case reflect.Uint32:
		return []uint32{}
	case reflect.Uint64:
		return []uint64{}
	case reflect.Float32:
		return []float32{}
	case reflect.Float64:
		return []float64{}
	case reflect.String:
		return []string{}
	// []string is just a slice of immutable []byte
	case reflect.Slice:
		et := et.Elem()
		if et.Kind() == reflect.Slice {
			et := et.Elem()
			if et.Kind() == reflect.Uint8 {
				return []string{}
			}
		}
		return nil
	default:
		return nil
	}
}

// SlicePrimitiveTypeDefName returns the string name of the supplied slice
// subject element's Go type definition that defines a built-in primitive as a
// custom type.
//
// If subject is not a slice of type definition over a built-in primitive,
// returns "".
func SlicePrimitiveTypeDefName(subject any) string {
	if !IsSlicePrimitiveTypeDef(subject) {
		return ""
	}
	return reflect.TypeOf(subject).Elem().Name()
}

// SliceUnderlyingPrimitive returns a slice of the underlying primitive of the
// supplied slice of Go type type definition that defines a built-in primitive
// as a custom type.
//
// If the subject is not a slice of primitive type definition or is itself nil,
// returns nil.
//
// Note that a type definition over `[]byte` is a slice, a scalar AND a type
// definition. When calling this function with a type definitino over `[]byte`,
// this function returns a `[]uint8`.
func SliceUnderlyingPrimitive(subject any) any {
	// Special case for a type definition over []byte
	if IsByteArrayTypeDef(subject) {
		v := reflect.ValueOf(subject)
		size := v.Len()
		out := make([]uint8, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = uint8(el.Uint())
		}
		return out
	}
	if !IsSlicePrimitiveTypeDef(subject) {
		return nil
	}

	t := reflect.TypeOf(subject)
	et := t.Elem()
	ek := et.Kind()

	v := reflect.ValueOf(subject)
	size := v.Len()
	switch ek {
	case reflect.Uint8:
		// NOTE(jaypipes): It should not be possible to get here considering
		// IsByteArrayTypeDef and !IsSlicePrimitiveTypeDef would cover the case
		// of a type definition over `[]byte` or passing in `[]byte` directly.
		panic("logic failure getting slice underlying primitive")
	case reflect.Bool:
		out := make([]bool, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = el.Bool()
		}
		return out
	case reflect.Int:
		out := make([]int, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = int(el.Int())
		}
		return out
	case reflect.Int8:
		out := make([]int8, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = int8(el.Int())
		}
		return out
	case reflect.Int16:
		out := make([]int16, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = int16(el.Int())
		}
		return out
	case reflect.Int32:
		out := make([]int32, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = int32(el.Int())
		}
		return out
	case reflect.Int64:
		out := make([]int64, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = el.Int()
		}
		return out
	case reflect.Uint:
		out := make([]uint, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = uint(el.Uint())
		}
		return out
	case reflect.Uint16:
		out := make([]uint16, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = uint16(el.Uint())
		}
		return out
	case reflect.Uint32:
		out := make([]uint32, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = uint32(el.Uint())
		}
		return out
	case reflect.Uint64:
		out := make([]uint64, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = uint64(el.Uint())
		}
		return out
	case reflect.Float32:
		out := make([]float32, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = float32(el.Float())
		}
		return out
	case reflect.Float64:
		out := make([]float64, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = el.Float()
		}
		return out
	case reflect.String:
		out := make([]string, size)
		for x := 0; x < size; x++ {
			el := v.Index(x)
			out[x] = el.String()
		}
		return out
	case reflect.Slice:
		et = et.Elem()
		if et.Kind() == reflect.Uint8 {
			out := make([][]uint8, size)
			for x := 0; x < size; x++ {
				el := v.Index(x)
				if IsByteArrayTypeDef(el.Interface()) {
					// The element is type definition over a `[]byte`, so
					// convert the type def to a `[]uint8`.
					elBA := SliceUnderlyingPrimitive(el.Interface())
					if elBA == nil {
						return nil
					}
					out[x] = elBA.([]uint8)
				} else {
					// Not a type definition over `[]byte`. the element is a
					// `[]byte` itself, so just type cast.
					out[x] = el.Interface().([]uint8)
				}
			}
			return out
		}
		return nil
	default:
		return nil
	}
}
