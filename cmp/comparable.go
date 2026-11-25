package cmp

import (
	"github.com/jaypipes/delta/types"
)

// Comparable returns true if the two supplied things are comparable according
// to any supplied options.
func Comparable(a any, b any, opts ...Option) bool {
	ac := types.CategoryOf(a)
	bc := types.CategoryOf(b)

	// We compare nils specially, so if either argument is nil, the pair is
	// comparable.
	if ac == types.Nil || bc == types.Nil {
		return true
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	aIsScalar := types.IsScalar(a)
	bIsScalar := types.IsScalar(b)

	// Both `string` and `MyByteArrayType` are scalars.  However, `string`
	// is not a slice, but `MyByteArrayType` *is* a slice.
	//
	// Here, we check to see if we are comparing a `[]string` (or
	// `[]MyByteArrayType`) with a `MyByteArrayType` and return false if
	// so.
	aIsBATD := types.IsByteArrayTypeDef(a)
	bIsBATD := types.IsByteArrayTypeDef(b)

	switch {
	case aIsScalar && bIsScalar:
		return comparableScalars(a, b, o)
	case aIsScalar && bIsBATD:
		return comparableScalars(a, b, o)
	case bIsScalar && aIsBATD:
		return comparableScalars(a, b, o)
	}

	aIsSlice := types.IsSlice(a)
	bIsSlice := types.IsSlice(b)

	if aIsSlice && bIsSlice {
		return comparableSlices(a, b, o)
	}

	// You cannot compare a slice to a non slice, so if only one argument is a
	// slice, we don't need to go any further.
	if aIsSlice || bIsSlice {
		return false
	}

	// If we get here, both a and b must be Record types.
	return comparableRecords(a, b, o)
}
