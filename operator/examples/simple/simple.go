// Using operators as function arguments
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
