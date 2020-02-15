/*
Package dialog implements simple cross platform native MessageBox/Alert
dialogs for Go.

Currently only Windows and Linux targets are supported.

On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio.

Usage:

    package main

    import "tawesoft.co.uk/go/dialog"

    func main() {
        dialog.Alert("Message")
        dialog.Alert("There are %d lights", 4)
    }

Home page https://tawesoft.co.uk/go

For source code see https://github.com/tawesoft/go/tree/master/dialog

For documentation see https://godoc.org/tawesoft.co.uk/go/dialog

*/
package dialog // import "tawesoft.co.uk/go/dialog"

