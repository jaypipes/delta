package cmp

import (
	"github.com/jaypipes/delta/types"
)

// Comparable returns true if the two supplied things are comparable according
// to any supplied options.
func Comparable(a any, b any, opts ...Option) bool {
	// We compare nils specially, so if either argument is nil, the pair is
	// comparable.
	if a == nil || b == nil {
		return true
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	aIsScalar := types.IsScalar(a)
	bIsScalar := types.IsScalar(b)

	aIsSlice := types.IsSlice(a)
	bIsSlice := types.IsSlice(b)

	// NOTE(jaypipes): A `[]byte` is BOTH a scalar and a slice type.

	switch {
	case aIsScalar && bIsScalar:
		return comparableScalars(a, b, o)
	case aIsSlice && bIsSlice:
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
