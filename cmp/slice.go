package cmp

import (
	"slices"

	"github.com/jaypipes/delta/types"
	"github.com/jaypipes/delta/types/cast"
	"github.com/jaypipes/delta/types/convert"
)

// comparableSlices returns true if the two supplied slices are comparable per
// the supplied Options.
func comparableSlices(a any, b any, opts Options) bool {
	aIsSliceScalar := types.IsSliceScalar(a)
	bIsSliceScalar := types.IsSliceScalar(b)

	if aIsSliceScalar && bIsSliceScalar {
		return comparableSliceScalars(a, b, opts)
	}

	// You cannot compare a slice of scalar to a slice of non-scalar, so if
	// only one argument is a slice of scalar, we don't need to go any further.
	if aIsSliceScalar || bIsSliceScalar {
		return false
	}

	if opts.NoResolveTypeDef || opts.NoCast || opts.NoConvert {
		return false
	}
	return convert.Convertible(a, b)
}

// comparableSlices returns true if the two supplied slices are comparable per
// the supplied Options.
func comparableSliceScalars(a any, b any, opts Options) bool {
	aIsSlicePrim := types.IsSlicePrimitive(a)
	bIsSlicePrim := types.IsSlicePrimitive(b)

	if aIsSlicePrim && bIsSlicePrim {
		return comparableSlicePrimitives(a, b, opts)
	}

	if opts.NoResolveTypeDef {
		return false
	}

	// One or both of the supplied slices of scalars are slices of primitive
	// type defs. Resolve to the slice of underlying primitive and compare
	// those.

	ap := a
	aIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(a)
	if aIsSlicePrimTypeDef {
		ap = types.SliceUnderlyingPrimitive(a)
	}

	bp := b
	bIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(b)
	if bIsSlicePrimTypeDef {
		bp = types.SliceUnderlyingPrimitive(b)
	}
	return comparableSlicePrimitives(ap, bp, opts)
}

// comparableSlicePrimitives returns true if the two supplied slices of
// primitives are comparable per the supplied Options.
func comparableSlicePrimitives(a any, b any, opts Options) bool {
	if opts.NoConvert {
		return comparableSlicePrimitivesNoConvert(a, b)
	}
	return convert.Convertible(a, b)
}

