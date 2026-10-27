package fieldpath_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jaypipes/delta/fieldpath"
)

func TestFromAny(t *testing.T) {
	cases := []struct {
		name    string
		subject any
		expErr  error
		exp     fieldpath.FieldPath
	}{
		{
			"string constructor",
			"a.b.c",
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			"[]string constructor",
			[]string{"a", "b", "c"},
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			"FieldPath constructor",
			fieldpath.FieldPath{"a", "b", "c"},
			nil,
			fieldpath.FieldPath{"a", "b", "c"},
		},
		{
			"invalid subject type",
			42,
			fieldpath.ErrInvalid,
			nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			got, err := fieldpath.FromAny(c.subject)
			if c.expErr != nil {
				require.ErrorIs(err, c.expErr)
			} else {
				require.Nil(err)
				require.Equal(c.exp, got)
			}
		})
	}
}

func TestPrefixed(t *testing.T) {
	fp := fieldpath.FieldPath{"generation"}
	cases := []struct {
		name     string
		prefixes []string
		exp      fieldpath.FieldPath
	}{
		{
			"no prefix",
			[]string{},
			fieldpath.FieldPath{"generation"},
		},
		{
			"single prefix",
			[]string{"spec"},
			fieldpath.FieldPath{"spec", "generation"},
		},
		{
			"two prefixes",
			[]string{"object", "spec"},
			fieldpath.FieldPath{"object", "spec", "generation"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			got := fieldpath.Prefixed(fp, c.prefixes...)
			require.Equal(c.exp, got)
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		name     string
		subject  fieldpath.FieldPath
		contains fieldpath.FieldPath
		exp      bool
	}{
		{
			"a.b.c is in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a", "b", "c"},
			true,
		},
		{
			"a.b is in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a", "b"},
			true,
		},
		{
			"a is in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"a"},
			true,
		},
		{
			"b is NOT in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"b"},
			false,
		},
		{
			"c is NOT in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"c"},
			false,
		},
		{
			"b.c is NOT in a.b.c",
			fieldpath.FieldPath{"a", "b", "c"},
			fieldpath.FieldPath{"b", "c"},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.Contains(c.contains))
		})
	}
}
