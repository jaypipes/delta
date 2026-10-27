package delta_test

import (
	"fmt"

	"github.com/jaypipes/delta"
	"github.com/jaypipes/delta/fieldpath"
)

func ExampleDiff_simple_scalar_variables() {
	var a = "foo"
	var b = "bar"

	d, err := delta.Diff(a, b)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
	// Output:
	// MODIFY foo -> bar
}

// Note that the comparison logic for simple scalar variables is deliberately
// naive. "1" is equal to 1. int64(1) is equal to uint(1), etc. We do this
// because the delta library was designed to compare JSON data that may have
// had numbers encoded as strings and where JSON does not differentiate like Go
// does between signed, unsigned and differently-sized integers.
func ExampleDiff_naive_comparisons() {
	var x = 1
	var y = "1"

	d, err := delta.Diff(x, y)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
	// Output:
	// <empty>
}

// Two maps representing two Kubernetes Deployments, `deploy1` and `deploy2`.
// Use `delta.Diff()` to get the difference between the two Deployments:
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

	d, err := delta.Diff(deploy1, deploy2)
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

	d, err := delta.Diff(deploy1, deploy2)
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

	d, err := delta.Diff(deploy1, deploy2)
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
