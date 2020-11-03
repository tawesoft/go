// SPDX-License-Identifier: MIT-0

// Package lxstrconv is an attempt at implementing locale-aware parsing of
// numbers.
//
// THIS IS A PREVIEW RELEASE, SUBJECT TO BREAKING CHANGES.
//
// This package integrates with `golang.org/x/text`. If that package is
// ever promoted to core then there will be a new version of this package
// named `lstrconv` (dropping the 'x').
//
// TODO: checks for integer overflow
//
// TODO: different representation of negative numbers e.g. `(123)` vs `-123`
//
// TODO: In cases where AcceptInteger/AcceptFloat reach a syntax error, they
// currently underestimate how many bytes they successfully parsed when
// the byte length of the string is not equal to the number of Unicode
// code points in the string.
//
// Usage: (playground: https://play.golang.org/p/PFZkjOLhoRb)
//
//    package main
//
//    import (
//        "fmt"
//        "golang.org/x/text/language"
//        "tawesoft.co.uk/go/lxstrconv"
//    )
//
//    func main() {
//        f := lxstrconv.NewDecimalFormat(language.French)
//
//        value, err := f.ParseFloat("1 234,56")
//        fmt.Printf("%f %v\n", value, err)
//    }
//
//
//
// You can give end-users examples of the input you expect for a given locale
// using the /x/text package (Playground: https://play.golang.org/p/zcj6ariGMX5):
//
//    message.NewPrinter(language.English).Print(number.Decimal(123456789))
//    // Prints: 123,456,789
//
package lxstrconv
