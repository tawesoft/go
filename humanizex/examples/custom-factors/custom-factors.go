// Example using custom time factors from the Battlestar Galactica 1978 TV
// series.
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
