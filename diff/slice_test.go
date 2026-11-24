package diff_test

import (
	"testing"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestDiff_SliceScalar_NoCast(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    diff.Delta
		expErr string
	}{
		{
			"empty slices",
			[]string{},
			[]string{},
			diff.Delta{},
			"",
		},
		{
			"primitives same values different type",
			[]string{"a"},
			[]any{"a"},
			diff.Delta{},
			"invalid type comparison",
		},
		{
			"primitives different values different type",
			[]string{"a"},
			[]any{"b"},
			diff.Delta{},
			"invalid type comparison",
		},
		{
			"primitives same type and value",
			[]string{"a"},
			[]string{"a"},
			diff.Delta{},
			"",
		},
		{
			"primitives same type different values",
			[]string{"a"},
			[]string{"b"},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      []string{"a"},
					To:        []string{"b"},
					Slice: &diff.SliceDifference{
						Sort:    false,
						ANotInB: []any{"a"},
						BNotInA: []any{"b"},
					},
				},
			},
			"",
		},
		{
			"slices same elements different order",
			[]string{"a", "b"},
			[]string{"b", "a"},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceSort,
					From:      []string{"a", "b"},
					To:        []string{"b", "a"},
					Slice: &diff.SliceDifference{
						Sort: true,
					},
				},
			},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := diff.Diff(c.a, c.b, cmp.NoCast)
			if c.expErr != "" {
				require.ErrorContains(err, c.expErr)
			} else {
				require.Nil(err)
				require.NotNil(d)
				require.Equal(c.exp, d)
			}
		})
	}
}
