# humanizex - locale-aware natural number formatting

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/humanizex"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_humanizex] ∙ [docs][docs_humanizex] ∙ [src][src_humanizex] | [MIT][copy_humanizex] | ✘ **no** |

[home_humanizex]: https://tawesoft.co.uk/go/humanizex
[src_humanizex]:  https://github.com/tawesoft/go/tree/master/humanizex
[docs_humanizex]: https://www.tawesoft.co.uk/go/doc/humanizex
[copy_humanizex]: https://github.com/tawesoft/go/tree/master/humanizex/LICENSE.txt

## About

Package humanizex is an elegant, general-purpose, extensible, modular,
locale-aware way to format and parse numbers and quantities - like distances,
bytes, and time - in a human-readable way ideal for config files and as a
building-block for fully translated ergonomic user interfaces.

If golang.org/x/text is ever promoted to core then there will be a new version
of this package named `humanize` (dropping the 'x').

What about dustin's go-humanize?

dustin's go-humanize (https://github.com/dustin/go-humanize) is 3.9 to 4.5
times faster formatting and 2 times faster parsing, if this is a bottleneck for
you. It's also quite mature, so is probably very well tested by now. If you're
only targeting the English language it also has more handy "out of the box"
features.

On the other hand, tawesoft's humanizex is more general purpose and has better
localisation support. Even with those extra features, tawesoft's humanizex
codebase is also smaller and simpler.


## Examples


Example formatting and parsing Byte quantities in various locales
```go
package main

import (
    "fmt"

    "golang.org/x/text/language"
    "tawesoft.co.uk/go/humanizex"
)

func mustInt64(v int64, err error) int64 {
    if err != nil { panic(err) }
    return v
}

func main() {
    hEnglish := humanizex.NewHumanizer(language.English)
    hDanish  := humanizex.NewHumanizer(language.Danish)
    hBengali := humanizex.NewHumanizer(language.Bengali)

    // prints 1.5 KiB
    fmt.Println(hEnglish.FormatBytesIEC(1024 + 512))

    // prints 1,5 KiB
    fmt.Println(hDanish.FormatBytesIEC(1024 + 512))

    // prints ১.৫ KiB
    fmt.Println(hBengali.FormatBytesIEC(1024 + 512))

    // prints 1536
    fmt.Println(mustInt64(hEnglish.ParseBytesIEC("1.5 KiB")))
}
```
Example leveraging the raw parts of FormatParts to handle durations in a
custom even nicer way for the english language.
```go
package main

import (
    "fmt"
    "time"

    "golang.org/x/text/language"
    "tawesoft.co.uk/go/humanizex"
)

func plural(x float64) string {
    if x > 0.99 && x < 1.01 { return "" }
    return "s"
}

func main() {

    duration := (2 * time.Hour) + (20 * time.Second)

    // prints "2 h 20 s"
    fmt.Printf("Basic time: %s\n", humanizex.NewHumanizer(language.English).FormatDuration(duration))

    // Get the raw format parts
    parts := humanizex.FormatParts(
        duration.Seconds(),
        humanizex.CommonUnits.Second,
        humanizex.CommonFactors.Time,
    )

    // prints Nice time: 2 hours and 20 seconds ago
    fmt.Printf("Nice time: ")
    if (len(parts) == 1) && (parts[0].Unit.Utf8 == "s") {
        fmt.Printf("just now\n")
    } else {
        for i, part := range parts {
            fmt.Printf("%d", int(part.Magnitude + 0.5))

            if part.Unit.Utf8 == "y" {
                fmt.Printf(" year%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "d" {
                fmt.Printf(" day%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "h" {
                fmt.Printf(" hour%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "min" {
                fmt.Printf(" minute%s", plural(part.Magnitude))
            } else if part.Unit.Utf8 == "s" {
                fmt.Printf(" second%s", plural(part.Magnitude))
            }

            if i + 1 < len(parts) {
                fmt.Printf(" and ")
            } else {
                fmt.Printf(" ago\n")
            }
        }
    }

}
```
Example using custom time factors from the Battlestar Galactica 1978 TV
series.
```go
package main

import (
    "fmt"

    "golang.org/x/text/language"
    "tawesoft.co.uk/go/humanizex"
)

func main() {
    factors := humanizex.Factors{
        Factors:    []humanizex.Factor{
            {1,                         humanizex.Unit{"millicenton", "millicenton"}, humanizex.FactorModeReplace},
            {60,                        humanizex.Unit{"centon", "centon"}, humanizex.FactorModeReplace},
            {60 * 60,                   humanizex.Unit{"centar", "centar"}, humanizex.FactorModeReplace},
            {24 * 60 * 60,              humanizex.Unit{"cycle", "cycle"}, humanizex.FactorModeReplace},
            {7 * 24 * 60 * 60,          humanizex.Unit{"secton", "secton"}, humanizex.FactorModeReplace},
            {28 * 24 * 60 * 60,         humanizex.Unit{"sectar", "sectar"}, humanizex.FactorModeReplace},
            {365 * 24 * 60 * 60,        humanizex.Unit{"yahren", "yahren"}, humanizex.FactorModeReplace},
            {100 * 365 * 24 * 60 * 60,  humanizex.Unit{"centauron", "centauron"}, humanizex.FactorModeReplace},
        },
        Components: 2,
    }

    h := humanizex.NewHumanizer(language.English)

    est := float64((2 * 365 * 24 * 60 * 60) + 1)

    fmt.Printf("Hey, I'll be with you in %s. Watch out for toasters!\n",
        h.Format(est, humanizex.Unit{"millicenton", "millicenton"}, factors).Utf8)

    // prints "Hey, I'll be with you in 2 yahren 1 millicenton. Watch out for toasters!"
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.