// Example leveraging the raw parts of FormatParts to handle durations in a
// custom even nicer way for the english language.
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
