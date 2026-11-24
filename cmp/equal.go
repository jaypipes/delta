package cmp

import (
	"github.com/jaypipes/delta/types"
)

// Equal returns true if the value of a and b are equal per the supplied
// Options, false otherwise.
//
// If a and b are not comparable per the supplied Options, returns false.
//
// Callers should use [Comparable] with the same Options to check that two
// subjects are comparable before calling Equal.
//
// Callers should use [EqualComparable] for a function that returns both
// equality and comparability.
func Equal(a any, b any, opts ...Option) bool {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	if !Comparable(a, b, opts...) {
		return false
	}
	return equal(a, b, o)
}

// EqualComparable returns two bools. The first bool indicates whether the
// value of a and b are equal per the supplied Options. The second bool
// indicates whether a and b were comparable per the supplied Options.
func EqualComparable(a any, b any, opts ...Option) (bool, bool) {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	if !Comparable(a, b, opts...) {
		return false, false
	}
	return equal(a, b, o), true
}

// equal returns true if the two comparable subjects are equal per
// the supplied Options.
func equal(a any, b any, opts Options) bool {
	// two nils are always equal.
	if a == nil && b == nil {
		return true
	}
	// conversely, if only one of the two subjects is nil, they cannot be
	// equal.
	if a == nil || b == nil {
		return false
	}

	aIsScalar := types.IsScalar(a)
	bIsScalar := types.IsScalar(b)

	if aIsScalar && bIsScalar {
		return equalScalars(a, b, opts)
	}

	// You cannot compare a scalar to a non-scalar, so if only one argument is
	// a scalar, we don't need to go any further.
	if aIsScalar || bIsScalar {
		return false
	}

	aIsSlice := types.IsSlice(a)
	bIsSlice := types.IsSlice(b)

	if aIsSlice && bIsSlice {
		return equalSlices(a, b, opts)
	}

	// You cannot compare a slice to a non slice, so if only one argument is a
	// slice, we don't need to go any further.
	if aIsSlice || bIsSlice {
		return false
	}
	// TODO
	return equalRecords(a, b, opts)
}
