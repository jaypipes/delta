package fieldpath_test

import (
	"testing"

	"github.com/jaypipes/delta/diff/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestFromAny(t *testing.T) {
	cases := []struct {
		subject any
		expErr  error
		exp     fieldpath.FieldPath
	}{
		{
			"a.b.c",
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			[]string{"a", "b", "c"},
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			42,
			fieldpath.ErrInvalid,
			nil,
		},
	}
	require := require.New(t)
	for _, c := range cases {
		got, err := fieldpath.FromAny(c.subject)
		if c.expErr != nil {
			require.ErrorIs(err, c.expErr)
		} else {
			require.Nil(err)
			require.Equal(c.exp, got)
		}
	}
}

func TestPrefixed(t *testing.T) {
	fp := fieldpath.FieldPath{"generation"}
	cases := []struct {
		prefixes []string
		exp      fieldpath.FieldPath
	}{
		{
			[]string{},
			fieldpath.FieldPath{"generation"},
		},
		{
			[]string{"spec"},
			fieldpath.FieldPath{"spec", "generation"},
		},
		{
			[]string{"object", "spec"},
			fieldpath.FieldPath{"object", "spec", "generation"},
		},
	}
	require := require.New(t)
	for _, c := range cases {
		got := fieldpath.Prefixed(fp, c.prefixes...)
		require.Equal(c.exp, got)
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		subject  fieldpath.FieldPath
		contains fieldpath.FieldPath
		exp      bool
	}{
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a", "b", "c"},
			true,
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a", "b"},
			true,
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a"},
			true,
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"b"},
			false,
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"c"},
			false,
		},
		{
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"b", "c"},
			false,
		},
	}
	require := require.New(t)
	for _, c := range cases {
		require.Equal(c.exp, c.subject.Contains(c.contains))
	}
}
