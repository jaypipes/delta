package delta_test

import (
	"testing"

	"github.com/jaypipes/delta"
	"github.com/jaypipes/delta/fieldpath"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestDiff_Scalar(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    delta.Delta
		expErr string
	}{
		{
			"nil to scalar adds",
			nil,
			1,
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To:        1,
				},
			},
			"",
		},
		{
			"scalar to nil deletes",
			1,
			nil,
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeRemove,
					From:      1,
					To:        nil,
				},
			},
			"",
		},
		{
			"scalar to different valued scalar same type modifies",
			1,
			2,
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeModify,
					From:      1,
					To:        2,
				},
			},
			"",
		},
		{
			"scalar to different valued scalar different type modifies",
			1,
			"2",
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeModify,
					From:      1,
					To:        "2",
				},
			},
			"",
		},
		{
			"scalar to same valued scalar same type no differences",
			1,
			1,
			delta.Delta{},
			"",
		},
		{
			"scalar to same valued scalar different type no differences",
			1,
			"1",
			delta.Delta{},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := delta.Diff(c.a, c.b)
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

func TestDiff_MapStringAny(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    delta.Delta
		expErr string
	}{
		{
			"empty maps",
			map[string]any{},
			map[string]any{},
			delta.Delta{},
			"",
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
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To: map[string]any{
						"foo": "valuefoo",
					},
				},
			},
			"",
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
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeRemove,
					From:      "valuex",
					To:        nil,
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
			},
			"",
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
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeModify,
					From:      "valuex1",
					To:        "valuex2",
				},
				{
					FieldPath: fieldpath.FromString("y.foo"),
					Type:      delta.DifferenceTypeModify,
					From:      "valuefoo1",
					To:        "valuefoo2",
				},
			},
			"",
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
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("bar"),
					Type:      delta.DifferenceTypeModify,
					From:      "bar",
					To:        "baz",
				},
			},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := delta.Diff(c.a, c.b)
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

func TestDiff_KubeMapTypes(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    delta.Delta
		expErr string
	}{
		{
			"empty unstructureds",
			unstructured.Unstructured{},
			unstructured.Unstructured{},
			delta.Delta{},
			"",
		},
		{
			"unstructured add all keys",
			unstructured.Unstructured{},
			unstructured.Unstructured{
				Object: map[string]any{
					"x": "valuex",
					"y": map[string]any{
						"foo": "valuefoo",
					},
				},
			},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To: map[string]any{
						"foo": "valuefoo",
					},
				},
			},
			"",
		},
		{
			"unstructured remove all keys",
			unstructured.Unstructured{
				Object: map[string]any{
					"x": "valuex",
					"y": map[string]any{
						"foo": "valuefoo",
					},
				},
			},
			unstructured.Unstructured{},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeRemove,
					From:      "valuex",
					To:        nil,
				},
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
			},
			"",
		},
		{
			"unstructured modify all keys",
			unstructured.Unstructured{
				Object: map[string]any{
					"x": "valuex1",
					"y": map[string]any{
						"foo": "valuefoo1",
					},
				},
			},
			unstructured.Unstructured{
				Object: map[string]any{
					"x": "valuex2",
					"y": map[string]any{
						"foo": "valuefoo2",
					},
				},
			},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeModify,
					From:      "valuex1",
					To:        "valuex2",
				},
				{
					FieldPath: fieldpath.FromString("y.foo"),
					Type:      delta.DifferenceTypeModify,
					From:      "valuefoo1",
					To:        "valuefoo2",
				},
			},
			"",
		},
		{
			"unstructured add remove and modify keys",
			unstructured.Unstructured{
				Object: map[string]any{
					"y": map[string]any{
						"foo": "valuefoo",
					},
					"bar": "bar",
				},
			},
			unstructured.Unstructured{
				Object: map[string]any{
					"x":   "valuex",
					"bar": "baz",
				},
			},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString("y"),
					Type:      delta.DifferenceTypeRemove,
					From: map[string]any{
						"foo": "valuefoo",
					},
					To: nil,
				},
				{
					FieldPath: fieldpath.FromString("x"),
					Type:      delta.DifferenceTypeAdd,
					From:      nil,
					To:        "valuex",
				},
				{
					FieldPath: fieldpath.FromString("bar"),
					Type:      delta.DifferenceTypeModify,
					From:      "bar",
					To:        "baz",
				},
			},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := delta.Diff(c.a, c.b)
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

func TestDiff_Slice(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    delta.Delta
		expErr string
	}{
		{
			"empty slices",
			[]string{},
			[]string{},
			delta.Delta{},
			"",
		},
		{
			"slices same values same type",
			[]string{"a"},
			[]string{"a"},
			delta.Delta{},
			"",
		},
		{
			"slices same values different type",
			[]string{"a"},
			[]any{"a"},
			delta.Delta{},
			"",
		},
		{
			"slices different values same type",
			[]string{"a"},
			[]string{"b"},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeModify,
					From:      []string{"a"},
					To:        []string{"b"},
				},
			},
			"",
		},
		{
			"slices different values different type",
			[]string{"a"},
			[]any{"b"},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeModify,
					From:      []string{"a"},
					To:        []any{"b"},
				},
			},
			"",
		},
		{
			"slices different order same type",
			[]string{"a", "b"},
			[]string{"b", "a"},
			delta.Delta{
				{
					FieldPath: fieldpath.FromString(""),
					Type:      delta.DifferenceTypeOrder,
					From:      []string{"a", "b"},
					To:        []string{"b", "a"},
				},
			},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := delta.Diff(c.a, c.b)
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

func TestDiff_KubeSliceTypes(t *testing.T) {
	cases := []struct {
		name   string
		a      any
		b      any
		exp    delta.Delta
		expErr string
	}{
		{
			"empty unstructured lists",
			unstructured.UnstructuredList{},
			unstructured.UnstructuredList{},
			delta.Delta{},
			"",
		},
		{
			"empty pointers to unstructured lists",
			&unstructured.UnstructuredList{},
			&unstructured.UnstructuredList{},
			delta.Delta{},
			"",
		},
		{
			"unstructured list values same type",
			unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					unstructured.Unstructured{
						Object: map[string]any{
							"a": "a",
						},
					},
				},
			},
			unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					unstructured.Unstructured{
						Object: map[string]any{
							"a": "a",
						},
					},
				},
			},
			delta.Delta{},
			"",
		},
		{
			"pointers to unstructured list values same type",
			&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					unstructured.Unstructured{
						Object: map[string]any{
							"a": "a",
						},
					},
				},
			},
			&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					unstructured.Unstructured{
						Object: map[string]any{
							"a": "a",
						},
					},
				},
			},
			delta.Delta{},
			"",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			d, err := delta.Diff(c.a, c.b)
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
