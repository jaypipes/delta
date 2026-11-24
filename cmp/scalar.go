package cmp

import (
	"bytes"
	"strconv"

	"github.com/jaypipes/delta/cast"
	"github.com/jaypipes/delta/convert"
	"github.com/jaypipes/delta/types"
)

// comparableScalars returns true if the two supplied scalars are comparable
// per the supplied Options.
func comparableScalars(a any, b any, opts Options) bool {
	aIsPrim := types.IsPrimitive(a)
	bIsPrim := types.IsPrimitive(b)

	if aIsPrim && bIsPrim {
		return comparablePrimitives(a, b, opts)
	}

	// At least one of the scalars is a type definition over a primitive, so
	// resolve the type definition to its underlying primitive and compare the
	// resolved primitives.

	if opts.NoResolveTypeDef {
		return false
	}

	ap := a
	aIsPrimTypeDef := types.IsPrimitiveTypeDef(a)
	if aIsPrimTypeDef {
		ap = types.UnderlyingPrimitive(a)
	}

	bp := b
	bIsPrimTypeDef := types.IsPrimitiveTypeDef(b)
	if bIsPrimTypeDef {
		bp = types.UnderlyingPrimitive(b)
	}
	return comparablePrimitives(ap, bp, opts)
}

// comparablePrimitives returns true if the two supplied primitives are
// comparable per the supplied Options.
func comparablePrimitives(a any, b any, opts Options) bool {
	switch {
	case opts.NoCast:
		return comparablePrimitivesNoCast(a, b)
	case opts.NoConvert:
		return comparablePrimitivesCast(a, b)
	default:
		return convert.Convertible(a, b)
	}
}

