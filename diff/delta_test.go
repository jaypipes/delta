package diff_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/fieldpath"
)

var (
	fpSpecReplicas        = fieldpath.FieldPath{"spec", "replicas"}
	fpSpecGeneration      = fieldpath.FieldPath{"spec", "generation"}
	fpStatusReadyReplicas = fieldpath.FieldPath{"status", "readyReplicas"}
)

var (
	diffAB_SpecReplicas = diff.Difference{
		FieldPath: fpSpecReplicas,
		Kind:      diff.DifferenceModify,
		From:      1,
		To:        2,
	}
	diffAB_StatusReadyReplicas = diff.Difference{
		FieldPath: fpStatusReadyReplicas,
		Kind:      diff.DifferenceModify,
		From:      0,
		To:        1,
	}
	diffAC_StatusReadyReplicas = diff.Difference{
		FieldPath: fpStatusReadyReplicas,
		Kind:      diff.DifferenceModify,
		From:      0,
		To:        1,
	}
	diffBC_SpecReplicas = diff.Difference{
		FieldPath: fpSpecReplicas,
		Kind:      diff.DifferenceModify,
		From:      2,
		To:        1,
	}
	diffBC_StatusReadyReplicas = diff.Difference{
		FieldPath: fpStatusReadyReplicas,
		Kind:      diff.DifferenceModify,
		From:      1,
		To:        0,
	}
)

var (
	emptyDelta = diff.Delta{}
	// B is different from A in spec.replicas and status.readyReplicas
	deltaAB = diff.Delta{
		diffAB_SpecReplicas,
		diffAB_StatusReadyReplicas,
	}
	// C is different from A only in status.readyReplicas
	deltaAC = diff.Delta{
		diffAC_StatusReadyReplicas,
	}
	// C is different from B in spec.replicas and status.readyReplicas
	deltaBC = diff.Delta{
		diffBC_SpecReplicas,
		diffBC_StatusReadyReplicas,
	}
)

func TestDelta_Different(t *testing.T) {
	cases := []struct {
		name    string
		subject diff.Delta
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			false,
		},
		{
			"non-empty delta",
			deltaAB,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.Different())
		})
	}
}

func TestDelta_ExtendAt(t *testing.T) {
	require := require.New(t)
	abc := diff.Delta{}
	abc.ExtendAt(fieldpath.FromString("ab"), deltaAB)
	abc.ExtendAt(fieldpath.FromString("ac"), deltaAC)
	expected := diff.Delta{
		{
			FieldPath: fieldpath.FromString("ab.spec.replicas"),
			Kind:      diff.DifferenceModify,
			From:      1,
			To:        2,
		},
		{
			FieldPath: fieldpath.FromString("ab.status.readyReplicas"),
			Kind:      diff.DifferenceModify,
			From:      0,
			To:        1,
		},
		{
			FieldPath: fieldpath.FromString("ac.status.readyReplicas"),
			Kind:      diff.DifferenceModify,
			From:      0,
			To:        1,
		},
	}
	require.Equal(expected, abc)
}

func TestDelta_DifferentAt(t *testing.T) {
	cases := []struct {
		name    string
		subject diff.Delta
		path    fieldpath.FieldPath
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			fpSpecReplicas,
			false,
		},
		{
			"AB spec.replicas",
			deltaAB,
			fpSpecReplicas,
			true,
		},
		{
			"AB spec.generation",
			deltaAB,
			fpSpecGeneration,
			false,
		},
		{
			"AB status.readyReplicas",
			deltaAB,
			fpStatusReadyReplicas,
			true,
		},
		{
			"AC spec.replicas",
			deltaAC,
			fpSpecReplicas,
			false,
		},
		{
			"AC spec.generation",
			deltaAC,
			fpSpecGeneration,
			false,
		},
		{
			"AC status.readyReplicas",
			deltaAC,
			fpStatusReadyReplicas,
			true,
		},
		{
			"BC spec.replicas",
			deltaBC,
			fpSpecReplicas,
			true,
		},
		{
			"BC spec.generation",
			deltaBC,
			fpSpecGeneration,
			false,
		},
		{
			"BC status.readyReplicas",
			deltaBC,
			fpStatusReadyReplicas,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferentAt(c.path))
		})
	}
}

