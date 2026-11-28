package diff

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/jaypipes/delta/diff/fieldpath"
)

// DifferenceKind differentiates kinds of Differences
type DifferenceKind int

const (
	// DifferenceAdd represents an addition of some field or item in a list/map
	// field value type.
	DifferenceAdd DifferenceKind = iota
	// DifferenceRemove represents the removal of some field or item in a
	// list/map field value type.
	DifferenceRemove
	// DifferenceModify represents a modification of some field or item in a
	// list/map field value type.
	DifferenceModify
	// DifferenceSort represents a change in the order of elements in a sorted
	// list field.
	DifferenceSort
	// DifferenceType indicates the that underlying type of the fields being
	// compared is different.
	DifferenceType
)

func (t DifferenceKind) String() string {
	switch t {
	case DifferenceAdd:
		return "ADD"
	case DifferenceRemove:
		return "REMOVE"
	case DifferenceModify:
		return "MODIFY"
	case DifferenceSort:
		return "SORT"
	case DifferenceType:
		return "TYPE"
	}
	return ""
}

// Difference represents a single difference between two things.
type Difference struct {
	// FieldPath contains the FieldPath where this difference exists.
	FieldPath fieldpath.FieldPath
	// Kind contains the type of difference.
	Kind DifferenceKind
	// From contains the "A" or "from" object's value at the FieldPath.
	From any
	// To contains the "B" or "to" object's value at the FieldPath.
	To any
	// Slice contains specific differences in a field containing a slice.
	Slice *SliceDifference
}

// String returns a formatted string representation of the Difference.
func (d Difference) String() string {
	fp := d.FieldPath.String()
	if len(fp) > 0 {
		fp += " "
	}
	switch d.Kind {
	case DifferenceAdd:
		return fmt.Sprintf("%sADD new value: %v", fp, d.To)
	case DifferenceRemove:
		return fmt.Sprintf("%sREMOVE previous value: %v", fp, d.From)
	case DifferenceModify:
		return fmt.Sprintf("%sMODIFY %v -> %v", fp, d.From, d.To)
	case DifferenceSort:
		return fmt.Sprintf("%sSORT", fp)
	case DifferenceType:
		return fmt.Sprintf("%sTYPE %T -> %T", fp, d.From, d.To)
	}
	return ""
}

// SliceDifference contains the elements in A that are not in B, the elements
// in B that are not in A, and a boolean representing whether A and B are in
// the same sort order.
type SliceDifference struct {
	// ANotInB contains elements in A that were not found in B.
	ANotInB []any
	// BNotInA contains elements in B that were not found in A.
	BNotInA []any
	// Sort is true if the sort order of the two slices was different.
	Sort bool
}

// Empty returns true if there were no differences between A and B, including
// sort order.
func (d SliceDifference) Empty() bool {
	return !d.Sort && len(d.ANotInB) == 0 && len(d.BNotInA) == 0
}

// ElementsEqual returns true if there were no differences between the elements
// in A and B.
func (d SliceDifference) ElementsEqual() bool {
	return len(d.ANotInB) == 0 && len(d.BNotInA) == 0
}

// NOTE(jaypipes): This code comes directly from sambe/lo's [Difference] but
// extended with a transform predicate.
//
// TODO(jaypipes): Contribute this back upstream.

// differenceSmallThreshold is the per-side length below which Difference uses a
// nested allocation-free scan instead of building two Keyify maps: for tiny
// inputs the map hashing + heap allocation overhead dominates the O(n*m) scan.
const differenceSmallThreshold = 8

// DifferenceBy returns the difference between two collections using the
// supplied transform against slice elements before checking for equality.
// The first value is the collection of elements absent from list2.
// The second value is the collection of elements absent from list1.
// Play: https://go.dev/play/p/pKE-JgzqRpz
func DifferenceBy[T comparable, Slice ~[]T](
	list1, list2 Slice,
	transform func(T, int) T,
) (Slice, Slice) {
	// Below the threshold an allocation-free nested O(n*m) scan is cheaper; above
	// it the map lookups (O(n+m)) win.
	if len(list1) <= differenceSmallThreshold && len(list2) <= differenceSmallThreshold {
		return differenceBySmall(list1, list2, transform)
	}
	return differenceByLarge(list1, list2, transform)
}

func differenceByLarge[T comparable, Slice ~[]T](
	list1, list2 Slice,
	transform func(T, int) T,
) (Slice, Slice) {
	left := make(Slice, 0, len(list1))
	right := make(Slice, 0, len(list2))

	seenLeft := lo.Keyify(lo.Map(list1, transform))
	seenRight := lo.Keyify(lo.Map(list2, transform))

	for i := range list1 {
		if _, ok := seenRight[list1[i]]; !ok {
			left = append(left, list1[i])
		}
	}

	for i := range list2 {
		if _, ok := seenLeft[list2[i]]; !ok {
			right = append(right, list2[i])
		}
	}

	return left, right
}

func differenceBySmall[T comparable, Slice ~[]T](
	list1, list2 Slice,
	transform func(T, int) T,
) (Slice, Slice) {
	left := make(Slice, 0, len(list1))
	right := make(Slice, 0, len(list2))

	// Same == equality as the map path: an element is kept only when no equal
	// element exists in the other list (NaN never matches, mirroring map keys).
	for i := range list1 {
		found := false
		for j := range list2 {
			if transform(list1[i], i) == transform(list2[j], j) {
				found = true
				break
			}
		}
		if !found {
			left = append(left, list1[i])
		}
	}

	for i := range list2 {
		found := false
		for j := range list1 {
			if transform(list2[i], i) == transform(list1[j], j) {
				found = true
				break
			}
		}
		if !found {
			right = append(right, list2[i])
		}
	}

	return left, right
}
