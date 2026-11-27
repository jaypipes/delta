package diff_test

import (
	"testing"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestDiff_SliceScalar(t *testing.T) {
	cases := []struct {
		name string
		a    any
		b    any
		opts []cmp.Option
		exp  diff.Delta
	}{
		{
			name: "empty slices",
			a:    []string{},
			b:    []string{},
			exp:  diff.Delta{},
		},
		{
			name: "primitives same values different type",
			a:    []string{"a"},
			b:    []any{"a"},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceType,
					From:      []string{"a"},
					To:        []any{"a"},
				},
			},
		},
		{
			name: "primitives different values different type NoCast",
			a:    []string{"a"},
			b:    []any{"b"},
			opts: []cmp.Option{cmp.NoCast},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceType,
					From:      []string{"a"},
					To:        []any{"b"},
				},
			},
		},
		{
			name: "primitives same type and value",
			a:    []string{"a"},
			b:    []string{"a"},
			exp:  diff.Delta{},
		},
		{
			name: "primitives same type different values",
			a:    []string{"a"},
			b:    []string{"b"},
			exp: diff.Delta{
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
		},
		{
			name: "slices same elements different order",
			a:    []string{"a", "b"},
			b:    []string{"b", "a"},
			exp: diff.Delta{
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
		},
		{
			name: "[]string and [][]byte considered slices",
			a:    []string{"a"},
			b:    [][]byte{[]byte("b")},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      []string{"a"},
					To:        [][]byte{[]byte("b")},
					Slice: &diff.SliceDifference{
						Sort:    false,
						ANotInB: []any{"a"},
						BNotInA: []any{"b"},
					},
				},
			},
		},
		{
			name: "[]string and [][]byte considered slices, diff order only",
			a:    []string{"a", "b"},
			b:    [][]byte{[]byte("b"), []byte("a")},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceSort,
					From:      []string{"a", "b"},
					To:        [][]byte{[]byte("b"), []byte("a")},
					Slice: &diff.SliceDifference{
						Sort: true,
					},
				},
			},
		},
		{
			name: "[]string and []byte are NOT comparable when using NoCast",
			a:    []string{"a"},
			b:    []byte("a"),
			opts: []cmp.Option{cmp.NoCast},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceType,
					From:      []string{"a"},
					To:        []byte("a"),
				},
			},
		},
		{
			name: "[]string and []byte ARE comparable with no options",
			a:    []string{"a"},
			b:    []byte("a"),
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      []string{"a"},
					To:        []byte("a"),
					Slice: &diff.SliceDifference{
						Sort:    false,
						ANotInB: []any{"a"},
						BNotInA: []any{"97"},
					},
				},
			},
		},
	}
	require := require.New(t)
	for _, c := range cases {
		d, err := diff.Diff(c.a, c.b, c.opts...)
		require.Nil(err)
		require.NotNil(d)
		require.Equal(c.exp, d)
	}
}
