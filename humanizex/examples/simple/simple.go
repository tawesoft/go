// Example formatting 1536 Bytes (1.5 KiB) in various locales
package main

import (
    "fmt"

    "golang.org/x/text/language"
    "tawesoft.co.uk/go/humanizex"
)

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
}
