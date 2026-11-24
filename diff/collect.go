package diff

import (
	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/jaypipes/delta/types"
)

// CollectAt recursively adds any detected differences in the supplied values
// to the Delta, adding a nested field path for each level of nested map or
// struct values examined.
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
		eq, ok := cmp.EqualComparable(a, b, opts...)
		if !ok {
			return cmp.InvalidTypeComparison(a, b)
		} else if eq {
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

	if aIsScalar || bIsScalar {
		// We already handled the case when both subjects are scalars, so if
		// either is a scalar, that represents a case where the field's
		// underlying data type changed...
		d.Push(
			Difference{
				Kind:      DifferenceType,
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
		// We already handled the case when both subjects are slices, so if
		// either is a slice, that represents a case where the field's
		// underlying data type changed...
		d.Push(
			Difference{
				Kind:      DifferenceType,
				FieldPath: fp,
				From:      a,
				To:        b,
			},
		)
		return nil
	}

	return d.collectRecordDifferences(fp, a, b, o)
}
