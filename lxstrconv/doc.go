// SPDX-License-Identifier: MIT-0
//
// tawesoft.co.uk/go/lxstrconv
//
// Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
// Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction,  including without limitation the rights
// to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
// copies  of  the  Software,  and  to  permit persons  to whom  the Software is
// furnished to do so.
//
// THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
// IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
//
// Usage:
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
