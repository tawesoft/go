// Using operators in lookup tables for a command-line calculator program
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
