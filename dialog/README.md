# dialog - simple cross-platform messagebox

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/dialog"
```




##FROZEN - PLEASE MIGRATE

These packages are moving to https://github.com/tawesoft/golib.

This is to increase security against possible supply chain attacks such as our domain name expiring in the future and being registered by someone else.

Please migrate to https://github.com/tawesoft/golib (when available) instead.

Most programs relying on a package in this monorepo, such as the dialog or lxstrconv packages, will continue to work for the foreseeable future.

Rarely used packages have been hidden for now - they are in the git commit history at https://github.com/tawesoft/go if you need to resurrect one.



|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_dialog] ∙ [docs][docs_dialog] ∙ [src][src_dialog] | [MIT-0][copy_dialog] | ✔ yes |

[home_dialog]: https://tawesoft.co.uk/go/dialog
[src_dialog]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_dialog]: https://www.tawesoft.co.uk/go/doc/dialog
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

### 2022-06-29

* Update dialog to never use a format string with no args

### 2019-11-16

* Fix incorrect formatting of multiple arguments in Linux stdio fallback

### 2019-10-16

* Remove title argument from Alert function

### 2019-10-01

* Fix string formatting bug in Windows build

### 2019-10-01

* Support Unicode in UTF16 Windows dialogs
* Use "golang.org/x/sys/windows" to provide WinAPI
* Removes CGO and windows.h implementation
* Linux stdio fallback alert no longer blocks waiting for input

### 2019-09-30

* First release


## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.