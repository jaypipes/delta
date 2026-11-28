package diff

import (
	"sort"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff/fieldpath"
	"github.com/samber/lo"
)

// collectRecordDifferences recursively adds any detected differences in the
// supplied records to the Delta, adding a nested field path for each level of
// map values examined.
func (d *Delta) collectRecordDifferences(
	fp fieldpath.FieldPath,
	a any,
	b any,
	opts cmp.Options,
) error {
	// Finally, we treat map[string]any specially since this is the type of
	// json.Raw
	switch a := a.(type) {
	case map[string]any:
		switch b := b.(type) {
		case map[string]any:
			return d.collectMapDifferences(fp, a, b, opts)
		}
	}
	return cmp.InvalidTypeComparison(a, b)
}

// collectMapDifferences recursively adds any detected differences in the
// supplied map of string to interface values to the Delta, adding a
// nested field path for each level of map values examined.
func (d *Delta) collectMapDifferences(
	fp fieldpath.FieldPath,
	a map[string]any,
	b map[string]any,
	opts cmp.Options,
) error {
	akeys := lo.Keys(a)
	sort.Strings(akeys)
	bkeys := lo.Keys(b)
	sort.Strings(bkeys)
	keysNotInB, keysNotInA := lo.Difference(akeys, bkeys)
	for _, k := range keysNotInB {
		err := d.CollectAt(
			fieldpath.Prefixed(fieldpath.FromString(k), fp...),
			a[k],
			nil,
		)
		if err != nil {
			return err
		}
	}
	for _, k := range keysNotInA {
		err := d.CollectAt(
			fieldpath.Prefixed(fieldpath.FromString(k), fp...),
			nil,
			b[k],
		)
		if err != nil {
			return err
		}
	}
	keysInBoth := lo.Intersect(akeys, bkeys)
	for _, k := range keysInBoth {
		av := a[k]
		bv := b[k]
		err := d.CollectAt(
			fieldpath.Prefixed(fieldpath.FromString(k), fp...),
			av,
			bv,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