// comparablePrimitivesNoCast returns true if the two supplied primitives are
// comparable with no type casting subjects to the greatest size of the
// underlying Go type before comparison.
func comparablePrimitivesNoCast(a any, b any) bool {
	switch a.(type) {
	case int:
		switch b.(type) {
		case int:
			return true
		default:
			return false
		}
	case int8:
		switch b.(type) {
		case int8:
			return true
		default:
			return false
		}
	case int16:
		switch b.(type) {
		case int16:
			return true
		default:
			return false
		}
	// Note that rune is a type alias for int32
	case int32:
		switch b.(type) {
		case int32:
			return true
		default:
			return false
		}
	case int64:
		switch b.(type) {
		case int64:
			return true
		default:
			return false
		}
	case uint:
		switch b.(type) {
		case uint:
			return true
		default:
			return false
		}
	// Note that byte is a type alias for uint8
	case uint8:
		switch b.(type) {
		case uint8:
			return true
		default:
			return false
		}
	case uint16:
		switch b.(type) {
		case uint16:
			return true
		default:
			return false
		}
	case uint32:
		switch b.(type) {
		case uint32:
			return true
		default:
			return false
		}
	case uint64:
		switch b.(type) {
		case uint64:
			return true
		default:
			return false
		}
	// string is just an immutable []byte
	case string, []byte:
		switch b.(type) {
		case string, []byte:
			return true
		default:
			return false
		}
	case float32:
		switch b.(type) {
		case float32:
			return true
		default:
			return false
		}
	case float64:
		switch b.(type) {
		case float64:
			return true
		default:
			return false
		}
	case bool:
		switch b.(type) {
		case bool:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// comparablePrimitivesNoConvert returns true if the two supplied primitives
// are comparable with type casting but no type conversion.
func comparablePrimitivesCast(a any, b any) bool {
	switch a.(type) {
	case
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return true
		default:
			return false
		}
	case float32:
		switch b.(type) {
		case float32:
			return true
		default:
			return false
		}
	case float64:
		switch b.(type) {
		case float64:
			return true
		default:
			return false
		}
	case string, []byte:
		switch b.(type) {
		case string, []byte:
			return true
		default:
			return false
		}
	case bool:
		switch b.(type) {
		case bool:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// compareScalars returns 0 if the value of the supplied subjects are equal per
// the supplied Options, -1 if the value of a is less than b and 1 if the value
// of a is greater than b.
func compareScalars(a any, b any, opts Options) int {
	//  TODO
	return -1
}

// equalScalars returns true if the two supplied scalars are equal
// per the supplied Options.
func equalScalars(a any, b any, opts Options) bool {
	aIsPrim := types.IsPrimitive(a)
	bIsPrim := types.IsPrimitive(b)

	if aIsPrim && bIsPrim {
		return equalPrimitives(a, b, opts)
	}

	if opts.NoResolveTypeDef {
		return false
	}

	ap := a
	aIsPrimTypeDef := types.IsPrimitiveTypeDef(a)
	if aIsPrimTypeDef {
		ap = types.UnderlyingPrimitive(a)
	}

	bp := b
	bIsPrimTypeDef := types.IsPrimitiveTypeDef(b)
	if bIsPrimTypeDef {
		bp = types.UnderlyingPrimitive(b)
	}
	return equalPrimitives(ap, bp, opts)
}

// equalPrimitives returns true if the two supplied primitives are
// equal per the supplied Options.
func equalPrimitives(a any, b any, opts Options) bool {
	switch {
	case opts.NoCast:
		return equalPrimitivesNoCast(a, b)
	case opts.NoConvert:
		return equalPrimitivesCast(a, b)
	default:
		return equalPrimitivesConvert(a, b, opts)
	}
}

// equalPrimitivesNoCast returns true if the two supplied primitives are equal
// with no type casting up to the largest size of the underlying type before
// comparison.
func equalPrimitivesNoCast(a any, b any) bool {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a == b
		default:
			return false
		}
	case int8:
		switch b := b.(type) {
		case int8:
			return a == b
		default:
			return false
		}
	case int16:
		switch b := b.(type) {
		case int16:
			return a == b
		default:
			return false
		}
	case int32:
		switch b := b.(type) {
		// Note that rune is a type alias for int32
		case int32:
			return a == b
		default:
			return false
		}
	case int64:
		switch b := b.(type) {
		case int64:
			return a == b
		default:
			return false
		}
	case uint:
		switch b := b.(type) {
		case uint:
			return a == b
		default:
			return false
		}
	case uint8:
		switch b := b.(type) {
		// Note that byte is a type alias for uint8
		case uint8:
			return a == b
		default:
			return false
		}
	case uint16:
		switch b := b.(type) {
		case uint16:
			return a == b
		default:
			return false
		}
	case uint32:
		switch b := b.(type) {
		case uint32:
			return a == b
		default:
			return false
		}
	case uint64:
		switch b := b.(type) {
		case uint64:
			return a == b
		default:
			return false
		}
	case string:
		switch b := b.(type) {
		case string:
			return a == b
		// string is just an immutable []byte
		case []byte:
			return bytes.Equal([]byte(a), b)
		default:
			return false
		}
	case []byte:
		switch b := b.(type) {
		case []byte:
			return bytes.Equal(a, b)
		// string is just an immutable []byte
		case string:
			return string(a) == b
		default:
			return false
		}
	case float32:
		switch b := b.(type) {
		case float32:
			return a == b
		default:
			return false
		}
	case float64:
		switch b := b.(type) {
		case float64:
			return a == b
		default:
			return false
		}
	case bool:
		switch b := b.(type) {
		case bool:
			return a == b
		default:
			return false
		}
	default:
		return false
	}
}

// equalPrimitivesCast returns true if the two supplied primitives are equal
// with type casting to the greatest size of the underlying Go type.
func equalPrimitivesCast(a any, b any) bool {
	switch a := a.(type) {
	case int, int8, int16, int32, int64:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			ai, _ := cast.ToInt64(a)
			bi, _ := cast.ToInt64(b)
			return ai == bi
		case uint, uint8, uint16, uint32, uint64:
			ai, _ := cast.ToInt64(a)
			bi, _ := cast.ToInt64(b)
			return ai == bi
		default:
			return false
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b := b.(type) {
		case int, int8, int16, int32, int64:
			av, _ := cast.ToUint64(a)
			bv, _ := cast.ToUint64(b)
			return av == bv
		case uint, uint8, uint16, uint32, uint64:
			av, _ := cast.ToUint64(a)
			bv, _ := cast.ToUint64(b)
			return av == bv
		default:
			return false
		}
	case float32:
		switch b := b.(type) {
		case float32:
			return a == b
		default:
			return false
		}
	case float64:
		switch b := b.(type) {
		case float64:
			return a == b
		default:
			return false
		}
	case string:
		switch b := b.(type) {
		case string:
			return a == b
		case []byte:
			return a == string(b)
		default:
			return false
		}
	case []byte:
		switch b := b.(type) {
		case string:
			return string(a) == b
		case []byte:
			return string(a) == string(b)
		default:
			return false
		}
	case bool:
		switch b := b.(type) {
		case bool:
			return a == b
		default:
			return false
		}
	default:
		return false
	}
}

// equalPrimitivesConvert returns true if the two supplied primitives are equal
// with type conversion.
func equalPrimitivesConvert(a any, b any, opts Options) bool {
	switch a := a.(type) {
	case int, int8, int16, int32, int64:
		ai, _ := cast.ToInt64(a)
		switch b := b.(type) {
		case
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			bi, _ := cast.ToInt64(b)
			return ai == bi
		case string, []byte:
			bs, _ := cast.ToString(b)
			bi, err := strconv.Atoi(bs)
			if err != nil {
				return false
			}
			return ai == int64(bi)
		case bool:
			return (ai == int64(0)) == b
		default:
			return false
		}
	case uint, uint8, uint16, uint32, uint64:
		au, _ := cast.ToUint64(a)
		switch b := b.(type) {
		case
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			bu, _ := cast.ToUint64(b)
			return au == bu
		case string, []byte:
			bs, _ := cast.ToString(b)
			bu, err := strconv.Atoi(bs)
			if err != nil {
				return false
			}
			return au == uint64(bu)
		case bool:
			return (au == uint64(0)) == b
		default:
			return false
		}
	case float32, float64:
		af, _ := cast.ToFloat64(a)
		switch b := b.(type) {
		case float32, float64:
			bf, _ := cast.ToFloat64(b)
			return af == bf
		case
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			bf, _ := cast.ToFloat64(b)
			return af == bf
		case string, []byte:
			// String and Float ScalarTypes are equivalent comparable. When
			// comparing String and Float scalars, [Equivalent] determines the
			// number of decimal places in the String scalar and uses
			// [fmt.Sprintf] to create a `string` from the underlying `float`
			// value and then compares the resulting strings to each other.
			bs, _ := cast.ToString(b)
			places, err := convert.DecimalPlacesInFloatString(bs)
			if err != nil {
				return false
			}
			as := convert.ToFloatString(af, places)
			return as == bs
		default:
			return false
		}
	case string, []byte:
		as, _ := cast.ToString(a)
		switch b := b.(type) {
		case string, []byte:
			bs, _ := cast.ToString(b)
			return as == bs
		case
			int, int8, int16, int32, int64:
			bi, _ := cast.ToInt64(b)
			bs := strconv.Itoa(int(bi))
			return as == bs
		case
			uint, uint8, uint16, uint32, uint64:
			bu, _ := cast.ToUint64(b)
			bs := strconv.Itoa(int(bu))
			return as == bs
		case float32, float64:
			bf, _ := cast.ToFloat64(b)
			places, err := convert.DecimalPlacesInFloatString(as)
			if err != nil {
				return false
			}
			bs := convert.ToFloatString(bf, places)
			return as == bs
		case bool:
			ab, err := strconv.ParseBool(as)
			if err != nil {
				return false
			}
			return ab == b
		default:
			return false
		}
	case bool:
		switch b := b.(type) {
		case bool:
			return a == b
		case string, []byte:
			bs, _ := cast.ToString(b)
			bb, err := strconv.ParseBool(bs)
			if err != nil {
				return false
			}
			return a == bb
		case
			int, int8, int16, int32, int64:
			bi, _ := cast.ToInt64(b)
			return a == (bi == int64(0))
		case
			uint, uint8, uint16, uint32, uint64:
			bu, _ := cast.ToUint64(b)
			return a == (bu == uint64(0))
		default:
			return false
		}
	default:
		return false
	}
}
