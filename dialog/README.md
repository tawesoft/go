# dialog - simple cross-platform messagebox

## About

Package dialog implements simple cross platform native MessageBox/Alert dialogs for Go.

Currently only Windows and Linux targets are supported.

On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT-0][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/dialog
[src_]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_]: https://godoc.org/tawesoft.co.uk/go/dialog
[copy_]: https://github.com/tawesoft/go/tree/master/dialog/COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/dialog
```

## Example:

```go
package main

import "tawesoft.co.uk/go/dialog"

func main() {
    dialog.Alert("Message")
    dialog.Alert("There are %d lights", 4)
}
```