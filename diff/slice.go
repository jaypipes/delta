package diff

import (
	"slices"

	"github.com/samber/lo"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff/fieldpath"
	"github.com/jaypipes/delta/types"
	"github.com/jaypipes/delta/types/cast"
	"github.com/jaypipes/delta/types/convert"
)

// collectSliceDifferences collects any detected differences in the supplied
// slice things to the Delta.
func (d *Delta) collectSliceDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	aIsSliceScalar := types.IsSliceScalar(a)
	bIsSliceScalar := types.IsSliceScalar(b)

	if aIsSliceScalar && bIsSliceScalar {
		return d.collectSliceScalarDifferences(fp, a, b, opts)
	}

	// You cannot compare a slice of scalar to a slice of non-scalar, so if
	// only one argument is a slice of scalar, we mark that there is a
	// difference in field type.

	if aIsSliceScalar || bIsSliceScalar {
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

	// TODO)
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
		return d.collectSlicePrimitiveDifferences(fp, a, b, opts)
	}

	// One or both of the slice of scalars is a slice of primitive type
	// definitions, so if we are allowed, perform the resolution to the
	// underlying primitive type and then collect the differences.

	if opts.NoResolveTypeDef {
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

	// One or both of the supplied slices of scalars are slices of primitive
	// type defs. Resolve to the slice of underlying primitive and compare
	// those.

	ap := a
	aIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(a)
	if aIsSlicePrimTypeDef {
		ap = types.SliceUnderlyingPrimitive(a)
	}

	bp := b
	bIsSlicePrimTypeDef := types.IsSlicePrimitiveTypeDef(b)
	if bIsSlicePrimTypeDef {
		bp = types.SliceUnderlyingPrimitive(b)
	}
	if ap == nil || bp == nil {
		// There is no easy way to differentiate between a string and a []byte.
		// A string is *not* a slice of primitive but a []byte *is* a slice of
		// primitive.
		//
		// So, if a is []string and b is MyByteArrayType, we fall into this
		// block because IsSlicePrimitive returns true for both []string and
		// MyByteArrayType but SliceUnderlyingPrimitive returns nil for
		// MyByteArrayType and []string for []string arguments.
		d.Push(
			Difference{
				Kind:      DifferenceType,
				FieldPath: fp,
				From:      a,
				To:        b,
			},
		)
	}
	return d.collectSlicePrimitiveDifferences(fp, ap, bp, opts)
}

// collectSlicePrimitiveDifferences collects any detected differences in the
// supplied slice of scalar primitives to the Delta.
func (d *Delta) collectSlicePrimitiveDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	switch {
	case opts.NoCast:
		return d.collectSlicePrimitiveDifferencesNoCast(fp, a, b)
	case opts.NoConvert:
		return d.collectSlicePrimitiveDifferencesCast(fp, a, b)
	default:
		return d.collectSlicePrimitiveDifferencesConvert(fp, a, b, opts)
	}
}

// collectSlicePrimitiveDifferencesNoCast collects any detected
// differences in the supplied slice of scalar primitives to the Delta without
// using any type casting.
func (d *Delta) collectSlicePrimitiveDifferencesNoCast(
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
		}
	default:
		goto addTypeDiff
	}

addTypeDiff:
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

// collectSlicePrimitiveDifferencesCast collects any detected differences in
// the supplied slice of scalar primitives to the Delta using type casting to
// compare element equality.
func (d *Delta) collectSlicePrimitiveDifferencesCast(
	fp fieldpath.FieldPath,
	a any,
	b any,
) error {
	dk := DifferenceModify
	sd := SliceDifference{}
	switch a := a.(type) {
	case []int, []int8, []int16, []int32, []int64:
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			ai, _ := cast.ToSliceInt64(a)
			bi, _ := cast.ToSliceInt64(b)
			if len(ai) != len(bi) {
				aNotInB, bNotInA := lo.Difference(ai, bi)
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
			if slices.Equal(ai, bi) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(ai, bi)
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
			goto addTypeDiff
		}
	case []uint, []uint8, []uint16, []uint32, []uint64:
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			au, _ := cast.ToSliceUint64(a)
			bu, _ := cast.ToSliceUint64(b)
			if len(au) != len(bu) {
				aNotInB, bNotInA := lo.Difference(au, bu)
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
			if slices.Equal(au, bu) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(au, bu)
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
		}
	default:
		goto addTypeDiff
	}

addTypeDiff:
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

// collectSliceScalarPrimitiveDifferencesConvert collects any detected
// differences in the supplied slice of scalar primitives to the Delta using
// type conversion to compare element equality.
func (d *Delta) collectSlicePrimitiveDifferencesConvert(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	dk := DifferenceModify
	sd := SliceDifference{}
	switch a := a.(type) {
	case []int, []int8, []int16, []int32, []int64:
		// Never an error for Integer types.
		ai, _ := cast.ToSliceInt64(a)
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			bi, _ := cast.ToSliceInt64(b)
			if len(ai) != len(bi) {
				aNotInB, bNotInA := lo.Difference(ai, bi)
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
			if slices.Equal(ai, bi) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(ai, bi)
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
			goto addTypeDiff
		}
	case []uint, []uint8, []uint16, []uint32, []uint64:
		// Never an error for Integer types.
		au, _ := cast.ToSliceUint64(a)
		switch b := b.(type) {
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64:
			// Never an error for Integer types.
			bu, _ := cast.ToSliceUint64(b)
			if len(au) != len(bu) {
				aNotInB, bNotInA := lo.Difference(au, bu)
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
			if slices.Equal(au, bu) {
				return nil
			}
			aNotInB, bNotInA := lo.Difference(au, bu)
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
			goto addTypeDiff
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
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64,
			[]bool:
			// Never an error converting slices of integer or bool to slices of
			// string.
			bs, _ := convert.ToSliceString(b)
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
			goto addTypeDiff
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
		case []int, []int8, []int16, []int32, []int64,
			[]uint, []uint8, []uint16, []uint32, []uint64,
			[]bool:
			// Never an error converting slices of integer or bool to slices of
			// string.
			bs, _ := convert.ToSliceString(b)
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
		default:
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
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
			goto addTypeDiff
		}
	default:
		goto addTypeDiff
	}

addTypeDiff:
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
