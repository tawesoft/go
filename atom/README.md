# atom - integer codes for strings

## About

Package atom defines an interface and simple implementations for uniquely mapping any set of known-constant strings to
a set of integers for efficient communication and equality operations.

Atoms are used by X11 (https://tronche.com/gui/x/xlib/window-information/XInternAtom.html),
Win32 (https://docs.microsoft.com/en-us/windows/win32/dataxchg/about-atom-tables),
and Go internally (https://godoc.org/golang.org/x/net/html/atom).

This is a very simple interface with a very simple implementation. Features like reference counting atoms, iterating
over atoms, alternative implementations, etc. are deliberately omitted at present.

The exact integer representation of an Atom is opaque. If two Atoms from the same collection compare equal then their
string representations also compare equal. If two strings from the same collection compare equal then their Atom
representations also compare equal.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT-0][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/atom
[src_]:  https://github.com/tawesoft/go/tree/master/atom
[docs_]: https://godoc.org/tawesoft.co.uk/go/atom
[copy_]: https://github.com/tawesoft/go/tree/master/atom/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/atom
```

## Example:

```go
package main

import "tawesoft.co.uk/go/atom"

func main() {
    var atoms = atom.SimpleAtoms()
    
    var atom1 = atoms.Get("Atom One")
    var atom2 = atoms.Get("Atom Two")
    
    if atom1 == atom2 { /* do something ... */ }
}
```