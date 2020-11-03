# dialog - simple cross-platform messagebox

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/dialog"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_dialog] ∙ [docs][docs_dialog] ∙ [src][src_dialog] | [MIT-0][copy_dialog] | ✔ yes |

[home_dialog]: https://tawesoft.co.uk/go/dialog
[src_dialog]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_dialog]: https://godoc.org/tawesoft.co.uk/go/dialog
[copy_dialog]: https://github.com/tawesoft/go/tree/master/dialog/LICENSE.txt

## About

Package dialog implements simple cross platform native MessageBox/Alert
dialogs for Go.

Currently, only supports Windows and Linux targets.

On Linux, uses (in order of preference) `zenity`, `xmessage`, or stdio.

## Example

Usage is quite simple:

```go
package main

import "tawesoft.co.uk/go/dialog"

func main() {
    dialog.Alert("Hello world!")
    dialog.Alert("There are %d lights", 4)
}
```

## Changes

### 2019-11-16
### 
* Fix incorrect formatting of multiple arguments in Linux stdio fallback
### 
### 2019-10-16
### 
* Remove title argument from Alert function
### 
### 2019-10-01
### 
* Fix string formatting bug in Windows build
### 
### 2019-10-01
### 
* Support Unicode in UTF16 Windows dialogs
* Use "golang.org/x/sys/windows" to provide WinAPI
* Removes CGO and windows.h implementation
* Linux stdio fallback alert no longer blocks waiting for input
### 
### 2019-09-30
### 
* First release
### 