[![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/delta.svg)](https://pkg.go.dev/github.com/jaypipes/delta)
[![Test Status](https://github.com/jaypipes/delta/actions/workflows/test-unit.yml/badge.svg?branch=main)](https://github.com/jaypipes/delta/actions)
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

## Scalar equivalency

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

Scalar equivalency is more formally defined in the `delta` library by the
following rules.

1. `string` or `[]byte` variables containing an integer are equivalent to an
   integer typed variable with the same value.
 - `"2"` is equivalent to `int(2)`
 - `"4"` is equivalent to `uint16(4)`

2. A `string` variable is equivalent to a `[]byte` variable with the same runes.
 - `[]byte("a")` is equivalent to `"a"`.
 - `"1"` is equivalent to `[]byte("1")`

3. Integer variables of any size or sign are equivalent to each other after
  being converted to an `int64`
 - `int8(0)` is equivalent to `uint16(0)`
 - `uint32(2048)` is equivalent to `2048` (implicit `int`)

4. Integer variables are never equivalent to a float variable.

5. `bool` variables are equivalent to `string` or `[]byte` variables using the
  `ParseBool` rules for parsing strings to booleans.
 - `true` is equivalent to any of `"1"`, `"true"`, `"t"` or `"TRUE"` along with
   their `[]byte` versions.
 - `false` is equivalent to any of `"0"`, `"false"`, `"f"` or `"FALSE"` along with
   their `[]byte` versions.

6. `bool` variables are equivalent to Integer variables only when the underlying
  Integer variable equals `0` or `1`.
 - `true` is equivalent to any of `1`, `int8(1)`, `uint64(1)`, etc
 - `false` is equivalent to any of `0`, `int16(0)`, `uint32(0)`, etc

7. `bool` variables are never equivalent to a float variable.

8. `string` or `[]byte` variables containing a floating point number are
   equivalent to a `float32` or `float64` variable when that float variable is
converted into a string with the same number of decimal point places as the
`string` or `[]byte` variable being compared to.
 - `"3.1415"` is equivalent to `float32(3.1415)` and`float64(3.1415000101)`
 - `[]byte("123.01")` is equivalent to `float32(123.0100112)`
   and`float64(123.01000000001110999)`

8. Float variables are never equivalent to another float variable.

> **NOTE**: Comparing float types to float types always returns false because
> there is no way to determine the number of decimal places of precision to use
> when comparing equivalency. This is in contrast to comparing a string to a
> float. When comparing a string to a float, we can determine the number of
> decimal places in the string representation of the floating point number and
> convert the float to a string at that decimal place precision.

> **Important**: equivalence is transitive. if `a` is equivalent to `b`, then
> `b` must be equivalent to `a`.

# Developers

[Contributions](CONTRIBUTING.md) to `delta` are welcomed! Fork the repo on GitHub
and submit a pull request with your proposed changes. Or, feel free to log an
issue for a feature request or bug report.