// comparableSlicePrimitivesNoConvert returns true if the two supplied slices
// of primitives are comparable with no type conversion.
func comparableSlicePrimitivesNoConvert(a any, b any) bool {
	switch a.(type) {
	case []int:
		switch b.(type) {
		case []int:
			return true
		default:
			return false
		}
	case []int8:
		switch b.(type) {
		case []int8:
			return true
		default:
			return false
		}
	case []int16:
		switch b.(type) {
		case []int16:
			return true
		default:
			return false
		}
	// Note that rune is a type alias for int32
	case []int32:
		switch b.(type) {
		case []int32:
			return true
		default:
			return false
		}
	case []int64:
		switch b.(type) {
		case []int64:
			return true
		default:
			return false
		}
	case []uint:
		switch b.(type) {
		case []uint:
			return true
		default:
			return false
		}
	// Note that byte is a type alias for uint8
	case []uint8:
		switch b.(type) {
		case []uint8:
			return true
		default:
			return false
		}
	case []uint16:
		switch b.(type) {
		case []uint16:
			return true
		default:
			return false
		}
	case []uint32:
		switch b.(type) {
		case []uint32:
			return true
		default:
			return false
		}
	case []uint64:
		switch b.(type) {
		case []uint64:
			return true
		default:
			return false
		}
	// []string is just a slice of immutable []byte
	case []string, [][]byte:
		switch b.(type) {
		case []string, [][]byte:
			return true
		default:
			return false
		}
	case []float32:
		switch b.(type) {
		case []float32:
			return true
		default:
			return false
		}
	case []float64:
		switch b.(type) {
		case []float64:
			return true
		default:
			return false
		}
	case []bool:
		switch b.(type) {
		case []bool:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// compareSlices returns 0 if the value of the supplied slices are equal per
// the supplied Options, -1 if the value of a is less than b and 1 if the value
// of a is greater than b.
func compareSlices(a any, b any, opts Options) int {
	//  TODO
	return -1
}

// equalSlices returns true if the two supplied slices are equal per the
// supplied Options.
func equalSlices(a any, b any, opts Options) bool {
	aIsSliceScalar := types.IsSliceScalar(a)
	bIsSliceScalar := types.IsSliceScalar(b)

	if aIsSliceScalar && bIsSliceScalar {
		return equalSliceScalars(a, b, opts)
	}

	// You cannot compare a slice of scalar to a slice of non-scalar, so if
	// only one argument is a slice of scalar, we don't need to go any further.

	if aIsSliceScalar || bIsSliceScalar {
		return false
	}

	// TODO)
	return false
}

// equalSliceScalars returns true if the two supplied slices of scalars are
// equal per the supplied Options.
func equalSliceScalars(a any, b any, opts Options) bool {
	aIsSlicePrim := types.IsSlicePrimitive(a)
	bIsSlicePrim := types.IsSlicePrimitive(b)

	if aIsSlicePrim && bIsSlicePrim {
		return equalSlicePrimitives(a, b, opts)
	}

	if opts.NoResolveTypeDef {
		return false
	}

	// One or both of the supplied slices of scalars are slices of primitive
	// type defs. Resolve to the slice of underlying primitive and compare
	// those.

	ap := a
	aIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(a)
	if aIsSlicePrimTypeDef {
		ap = types.SliceUnderlyingPrimitive(a)
	}

	bp := b
	bIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(b)
	if bIsSlicePrimTypeDef {
		bp = types.SliceUnderlyingPrimitive(b)
	}
	if ap == nil || bp == nil {
		// There is no easy way to differentiate between a string and a []byte.
		// A string is *not* a slice of primitive but a []byte *is* a slice of
		// primitive.
		//
		// So, if a is []string and b is MyByteArrayType, we fall into this
		// block because IsSlicePrimitive returns true for both []string and
		// MyByteArrayType but SliceUnderlyingPrimitive returns nil for
		// MyByteArrayType and []string for []string arguments.
		return false
	}
	return equalSlicePrimitives(ap, bp, opts)
}

// equalSlicePrimitives returns true if the two supplied primitives are
// equal per the supplied Options.
func equalSlicePrimitives(a any, b any, opts Options) bool {
	switch {
	case opts.NoCast:
		return equalSlicePrimitivesNoCast(a, b)
	case opts.NoConvert:
		return equalSlicePrimitivesCast(a, b)
	default:
		return equalSlicePrimitivesConvert(a, b)
	}
}

// equalSlicePrimitivesNoCast returns true if the two supplied slices of
// primitives are equal without any type casting or conversion.
func equalSlicePrimitivesNoCast(a any, b any) bool {
	switch a := a.(type) {
	case []int:
		switch b := b.(type) {
		case []int:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []int8:
		switch b := b.(type) {
		case []int8:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []int16:
		switch b := b.(type) {
		case []int16:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []int32:
		switch b := b.(type) {
		// Note that rune is a type alias for int32
		case []int32:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []int64:
		switch b := b.(type) {
		case []int64:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []uint:
		switch b := b.(type) {
		case []uint:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []uint8:
		switch b := b.(type) {
		// Note that byte is a type alias for uint8
		case []uint8:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []uint16:
		switch b := b.(type) {
		case []uint16:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []uint32:
		switch b := b.(type) {
		case []uint32:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []uint64:
		switch b := b.(type) {
		case []uint64:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []string:
		switch b := b.(type) {
		case []string:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		// string is just an immutable []byte
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			bs := make([]string, len(b))
			for x, el := range b {
				bs[x] = string(el)
			}
			slices.Sort(a)
			slices.Sort(bs)
			return slices.Equal(a, bs)
		default:
			return false
		}
	case [][]byte:
		switch b := b.(type) {
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			as := make([]string, len(a))
			bs := make([]string, len(b))
			for x := range len(a) {
				as[x] = string(a[x])
				bs[x] = string(b[x])
			}
			slices.Sort(as)
			slices.Sort(bs)
			return slices.Equal(as, bs)
		// string is just an immutable []byte
		case []string:
			if len(a) != len(b) {
				return false
			}
			as := make([]string, len(a))
			for x, el := range a {
				as[x] = string(el)
			}
			slices.Sort(as)
			slices.Sort(b)
			return slices.Equal(as, b)
		default:
			return false
		}
	case []float32:
		switch b := b.(type) {
		case []float32:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []float64:
		switch b := b.(type) {
		case []float64:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []bool:
		switch b := b.(type) {
		case []bool:
			return slices.Equal(a, b)
		default:
			return false
		}
	default:
		return false
	}
}

// equalSlicePrimitivesCast returns true if the two supplied slices of
// primitives are equal to each other after type casting to the largest size of
// the underlying Go type.
func equalSlicePrimitivesCast(a any, b any) bool {
	switch a := a.(type) {
	case []int, []int8, []int16, []int32, []int64:
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			ai, _ := cast.ToSliceInt64(a)
			bi, _ := cast.ToSliceInt64(b)
			if len(ai) != len(bi) {
				return false
			}
			slices.Sort(ai)
			slices.Sort(bi)
			return slices.Equal(ai, bi)
		default:
			return false
		}
	case []uint, []uint8, []uint16, []uint32, []uint64:
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			au, _ := cast.ToSliceUint64(a)
			bu, _ := cast.ToSliceUint64(b)
			if len(au) != len(bu) {
				return false
			}
			slices.Sort(au)
			slices.Sort(bu)
			return slices.Equal(au, bu)
		default:
			return false
		}
	case []string:
		switch b := b.(type) {
		case []string:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		// string is just an immutable []byte
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			bs := make([]string, len(b))
			for x, el := range b {
				bs[x] = string(el)
			}
			slices.Sort(a)
			slices.Sort(bs)
			return slices.Equal(a, bs)
		default:
			return false
		}
	case [][]byte:
		switch b := b.(type) {
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			as := make([]string, len(a))
			bs := make([]string, len(b))
			for x := range len(a) {
				as[x] = string(a[x])
				bs[x] = string(b[x])
			}
			slices.Sort(as)
			slices.Sort(bs)
			return slices.Equal(as, bs)
		// string is just an immutable []byte
		case []string:
			if len(a) != len(b) {
				return false
			}
			as := make([]string, len(a))
			for x, el := range a {
				as[x] = string(el)
			}
			slices.Sort(as)
			slices.Sort(b)
			return slices.Equal(as, b)
		default:
			return false
		}
	case []float32:
		switch b := b.(type) {
		case []float32:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []float64:
		switch b := b.(type) {
		case []float64:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []bool:
		switch b := b.(type) {
		case []bool:
			return slices.Equal(a, b)
		default:
			return false
		}
	default:
		return false
	}
}

// equalSlicePrimitivesConvert returns true if the two supplied slices of
// primitives are equal to each other after type converting the subjects to the
// type of the first subject.
func equalSlicePrimitivesConvert(a any, b any) bool {
	switch a := a.(type) {
	case []int, []int8, []int16, []int32, []int64:
		// Never an error for Integer types.
		ai, _ := cast.ToSliceInt64(a)
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			bi, _ := cast.ToSliceInt64(b)
			if len(ai) != len(bi) {
				return false
			}
			slices.Sort(ai)
			slices.Sort(bi)
			return slices.Equal(ai, bi)
		case
			[]string, [][]byte,
			[]float32, []float64:
			bi, ok := convert.ToSliceInt64(b)
			if !ok {
				return false
			}
			if len(ai) != len(bi) {
				return false
			}
			slices.Sort(ai)
			slices.Sort(bi)
			return slices.Equal(ai, bi)
		default:
			return false
		}
	case []uint, []uint8, []uint16, []uint32, []uint64:
		// Never an error for Integer types.
		au, _ := cast.ToSliceUint64(a)
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			bu, _ := cast.ToSliceUint64(b)
			if len(au) != len(bu) {
				return false
			}
			slices.Sort(au)
			slices.Sort(bu)
			return slices.Equal(au, bu)
		case
			[]string, [][]byte,
			[]float32, []float64:
			bu, ok := convert.ToSliceUint64(b)
			if !ok {
				return false
			}
			if len(au) != len(bu) {
				return false
			}
			slices.Sort(au)
			slices.Sort(bu)
			return slices.Equal(au, bu)
		default:
			return false
		}
	case []string:
		switch b := b.(type) {
		case []string:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		// string is just an immutable []byte
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			bs := make([]string, len(b))
			for x, el := range b {
				bs[x] = string(el)
			}
			slices.Sort(a)
			slices.Sort(bs)
			return slices.Equal(a, bs)
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64,
			[]bool:
			// Never an error converting slices of integer or bool to slices of
			// string.
			bs, _ := convert.ToSliceString(b)
			if len(a) != len(bs) {
				return false
			}
			slices.Sort(a)
			slices.Sort(bs)
			return slices.Equal(a, bs)
		default:
			return false
		}
	case [][]byte:
		as, _ := convert.ToSliceString(a)
		switch b := b.(type) {
		case [][]byte:
			if len(a) != len(b) {
				return false
			}
			bs := make([]string, len(b))
			for x := range len(a) {
				bs[x] = string(b[x])
			}
			slices.Sort(as)
			slices.Sort(bs)
			return slices.Equal(as, bs)
		// string is just an immutable []byte
		case []string:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(as)
			slices.Sort(b)
			return slices.Equal(as, b)
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64,
			[]bool:
			// Never an error converting slices of integer or bool to slices of
			// string.
			bs, _ := convert.ToSliceString(b)
			if len(a) != len(bs) {
				return false
			}
			slices.Sort(as)
			slices.Sort(bs)
			return slices.Equal(as, bs)
		default:
			return false
		}
	case []float32:
		switch b := b.(type) {
		case []float32:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []float64:
		switch b := b.(type) {
		case []float64:
			if len(a) != len(b) {
				return false
			}
			slices.Sort(a)
			slices.Sort(b)
			return slices.Equal(a, b)
		default:
			return false
		}
	case []bool:
		switch b := b.(type) {
		case []bool:
			return slices.Equal(a, b)
		default:
			return false
		}
	default:
		return false
	}
}
