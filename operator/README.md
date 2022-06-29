# operator - operators as functions

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/operator"
```




## FROZEN - PLEASE MIGRATE

These packages are moving to https://github.com/tawesoft/golib.

This is to increase security against possible supply chain attacks such as our domain name expiring in the future and being registered by someone else.

Please migrate to https://github.com/tawesoft/golib (when available) instead.

Most programs relying on a package in this monorepo, such as the dialog or lxstrconv packages, will continue to work for the foreseeable future.

Rarely used packages have been hidden for now - they are in the git commit history at https://github.com/tawesoft/go if you need to resurrect one.



|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_operator] ∙ [docs][docs_operator] ∙ [src][src_operator] | [MIT-0][copy_operator] | ✔ yes |

[home_operator]: https://tawesoft.co.uk/go/operator
[src_operator]:  https://github.com/tawesoft/go/tree/master/operator
[docs_operator]: https://www.tawesoft.co.uk/go/doc/operator
[copy_operator]: https://github.com/tawesoft/go/tree/master/operator/LICENSE.txt

## About

Package operator implements logical, arithmetic, bitwise and comparison
operators as functions (like the Python operator module). Includes unary,
binary, and n-ary functions with overflow checked variants.


## Examples


Using operators as function arguments
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
Using operators in lookup tables for a command-line calculator program
```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    
    "tawesoft.co.uk/go/operator"
)

type checkedOperation func(float64, float64) (float64, error)

var reader = bufio.NewReader(os.Stdin)

var operations = map[string]checkedOperation {
    "+": operator.Float64Checked.Binary.Add,
    "-": operator.Float64Checked.Binary.Sub,
    "*": operator.Float64Checked.Binary.Mul,
    "/": operator.Float64Checked.Binary.Div,
}

func getNumber(prompt string) float64 {
    for {
        fmt.Print(prompt)
        var text, _ = reader.ReadString('\n')
        var result, err = strconv.ParseFloat(strings.TrimSpace(text), 64)
        if err != nil {
            fmt.Println("Sorry, try again. (%v)", err)
            continue
        }
        return result
    }
}

func getOperation(prompt string) checkedOperation {
    for {
        fmt.Print(prompt)
        var text, _ = reader.ReadString('\n')
        var operator, ok = operations[strings.TrimSpace(text)]
        if !ok {
            fmt.Println("Sorry, try again.")
            continue
        }
        return operator
    }
}

func main() {
    var firstNumber = getNumber("Enter a number (then press enter): ")
    var operation = getOperation("Enter +, -, * or / (then press enter) for add, subtract, multiply, or divide: ")
    var secondNumber = getNumber("Enter another number (then press enter): ")
    var result, err = operation(firstNumber, secondNumber)
    if err != nil {
        fmt.Printf("Sorry, something went wrong: %v\n", err)
    } else {
        fmt.Printf("The result is %.2f!\n", result)
    }
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.