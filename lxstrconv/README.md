# lxstrconv - locale-aware number parsing

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/lxstrconv"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_lxstrconv] ∙ [docs][docs_lxstrconv] ∙ [src][src_lxstrconv] | [MIT][copy_lxstrconv] | ✘ **no** |

[home_lxstrconv]: https://tawesoft.co.uk/go/lxstrconv
[src_lxstrconv]:  https://github.com/tawesoft/go/tree/master/lxstrconv
[docs_lxstrconv]: https://www.tawesoft.co.uk/go/doc/lxstrconv
[copy_lxstrconv]: https://github.com/tawesoft/go/tree/master/lxstrconv/LICENSE.txt

## About

Package lxstrconv is an attempt at implementing locale-aware parsing of
numbers that integrates with golang.org/x/text.

If golang.org/x/text is ever promoted to core then there will be a new version
of this package named `lstrconv` (dropping the 'x').


## Package Stability


THIS IS A PREVIEW RELEASE, SUBJECT TO BREAKING CHANGES.

Todo:

* checks for integer overflow

* different representations of negative numbers e.g. `(123)` vs `-123`

* In cases where AcceptInteger/AcceptFloat reach a syntax error, they
currently underestimate how many bytes they successfully parsed when
the byte length of the string is not equal to the number of Unicode
code points in the string.


## Example


This example demonstrates British, Dutch, and Arabic locale number parsing.


```go
package main

import (
    "fmt"
    "golang.org/x/text/language"
    "tawesoft.co.uk/go/lxstrconv"
)

func checked(f float64, e error) float64 {
    if e != nil {
        panic(e)
    }
    return f
}

func main() {
    dutch   := lxstrconv.NewDecimalFormat(language.Dutch)
    british := lxstrconv.NewDecimalFormat(language.BritishEnglish)
    arabic  := lxstrconv.NewDecimalFormat(language.Arabic)

    fmt.Printf("%f\n", checked(british.ParseFloat("1,234.56")))
    fmt.Printf("%f\n", checked(dutch.ParseFloat("1.234,56")))
    fmt.Printf("%f\n", checked(arabic.ParseFloat("١٬٢٣٤٫٥٦")))
}
```

Example

You can give end-users examples of the input you expect for a given locale
using the /x/text package:


```go
package main

import (
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "golang.org/x/text/number"
)

func main() {

    message.NewPrinter(language.English).Println(number.Decimal(123456789))
    // Prints 123,456,789

    message.NewPrinter(language.Dutch).Println(number.Decimal(123456789))
    // Prints 123.456.789

    message.NewPrinter(language.Malayalam).Println(number.Decimal(123456789))
    // Prints 12,34,56,789

    message.NewPrinter(language.Bengali).Println(number.Decimal(123456789))
    // Prints ১২,৩৪,৫৬,৭৮৯
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.