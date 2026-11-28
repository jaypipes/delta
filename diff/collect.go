package diff

import (
	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff/fieldpath"
	"github.com/jaypipes/delta/types"
)

// CollectAt recursively adds any detected differences in the supplied values
// to the Delta, adding a nested field path for each level of nested map or
// struct values examined.
//
// The supplied subject values represent the contents of a field at the
// supplied field path. Therefore, if a is nil and b is not nil, that indicates
// the kind of difference is additive. If a is not nil and b is nil, the kind
// of difference is subtractive. If both a and b is not nil, any difference
// kind is mutative.
func (d *Delta) CollectAt(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts ...cmp.Option,
) error {
	differ, ok := a.(Differ)
	if ok {
		other, err := differ.Diff(b)
		if err != nil {
			return err
		}
		d.ExtendAt(fp, other)
		return nil
	}

	// First handle `nil`.
	//
	// NOTE(jaypipes): `nil` is not the same an interface referencing nil.
	//
	// If both a and b are `nil` there is no difference to collect.
	//
	// If a is `nil` and b is not `nil`, this indicates an additiive
	// difference.
	//
	// If a is not `nil` and b is `nil`, this indicates a subtractive
	// difference.
	if a == nil {
		if b == nil {
			// No difference... so just return.
			return nil
		}
		// the `b` thing has a value. the `a` thing does not, so this is an
		// additive difference.
		d.Push(
			Difference{
				Kind:      DifferenceAdd,
				FieldPath: fp,
				From:      nil,
				To:        b,
			},
		)
		return nil
	} else {
		if b == nil {
			// the `a` thing has a value. the `b` thing does not, so this is a
			// removing difference.
			d.Push(
				Difference{
					Kind:      DifferenceRemove,
					FieldPath: fp,
					From:      a,
					To:        nil,
				},
			)
			return nil
		}
	}

	aIsScalar := types.IsScalar(a)
	bIsScalar := types.IsScalar(b)

	if aIsScalar && bIsScalar {
		if cmp.Equal(a, b, opts...) {
			return nil
		}
		d.Push(
			Difference{
				Kind:      DifferenceModify,
				FieldPath: fp,
				From:      a,
				To:        b,
			},
		)
		return nil
	}

	o := cmp.Options{}
	for _, opt := range opts {
		opt(&o)
	}

	aIsSlice := types.IsSlice(a)
	bIsSlice := types.IsSlice(b)

	if aIsSlice && bIsSlice {
		return d.collectSliceDifferences(fp, a, b, o)
	}

	if aIsSlice || bIsSlice {
		// We cannot compare a slice to a non-slice
		return cmp.InvalidTypeComparison(a, b)
	}

	return d.collectRecordDifferences(fp, a, b, o)
}
