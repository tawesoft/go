tawesoft.co.uk/go/dialog
========================

Simple cross platform¹ MessageBox/Alert dialogs for Go.

¹ currently only Windows and Linux

Stable API, but more features are a work in progress.

On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio

Download
--------

```sh
go get tawesoft.co.uk/go/dialog
```

Usage
-----

```go
package main

import "tawesoft.co.uk/go/dialog"

func main() {
    dialog.Alert("Message")
    dialog.Alert("There are %d lights", 4)
}
```

Links
-----

* Home: https://www.tawesoft.co.uk/go/dialog
* Source code: https://github.com/tawesoft/dialog
* Documentation: https://godoc.org/tawesoft.co.uk/go/dialog
