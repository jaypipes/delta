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
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return compare(a, b, o)
}

// Compare returns 0 if the value of a and b are equal per the supplied
// Options, -1 if the value of a is less than b and 1 if the value of a is
// greater than b.
func compare(a any, b any, opts Options) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	aIsScalar := types.IsScalar(a)
	bIsScalar := types.IsScalar(b)

	aIsSlice := types.IsSlice(a)
	bIsSlice := types.IsSlice(b)

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
		return compareScalars(a, b, opts)
	case aIsScalar && bIsBATD:
		return compareScalars(a, b, opts)
	case bIsScalar && aIsBATD:
		return compareScalars(a, b, opts)
	case aIsSlice && bIsSlice:
		return compareSlices(a, b, opts)
	case aIsSlice && bIsBATD:
		return compareSlices(a, b, opts)
	case bIsSlice && aIsBATD:
		return compareSlices(a, b, opts)
	case aIsSlice || bIsSlice:
		// You cannot compare a slice to a non slice, so if only one argument
		// is a slice, we don't need to go any further.
		return -1
	}

	// If we get here, both a and b must be Record types.
	// TODO
	return compareRecords(a, b, opts)
}
