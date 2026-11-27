package diff_test

import (
	"fmt"
	"testing"

	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestDiff_Scalar(t *testing.T) {
	cases := []struct {
		name string
		a    any
		b    any
		opts []cmp.Option
		exp  diff.Delta
	}{
		{
			name: "scalar to same valued scalar same type no differences",
			a:    1,
			b:    1,
			exp:  diff.Delta{},
		},
		{
			name: "scalar to same valued scalar different type no differences",
			a:    1,
			b:    "1",
			exp:  diff.Delta{},
		},
		{
			name: "nil to scalar adds",
			a:    nil,
			b:    1,
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceAdd,
					From:      nil,
					To:        1,
				},
			},
		},
		{
			name: "scalar to nil deletes",
			a:    1,
			b:    nil,
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceRemove,
					From:      1,
					To:        nil,
				},
			},
		},
		{
			name: "scalar to different valued scalar same type modifies",
			a:    1,
			b:    2,
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      1,
					To:        2,
				},
			},
		},
		{
			name: "scalar to different valued scalar different type modifies",
			a:    1,
			b:    "2",
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      1,
					To:        "2",
				},
			},
		},
		{
			name: "string and []byte considered scalars",
			a:    "foo",
			b:    []byte("bar"),
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      "foo",
					To:        []byte("bar"),
				},
			},
		},
		{
			name: "string and int comparison no difference",
			a:    "42",
			b:    42,
			exp:  diff.Delta{},
		},
		{
			name: "string and int comparison NoCast yields diff",
			a:    "42",
			b:    42,
			opts: []cmp.Option{cmp.NoCast},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      "42",
					To:        42,
				},
			},
		},
		{
			name: "[]string and []int comparison no difference",
			a:    []string{"42"},
			b:    []int{42},
			exp:  diff.Delta{},
		},
		{
			name: "[]string and []int comparison NoCast diff kind type",
			a:    []string{"42"},
			b:    []int{42},
			opts: []cmp.Option{cmp.NoCast},
			exp: diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceType,
					From:      []string{"42"},
					To:        []int{42},
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

func ExampleDiff_simple_scalars() {
	var a = "foo"
	var b = "bar"

	d, err := diff.Diff(a, b)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
	// Output:
	// MODIFY foo -> bar
}

// The default comparison logic is to convert subjects into Go types that can
// be compared for "equivalence".
//
// "1" is equivalent to 1. int64(1) is equivalent to uint(1), etc. We do this
// because the delta library was designed to compare JSON data that may have
// had numbers encoded as strings and where JSON does not differentiate like Go
// does between signed, unsigned and differently-sized integers.
func ExampleDiff_scalar_equivalence() {
	var x = 1
	var y = "1"

	d, err := diff.Diff(x, y)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
	// Output:
	// <empty>
}
