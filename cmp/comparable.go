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

	aIsScalar := ac == types.Scalar
	bIsScalar := bc == types.Scalar

	if aIsScalar && bIsScalar {
		return comparableScalars(a, b, o)
	}

	// You cannot compare a scalar to a non-scalar, so if only one argument is
	// a scalar, we don't need to go any further.
	if aIsScalar || bIsScalar {
		return false
	}

	aIsSlice := ac == types.Slice
	bIsSlice := bc == types.Slice

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
