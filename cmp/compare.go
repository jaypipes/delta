package cmp

import "github.com/jaypipes/delta/types"

// Compare returns 0 if the value of a and b are equal per the supplied
// Options, -1 if the value of a is less than b and 1 if the value of a is
// greater than b.
//
// If a and b are not comparable per the supplied Options, returns -1. Callers
// should use [Comparable] to check that two subjects are comparable before
// calling Compare.
func Compare(a any, b any, opts ...Option) int {
	if !Comparable(a, b, opts...) {
		return -1
	}

	ac := types.CategoryOf(a)
	bc := types.CategoryOf(b)

	// We compare nils specially. If both arguments are nil, we return 0. If a
	// is nil we return -1. If b is nil, we return 1.
	switch {
	case ac == types.Nil && bc == types.Nil:
		return 0
	case ac == types.Nil:
		return -1
	case bc == types.Nil:
		return 1
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// Note that if we got here, both arguments must be the same Category,
	// which is why we can shorten below to the switch case on a's category.
	switch ac {
	case types.Scalar:
		return compareScalars(a, b, o)
	case types.Slice:
		return compareSlices(a, b, o)
	default:
		return compareRecords(a, b, o)
	}
}
