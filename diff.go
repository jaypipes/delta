package delta

import (
	"reflect"

	"github.com/jaypipes/delta/fieldpath"
)

// Differ allows custom types that know how to determine differences between
// two things to override the more generic comparison functionality in
// Delta.CollectAt().
type Differ interface {
	// Diff should return a Delta representing the set of differences between
	// the thing and the supplied parameter. Should return ErrTypeComparison if
	// the supplied parameter cannot be compared to the thing implementing the
	// Differ interface.
	Diff(any) (Delta, error)
}

// Diff returns a Delta representing the set of differences between two
// supplied things `a` and `b`.
//
// If `a` implements the `Differ` interface, then `Diff` will call
// `a.Diff(b)`.
//
// If `a` does not implement the `Differ` interface and the type of `a` and `b`
// is not comparable, ErrTypeComparison will be returned.
func Diff(a, b any) (Delta, error) {
	differ, ok := a.(Differ)
	if ok {
		return differ.Diff(b)
	}
	d := Delta{}
	err := d.CollectAt(fieldpath.FromString(""), a, b)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// areComparable returns true if the types of the supplied subjects are naively
// comparable.
func areComparable(a, b any) bool {
	if a == nil || b == nil {
		// comparing anything to nil is possible...
		return true
	}
	if IsScalar(a) {
		return IsScalar(b)
	}
	if IsScalar(b) {
		return false
	}
	if IsSlice(a) && IsSlice(b) {
		return areComparable(
			reflect.TypeOf(a).Elem(),
			reflect.TypeOf(b).Elem(),
		)
	}
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
