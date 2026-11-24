`github.com/jaypipes/delta/types` contains type inspection and discovery
functions.

# Category

`delta` groups data and variables into these type categories:

* `Scalar` is for things that represent a singular value
* `Slice` is for things contain zero or more values addressable via an
  integer index
* `Record` is for things that contain zero or more values addressable via a
  named key or field

Use `CategoryOf` to determine the `Category` of a subject.

<details>
<summary>Use `CategoryOf` to determine the `Category` of a subject</summary>

```go
package main

import (
    "fmt"

    "github.com/jaypipes/delta/types"
)

type MyIntType int

const (
    MyIntTypeZero MyIntType = 0
)

type MyStructType struct {
    A int
}

func main() {
    subjects := []any{
        int(0),
        uint(0),
        float32(0.0),
        int8(0),
        MyIntType(0),
        []int64{0, 1, 2},
        MyStructType{
            A: 0,
        },
        map[string]string{
            "A": "a",    
        },
    }

    for _, subject := range subjects {
        cat := types.Category(subject)
        fmt.Printf(
            "%v (%T) is types.%s\n",
            subject, cat,
        )
    }
}
```

will output:

```
0 (int) is types.Scalar
0 (uint) is types.Scalar
0 (float32) is types.Scalar
0 (int8) is types.Scalar
0 (MyIntType) is types.Scalar
[0,1,2] ([]int64) is types.Slice
{A:0} (MyStructType) is types.Record
{A:0} (map[string]string) is types.Record
```
</details>