// Use `Delta.DifferencesAt()` to return the differences at a particular field
// path. All differences at the specified field path "and below" are returned.
//
// If we wanted to get the differences in the `status` field of our Deployment
// maps, we would use `Delta.DifferencesAt(fieldpath.FromString("status"))`:
func ExampleDelta_DifferencesAt() {
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
	statusDiffs := d.DifferencesAt(fieldpath.FromString("status"))
	for _, statusDiff := range statusDiffs {
		fmt.Println(statusDiff)
	}
	// Output:
	// status.observedGeneration REMOVE previous value: 1
	// status.readyReplicas MODIFY 1 -> 2
}

func TestDelta_DifferentExcept(t *testing.T) {
	cases := []struct {
		name    string
		subject diff.Delta
		paths   []fieldpath.FieldPath
		exp     bool
	}{
		{
			"empty delta",
			emptyDelta,
			[]fieldpath.FieldPath{fpSpecReplicas},
			false,
		},
		{
			"AB except spec.replicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas},
			true, // because also diff at status.readyReplicas
		},
		{
			"AB except spec.replicas and status.readyReplicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas, fpStatusReadyReplicas},
			false, // because above paths are the only diffs
		},
		{
			"AC except status.readyReplicas",
			deltaAC,
			[]fieldpath.FieldPath{fpStatusReadyReplicas},
			false, // because above path is the only diffs
		},
		{
			"no supplied field paths and non-empty delta",
			deltaAB,
			nil,
			true, // because there are differences
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferentExcept(c.paths...))
		})
	}
}

func TestDelta_DifferencesExcept(t *testing.T) {
	cases := []struct {
		name    string
		subject diff.Delta
		paths   []fieldpath.FieldPath
		exp     []diff.Difference
	}{
		{
			"empty delta",
			emptyDelta,
			[]fieldpath.FieldPath{fpSpecReplicas},
			[]diff.Difference{},
		},
		{
			"AB except spec.replicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas},
			[]diff.Difference{
				{
					FieldPath: fieldpath.FromString("status.readyReplicas"),
					Kind:      diff.DifferenceModify,
					From:      0,
					To:        1,
				},
			},
		},
		{
			"AB except spec.replicas and status.readyReplicas",
			deltaAB,
			[]fieldpath.FieldPath{fpSpecReplicas, fpStatusReadyReplicas},
			[]diff.Difference{},
		},
		{
			"AC except status.readyReplicas",
			deltaAC,
			[]fieldpath.FieldPath{fpStatusReadyReplicas},
			[]diff.Difference{}, // because above path is the only diffs
		},
		{
			"no supplied field paths and non-empty delta",
			deltaAB,
			nil,
			[]diff.Difference{
				{
					FieldPath: fieldpath.FromString("spec.replicas"),
					Kind:      diff.DifferenceModify,
					From:      1,
					To:        2,
				},
				{
					FieldPath: fieldpath.FromString("status.readyReplicas"),
					Kind:      diff.DifferenceModify,
					From:      0,
					To:        1,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(c.exp, c.subject.DifferencesExcept(c.paths...))
		})
	}
}

// Use `Delta.DifferencesExcept()` to return any differences *other than those
// at a particular field path*.
//
// If we wanted to get the differences between our Deployment maps but did not
// care about the `status` values, we would use
// `Delta.DifferencesExcept(fieldpath.FromString("status"))`:
func ExampleDelta_DifferencesExcept() {
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
	notStatusDiffs := d.DifferencesExcept(fieldpath.FromString("status"))
	for _, notStatusDiff := range notStatusDiffs {
		fmt.Println(notStatusDiff)
	}
	// Output:
	// spec.replicas MODIFY 1 -> 3
}
