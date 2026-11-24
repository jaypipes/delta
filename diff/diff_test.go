package diff_test

import (
	"fmt"
	"testing"

	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestDiff_Scalar(t *testing.T) {
	cases := []struct {
		name string
		a    any
		b    any
		exp  diff.Delta
	}{
		{
			"nil to scalar adds",
			nil,
			1,
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceAdd,
					From:      nil,
					To:        1,
				},
			},
		},
		{
			"scalar to nil deletes",
			1,
			nil,
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceRemove,
					From:      1,
					To:        nil,
				},
			},
		},
		{
			"scalar to different valued scalar same type modifies",
			1,
			2,
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      1,
					To:        2,
				},
			},
		},
		{
			"scalar to different valued scalar different type modifies",
			1,
			"2",
			diff.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Kind:      diff.DifferenceModify,
					From:      1,
					To:        "2",
				},
			},
		},
		{
			"scalar to same valued scalar same type no differences",
			1,
			1,
			diff.Delta{},
		},
		{
			"scalar to same valued scalar different type no differences",
			1,
			"1",
			diff.Delta{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := diff.Diff(c.a, c.b)
			require.Nil(err)
			require.NotNil(d)
			require.Equal(c.exp, d)
		})
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
