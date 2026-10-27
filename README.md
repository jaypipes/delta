[![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/delta.svg)](https://pkg.go.dev/github.com/jaypipes/delta)
[![Test Status](https://github.com/jaypipes/delta/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jaypipes/delta/actions)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

`delta` is a Go library that facilitates efficient diffing of things.

# Usage

The most basic use of `delta` is to get the difference between two things using
the `delta.Diff()` function which returns a `delta.Delta`:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta"
)

func main() {
    var a = "x"
    var b = "y"

    d, err := delta.Diff(a, b)
    if err != nil {
        fmt.Printf("error diffing vars: %s\n", err)
        return
    }
    fmt.Println(d)
}
```

will output:

```
MODIFY x -> y
```
<details>
<summary>A note about naive comparisons</summary>

The comparison logic for simple scalar variables is deliberately naive. `"1"`
is equal to `1`. `int64(1)` is equal to `uint(1)`, etc. We do this because the
`delta` library was designed to compare JSON data that may have had numbers
encoded as strings and where JSON does not differentiate like Go does between
signed, unsigned and differently-sized integers.

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta"
)

func main() {
	var x = 1
	var y = "1"

	d, err := delta.Diff(x, y)
	if err != nil {
		fmt.Printf("error diffing vars: %s\n", err)
		return
	}
	fmt.Println(d)
}
```

will output:

```
<empty>
```
</details>

While the above is a perfectly acceptable use of `delta`, the library really
shines when comparing non-scalar (map, struct) things.

In the example below, we have two maps representing two simplified Kubernetes
Deployments, `deploy1` and `deploy2`. Use `delta.Diff()` to get the difference
between the two Deployments:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta"
)

func main() {
	var deploy1 = map[string]any{
        "spec": map[string]any{
            "replicas": 1,
        },
        "status": map[string]any{
            "readyReplicas": 1,
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
}
```

will output:

```
spec.replicas MODIFY 1 -> 3
status.observedGeneration REMOVE previous value: 1
status.readyReplicas MODIFY 1 -> 2
```

The differences contained in the `delta.Delta` indicate specific *field paths*.

A field path is a dotted string notation that identifies a specific "field"
value within the non-scalar variable. For example, `status.observedGeneration`
is a field path that refers to the `observedGeneration` field in the `status`
top-level field of the map variables.

## Differences at a field path

Use `Delta.DifferencesAt()` to return the differences at a particular field
path. All differences at the specified field path "and below" are returned.

If we wanted to get the differences in the `status` field of our Deployment
maps, we would use `Delta.DifferencesAt(fieldpath.FromString("status"))`:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta"
    "github.com/jaypipes/delta/fieldpath"
)

func main() {
	var deploy1 = map[string]any{
        "spec": map[string]any{
            "replicas": 1,
        },
        "status": map[string]any{
            "readyReplicas": 1,
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
}
```

will output:

```
status.observedGeneration REMOVE previous value: 1
status.readyReplicas MODIFY 1 -> 2
```

## Differences except at a field path

Use `Delta.DifferencesExcept()` to return any differences *other than those at
a particular field path*.

If we wanted to get the differences between our Deployment maps but did not
care about the `status` values, we would use
`Delta.DifferencesExcept(fieldpath.FromString("status"))`:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta"
    "github.com/jaypipes/delta/fieldpath"
)

func main() {
	var deploy1 = map[string]any{
        "spec": map[string]any{
            "replicas": 1,
        },
        "status": map[string]any{
            "readyReplicas": 1,
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
}
```

will output:

```
spec.replicas MODIFY 1 -> 3
```

# Developers

[Contributions](CONTRIBUTING.md) to `delta` are welcomed! Fork the repo on GitHub
and submit a pull request with your proposed changes. Or, feel free to log an
issue for a feature request or bug report.
