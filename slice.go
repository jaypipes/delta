package delta

import (
	"reflect"

	"github.com/samber/lo"
)

// IsSlice returns true if the supplied subject is a slice of something.
func IsSlice(subject any) bool {
	return reflect.TypeOf(subject).Kind() == reflect.Slice
}

// ToAnySlice returns a []any for the supplied scalar slice-valued thing.
//
// If the supplied argument is not a slice type, returns nil.
func ToAnySlice(subject any) []any {
	if !IsSlice(subject) {
		return nil
	}
	switch subject := subject.(type) {
	case []any:
		return subject
	case []rune:
		return lo.ToAnySlice(subject)
	case []string:
		return lo.ToAnySlice(subject)
	case []int:
		return lo.ToAnySlice(subject)
	case []int8:
		return lo.ToAnySlice(subject)
	case []int16:
		return lo.ToAnySlice(subject)
	case []int64:
		return lo.ToAnySlice(subject)
	case []uint:
		return lo.ToAnySlice(subject)
	case []uint8:
		return lo.ToAnySlice(subject)
	case []uint16:
		return lo.ToAnySlice(subject)
	case []uint32:
		return lo.ToAnySlice(subject)
	case []uint64:
		return lo.ToAnySlice(subject)
	case []float32:
		return lo.ToAnySlice(subject)
	case []float64:
		return lo.ToAnySlice(subject)
	case []bool:
		return lo.ToAnySlice(subject)
	}
	return nil
}

// SliceDifferenceResult contains the elements in A that are not in B, the
// elements in B that are not in A, and a boolean representing whether A and B
// are in the same sort order.
type SliceDifferenceResult struct {
	ANotInB       []any
	BNotInA       []any
	SameSortOrder bool
}

// Equal returns true if there were no differences between A and B, including
// sort order.
func (r SliceDifferenceResult) Equal() bool {
	return r.SameSortOrder && len(r.ANotInB) == 0 && len(r.BNotInA) == 0
}

// ElementsEqual returns true if there were no differences between the elements
// in A and B.
func (r SliceDifferenceResult) ElementsEqual() bool {
	return len(r.ANotInB) == 0 && len(r.BNotInA) == 0
}

// SliceDifference returns the difference between the supplied slices.
// equality and sort order.
func SliceDifference(
	a []any,
	b []any,
) SliceDifferenceResult {
	r := SliceDifferenceResult{}
	r.ANotInB, r.BNotInA = lo.Difference(a, b)
	r.SameSortOrder = SliceEqualOrder(a, b)
	return r
}

// SliceEqualOrder returns true if the supplied slices are in the same sort
// order.
//
// Important: naive equivalence is used to determine whether elements at the
// same index in compared slices are the same. So, SliceEqualOrder will return
// true for the following supplied slice arguments:
//
// A: [1, "2"]
// B: ["1", 2]
func SliceEqualOrder(
	a []any,
	b []any,
) bool {
	if len(a) != len(b) {
		return false
	}
	for x, av := range a {
		bv := b[x]
		if !ScalarEquivalent(av, bv) {
			return false
		}
	}
	return true
}
