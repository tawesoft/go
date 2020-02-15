# operator - operators as functions

## About

Package operator implements logical, arithmetic, bitwise and comparison
operators as functions (like the Python operator module). Includes unary,
binary, and nary functions with overflow checked variants.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT-0][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/operator
[src_]:  https://github.com/tawesoft/go/tree/master/operator
[docs_]: https://godoc.org/tawesoft.co.uk/go/operator
[copy_]: https://github.com/tawesoft/go/tree/master/operator/COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/operator
```

## Example:

```go
package main

import (
    "fmt"
    "tawesoft.co.uk/go/operator"
)

func foo(op func(int, int) int, a int, b int) int {
    return op(a, b)
}

func fooChecked(op func(int8, int8) (int8, error), a int8, b int8) (int8, error) {
    return op(a, b)
}

func main() {
    fmt.Println(foo(operator.Int.Binary.Add, 5, 3))
    fmt.Println(foo(operator.Int.Binary.Sub, 5, 3))
    
    var result, err = fooChecked(operator.Int8Checked.Binary.Add, 126, 2) // max int8 is 127!
    if err != nil {
        fmt.Printf("error: %v (expected!)\n", err)
    } else {
        fmt.Println(result)
    }
}
```