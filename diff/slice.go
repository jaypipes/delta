package diff

import (
	"slices"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/jaypipes/delta/types"
	"github.com/samber/lo"
)

// collectSliceDifferences collects any detected differences in the supplied
// slice things to the Delta.
func (d *Delta) collectSliceDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	aec, ok := types.CategoryOfElem(a)
	if !ok {
		return cmp.InvalidTypeComparison(a, b)
	}
	bec, ok := types.CategoryOfElem(b)
	if !ok {
		return cmp.InvalidTypeComparison(a, b)
	}

	if aec != bec {
		return cmp.InvalidTypeComparison(a, b)
	}

	switch aec {
	case types.Scalar:
		return d.collectSliceScalarDifferences(fp, a, b, opts)
	}
	return nil
}

// collectSliceScalarDifferences collects any detected differences in the
// supplied slice of scalars to the Delta.
func (d *Delta) collectSliceScalarDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	aIsSlicePrim := types.IsSlicePrimitive(a)
	bIsSlicePrim := types.IsSlicePrimitive(b)

	if aIsSlicePrim && bIsSlicePrim {
		return d.collectSliceScalarPrimitiveDifferences(fp, a, b, opts)
	}

	// One or both of the slice of scalars is a slice of primitive type
	// definitions, so if we are allowed, perform the resolution to the
	// underlying primitive type and then collect the differences.

	if opts.NoResolveTypeDef {
		return cmp.InvalidTypeComparison(a, b)
	}

	// TODO
	return nil
}

// collectSliceScalarPrimitiveDifferences collects any detected differences in
// the supplied slice of scalar primitives to the Delta.
func (d *Delta) collectSliceScalarPrimitiveDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	switch {
	case opts.NoCast:
		return d.collectSliceScalarPrimitiveDifferencesNoCast(fp, a, b)
	case opts.NoConvert:
		return d.collectSliceScalarPrimitiveDifferencesCast(fp, a, b)
	default:
		return d.collectSliceScalarPrimitiveDifferencesConvert(fp, a, b, opts)
	}
}

// collectSliceScalarPrimitiveDifferencesNoCast collects any detected
// differences in the supplied slice of scalar primitives to the Delta without
// using any type casting.
func (d *Delta) collectSliceScalarPrimitiveDifferencesNoCast(
	fp fieldpath.FieldPath,
	a any,
	b any,
) error {
	dk := DifferenceModify
	sd := SliceDifference{}
	switch a := a.(type) {
	case []int:
		switch b := b.(type) {
		case []int:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      dk,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []int8:
		switch b := b.(type) {
		case []int8:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []int16:
		switch b := b.(type) {
		case []int16:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []int32:
		switch b := b.(type) {
		// Note that rune is a type alias for int32
		case []int32:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []int64:
		switch b := b.(type) {
		case []int64:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []uint:
		switch b := b.(type) {
		case []uint:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []uint8:
		switch b := b.(type) {
		// Note that byte is a type alias for uint8
		case []uint8:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []uint16:
		switch b := b.(type) {
		case []uint16:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []uint32:
		switch b := b.(type) {
		case []uint32:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []uint64:
		switch b := b.(type) {
		case []uint64:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []string:
		switch b := b.(type) {
		case []string:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		// string is just an immutable []byte
		case [][]byte:
			bs := make([]string, len(b))
			for x, el := range b {
				bs[x] = string(el)
			}
			if len(a) != len(bs) {
				aNotInB, bNotInA := lo.Difference(a, bs)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, bs) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, bs)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case [][]byte:
		as := make([]string, len(a))
		for x, el := range a {
			as[x] = string(el)
		}
		switch b := b.(type) {
		case [][]byte:
			bs := make([]string, len(b))
			for x := range len(a) {
				bs[x] = string(b[x])
			}
			if len(as) != len(bs) {
				aNotInB, bNotInA := lo.Difference(as, bs)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(as, bs) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(as, bs)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		// string is just an immutable []byte
		case []string:
			if len(as) != len(b) {
				aNotInB, bNotInA := lo.Difference(as, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(as, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(as, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []float32:
		switch b := b.(type) {
		case []float32:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []float64:
		switch b := b.(type) {
		case []float64:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			if slices.Equal(a, b) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// lo.Difference returns empty differences even if
				// slices.Equal() returns false because lo.Difference does not
				// account for sort order. So, if we get here, we know the only
				// difference betewen the two slices was the sort order.
				dk = DifferenceSort
				sd.Sort = true
			} else {
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
			}
			d.Push(
				Difference{
					Kind:      dk,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	case []bool:
		switch b := b.(type) {
		case []bool:
			if len(a) != len(b) {
				aNotInB, bNotInA := lo.Difference(a, b)
				sd.ANotInB = lo.ToAnySlice(aNotInB)
				sd.BNotInA = lo.ToAnySlice(bNotInA)
				d.Push(
					Difference{
						Kind:      DifferenceModify,
						FieldPath: fp,
						From:      a,
						To:        b,
						Slice:     &sd,
					},
				)
				return nil
			}
			aNotInB, bNotInA := lo.Difference(a, b)
			if len(aNotInB) == 0 && len(bNotInA) == 0 {
				// Same  sort order, same elements, no difference to collect.
				return nil
			}
			// NOTE(jaypipes): You cannot sort slices of bool.
			sd.ANotInB = lo.ToAnySlice(aNotInB)
			sd.BNotInA = lo.ToAnySlice(bNotInA)
			d.Push(
				Difference{
					Kind:      DifferenceModify,
					FieldPath: fp,
					From:      a,
					To:        b,
					Slice:     &sd,
				},
			)
			return nil
		default:
			// No cast means comparing to any other slice element type is
			// invalid.
			return cmp.InvalidTypeComparison(a, b)
		}
	default:
		return cmp.InvalidTypeComparison(a, b)
	}
}

// collectSliceScalarPrimitiveDifferencesCast collects any detected differences
// in the supplied slice of scalar primitives to the Delta using type casting
// to compare element equality.
func (d *Delta) collectSliceScalarPrimitiveDifferencesCast(
	fp fieldpath.FieldPath,
	a any,
	b any,
) error {
	// TODO
	return cmp.InvalidTypeComparison(a, b)
}

// collectSliceScalarPrimitiveDifferencesConvert collects any detected differences
// in the supplied slice of scalar primitives to the Delta using type conversion
// to compare element equality.
func (d *Delta) collectSliceScalarPrimitiveDifferencesConvert(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	// TODO
	return cmp.InvalidTypeComparison(a, b)
}
