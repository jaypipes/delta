# Compare

* [`StrictEqual`](#strict-equal) returns true if two arguments are
  [strict equal comparable](#strict-equal-comparable) and have equal values.
* [`Equal`](#equal) returns true if two arguments are
  [equal comparable](#equal-comparable) and have equal values.
* [`Equivalent`](#equivalent) returns true if two arguments are
  [equivalent comparable](#equivalent-comparable) and have equivalent values.

## `StrictEqual`

`StrictEqual` returns a bool indicating two arguments are [strict equal
comparable](#strict-equal-comparable) have equal values.



`rune` and `int32` are the exact same type since `rune` is a type alias for
`int32`.

`byte` and `uint8` are the exact same type since `byte` is a type alias for
`uint8`.

`string` and `[]byte` are considered the exact same type since `string` is an
immutable `[]byte`.

If both arguments are `nil`, returns `true`.

> [!CAUTION]
> `nil` is not the same as an empty interface that points to `nil`.

Integer and Float types of different sizes and sign are **different types**.
`uint` and `int` are different types. `int8` and `int16` are different types.
`float32` and `float64` are different types.

`int` and `int64`/`int32` are different types because `int` has a different
size on different underlying architectures. Same for `uint` and
`uint64`/`uint32`

Use `Equal` instead of `StrictEqual` if your comparison logic considers all
Integer types to be fundamentally the same.

Use `Equivalent` instead of `Equal` when your code needs to compare potentially
different sizes of Float types.

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta/cmp"
)

func main() {
	s1 := "1"
	bs1 := []byte("1")

	fmt.Printf("Are s1 and bs1 strictly equal? %t\n", compare.StrictEqual(s1, bs1))

	r1 := rune('1')
	i32 := int32('1')

	fmt.Printf("Are r1 and i32 strictly equal? %t\n", compare.StrictEqual(r1, i32))

	b1 := byte('1')
	u8 := uint8('1')

	fmt.Printf("Are b1 and u8 strictly equal? %t\n", compare.StrictEqual(b1, u8))

	i1 := 1
	ui1 := uint(1)

	fmt.Printf("Are i1 and ui1 strictly equal? %t\n", compare.StrictEqual(i1, ui1))

	f32 := float32(1)
	f64 := float64(1)

	fmt.Printf("Are f32 and f64 strictly equal? %t\n", compare.StrictEqual(f32, f64))
}
```

will output:

```
Are s1 and bs1 strictly equal? true
Are r1 and i32 strictly equal? true
Are b1 and u8 strictly equal? true
Are i1 and ui1 strictly equal? false
Are f32 and f64 strictly equal? false
```

# `Equal`

`Equal` returns a bool indicating whether two [equal
comparable][#equal-comparable] arguments are considered equal.

If the types of the two arguments are not equal comparable, `Equal` returns
false.

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta/cmp"
)

func main() {
    s1 := "1"
    i1 := 1

    b1 := true
    f1 := float32(1)

    fmt.Printf("Are s1 and i1 equal? %t\n", compare.Equal(s1, i1))
    fmt.Printf("Are s1 and b1 equal? %t\n", compare.Equal(s1, b1))
    fmt.Printf("Are s1 and f1 equal? %t\n", compare.Equal(s1, f1))
}
```

will output:

```
Are s1 and i1 equal? t
Are s1 and b1 equal? t
Are s1 and f1 equal? f
```

# `Equivalent`

`Equivalent` returns a bool indicating whether two [equivalent
comparable][#equivalent-comparable] arguments are considered equivalent.

If the types of the two arguments are not equivalent comparable, `Equivalent`
returns false.

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta/cmp"
)

func main() {
    s1 := "1"
    i1 := 1

    b1 := true
    f1 := float32(1)

    fmt.Printf("Are s1 and i1 equal? %t\n", compare.Equal(s1, i1))
    fmt.Printf("Are s1 and b1 equal? %t\n", compare.Equal(s1, b1))
    fmt.Printf("Are s1 and f1 equal? %t\n", compare.Equal(s1, f1))
}
```

will output:

```
Are s1 and i1 equal? t
Are s1 and b1 equal? t
Are s1 and f1 equal? t
```

## Equal comparable

Variables that are equal comparable are comparable using simple type casts.

### Integer types

When both variables being compared are Integer types, but the two variables
being compared are not the exact same type, the variables are cast to the
greatest size of Integer corresponding to the sign of the first variable.

For example, if the first variable being compared is a `uint8` and the second
variable being compared is an `int64`, both variables are type-cast to a
`uint64` and then compared.

If the first variable being compared is an `int64` and the second variable
being compared is a `byte` (which is a type alias for `uint8`), the second
variable is type-cast to an `int64` and the value of the variables is then
compared.

### String types

Because `string` is simply an immutable representation of a `[]byte`, variables
of either `string` or `[]byte` are type-cast compatible and therefore equal
comparable.

### Float types

Float types of different sizes are not equal comparable because type-casting
between a `float64` and `float32` potentially loses precision. Use `Equivalent`
when your code needs to compare potentially different sizes of Float types.

### Bool

Bool types are not equal comparable with anything other than another Bool type.
Type-casting between a `bool` and anthing else loses (almost all) data. Use
`Equivalent` when your code needs to compare strings like `"t"` or integers
like `0` to boolean values `true` and `false`.

## Equivalent comparable

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
