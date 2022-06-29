[![Tawesoft](https://www.tawesoft.co.uk/media/0/logo-240r.png)](https://tawesoft.co.uk/go)
================================================================================

A monorepo for small Go modules maintained by [Tawesoft®](https://www.tawesoft.co.uk/go)

This is permissively-licensed open source software but exact licences may vary between modules.



##FROZEN - PLEASE MIGRATE

These packages are moving to https://github.com/tawesoft/golib.

This is to increase security against possible supply chain attacks such as our domain name expiring in the future and being registered by someone else.

Please migrate to https://github.com/tawesoft/golib (when available) instead.

Most programs relying on a package in this monorepo, such as the dialog or lxstrconv packages, will continue to work for the foreseeable future.

Rarely used packages have been hidden for now - they are in the git commit history at https://github.com/tawesoft/go if you need to resurrect one.



Download
--------

    go get -u tawesoft.co.uk/go

Contents
--------


### dialog - simple cross-platform messagebox

Package dialog implements simple cross platform native MessageBox/Alert
dialogs for Go.

```go
import "tawesoft.co.uk/go/dialog"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_dialog] ∙ [docs][docs_dialog] ∙ [src][src_dialog] | [MIT-0][copy_dialog] | ✔ yes |

[home_dialog]: https://tawesoft.co.uk/go/dialog
[src_dialog]:  https://github.com/tawesoft/go/tree/master/dialog
[docs_dialog]: https://www.tawesoft.co.uk/go/doc/dialog
[copy_dialog]: https://github.com/tawesoft/go/tree/master/dialog/LICENSE.txt
### glcaps - read and check OpenGL capabilities

Package glcaps provides a nice interface to declare OpenGL capabilities you
care about, including minimum required extensions or capabilities. Glcaps has
no dependencies and is agnostic to the exact OpenGL binding used.

```go
import "tawesoft.co.uk/go/glcaps"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_glcaps] ∙ [docs][docs_glcaps] ∙ [src][src_glcaps] | [MIT][copy_glcaps] | ✔ yes |

[home_glcaps]: https://tawesoft.co.uk/go/glcaps
[src_glcaps]:  https://github.com/tawesoft/go/tree/master/glcaps
[docs_glcaps]: https://www.tawesoft.co.uk/go/doc/glcaps
[copy_glcaps]: https://github.com/tawesoft/go/tree/master/glcaps/LICENSE.txt
### humanizex - locale-aware natural number formatting

Package humanizex is an elegant, general-purpose, extensible, modular,
locale-aware way to format and parse numbers and quantities - like distances,
bytes, and time - in a human-readable way ideal for config files and as a
building-block for fully translated ergonomic user interfaces.

```go
import "tawesoft.co.uk/go/humanizex"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_humanizex] ∙ [docs][docs_humanizex] ∙ [src][src_humanizex] | [MIT][copy_humanizex] | ✔ yes |

[home_humanizex]: https://tawesoft.co.uk/go/humanizex
[src_humanizex]:  https://github.com/tawesoft/go/tree/master/humanizex
[docs_humanizex]: https://www.tawesoft.co.uk/go/doc/humanizex
[copy_humanizex]: https://github.com/tawesoft/go/tree/master/humanizex/LICENSE.txt
### lxstrconv - locale-aware number parsing

Package lxstrconv is an attempt at implementing locale-aware parsing of
numbers that integrates with golang.org/x/text.

```go
import "tawesoft.co.uk/go/lxstrconv"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_lxstrconv] ∙ [docs][docs_lxstrconv] ∙ [src][src_lxstrconv] | [MIT][copy_lxstrconv] | ✔ yes |

[home_lxstrconv]: https://tawesoft.co.uk/go/lxstrconv
[src_lxstrconv]:  https://github.com/tawesoft/go/tree/master/lxstrconv
[docs_lxstrconv]: https://www.tawesoft.co.uk/go/doc/lxstrconv
[copy_lxstrconv]: https://github.com/tawesoft/go/tree/master/lxstrconv/LICENSE.txt
### operator - operators as functions

Package operator implements logical, arithmetic, bitwise and comparison
operators as functions (like the Python operator module). Includes unary,
binary, and n-ary functions with overflow checked variants.

```go
import "tawesoft.co.uk/go/operator"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_operator] ∙ [docs][docs_operator] ∙ [src][src_operator] | [MIT-0][copy_operator] | ✔ yes |

[home_operator]: https://tawesoft.co.uk/go/operator
[src_operator]:  https://github.com/tawesoft/go/tree/master/operator
[docs_operator]: https://www.tawesoft.co.uk/go/doc/operator
[copy_operator]: https://github.com/tawesoft/go/tree/master/operator/LICENSE.txt

Links
-----

* Home: [tawesoft.co.uk/go](https://tawesoft.co.uk/go)
* Docs hub: [tawesoft.co.uk/go/doc/](https://www.tawesoft.co.uk/go/doc/)
* Repository: [github.com/tawesoft/go](https://github.com/tawesoft/go)
* Or [search "tawesoft"](https://pkg.go.dev/search?q=tawesoft) on [go.dev](https://go.dev/)

Support
-------

### Free and Community Support

* [GitHub issues](https://github.com/tawesoft/go/issues)
* Email open-source@tawesoft.co.uk (feedback welcomed, but support is "best
 effort")

### Commercial Support

Open source software from Tawesoft® backed by commercial support options.

Email open-source@tawesoft.co.uk or visit [tawesoft.co.uk/products/open-source-software](https://www.tawesoft.co.uk/products/open-source-software)
to learn more.
