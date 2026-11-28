package diff

import (
	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff/fieldpath"
)

// Diff returns a Delta representing the set of differences between two
// supplied things `a` and `b`.
//
// If `a` implements the `Differ` interface, then `Diff` will call
// `a.Diff(b)`.
//
// If `a` does not implement the `Differ` interface and the type of `a` and `b`
// is not comparable, ErrTypeComparison will be returned.
func Diff(a any, b any, opts ...cmp.Option) (Delta, error) {
	differ, ok := a.(Differ)
	if ok {
		return differ.Diff(b, opts...)
	}
	d := Delta{}
	err := d.CollectAt(fieldpath.FromString(""), a, b, opts...)
	if err != nil {
		return nil, err
	}
	return d, nil
}
