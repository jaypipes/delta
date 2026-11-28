package diff_test

import (
	"fmt"
	"testing"

	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/diff/fieldpath"
	"github.com/stretchr/testify/require"
)

func TestDiff_MapStringAny(t *testing.T) {
	cases := []struct {
		name string
		a    any
		b    any
		exp  diff.Delta
	}{
		{
			"empty maps",
			map[string]any{},
			map[string]any{},
			diff.Delta{},
		},
		{
			"map add all keys",
			map[string]any{},
			map[string]any{
				"x": "valuex",
				"y": map[string]any{
					"foo": "valuefoo",
				},
			},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Kind:      diff.DifferenceAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Kind:      diff.DifferenceAdd,
					From:      nil,
					To: map[string]any{
						"foo": "valuefoo",
					},
				},
			},
		},
		{
			"map remove all keys",
			map[string]any{
				"x": "valuex",
				"y": map[string]any{
					"foo": "valuefoo",
				},
			},
			map[string]any{},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Kind:      diff.DifferenceRemove,
					From:      "valuex",
					To:        nil,
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Kind:      diff.DifferenceRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
			},
		},
		{
			"map modify all keys",
			map[string]any{
				"x": "valuex1",
				"y": map[string]any{
					"foo": "valuefoo1",
				},
			},
			map[string]any{
				"x": "valuex2",
				"y": map[string]any{
					"foo": "valuefoo2",
				},
			},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Kind:      diff.DifferenceModify,
					From:      "valuex1",
					To:        "valuex2",
				},
				{
					FieldPath: fieldpath.FromString("y.foo"),
					Kind:      diff.DifferenceModify,
					From:      "valuefoo1",
					To:        "valuefoo2",
				},
			},
		},
		{
			"map add remove and modify keys",
			map[string]any{
				"y": map[string]any{
					"foo": "valuefoo",
				},
				"bar": "bar",
			},
			map[string]any{
				"x":   "valuex",
				"bar": "baz",
			},
			diff.Delta{
				{
					FieldPath: fieldpath.FromString("y"),
					Kind:      diff.DifferenceRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
				{
					FieldPath: fieldpath.FromString("x"),
					Kind:      diff.DifferenceAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("bar"),
					Kind:      diff.DifferenceModify,
					From:      "bar",
					To:        "baz",
				},
			},
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

// Two maps representing two Kubernetes Deployments, `deploy1` and `deploy2`.
// Use `diff.Diff()` to get the difference between the two Deployments:
func ExampleDiff_nested_maps() {
	var deploy1 = map[string]any{
		"spec": map[string]any{
			"replicas": 1,
		},
		"status": map[string]any{
			"readyReplicas":      1,
			"observedGeneration": 1,
		},
	}
	var deploy2 = map[string]any{
		"spec": map[string]any{
			"replicas": 3,
		},
		"status": map[string]any{
			"readyReplicas": 2,
		},
	}

	d, err := diff.Diff(deploy1, deploy2)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
	// Output:
	// spec.replicas MODIFY 1 -> 3
	// status.observedGeneration REMOVE previous value: 1
	// status.readyReplicas MODIFY 1 -> 2
}
