package main

import (
    "tawesoft.co.uk/go/dialog"
)

func main() {
    // test a windows version:
    // CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1  go build -trimpath dialog/example/main.go
    // wine ./main.exe

    dialog.Alert("Hello world!")
    dialog.Alert("There are %d lights", 4)
    dialog.Alert("Hello %d world!") // safe because there are no args
    dialog.Alert("Unicode £GBP €EUR")
}
