# ximage/xcolor - extended color types

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/ximage/xcolor"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_ximage/xcolor] ∙ [docs][docs_ximage/xcolor] ∙ [src][src_ximage/xcolor] | [BSD-3-Clause][copy_ximage/xcolor] | ✔ yes |

[home_ximage/xcolor]: https://tawesoft.co.uk/go/ximage/xcolor
[src_ximage/xcolor]:  https://github.com/tawesoft/go/tree/master/ximage/xcolor
[docs_ximage/xcolor]: https://www.tawesoft.co.uk/go/doc/ximage/xcolor
[copy_ximage/xcolor]: https://github.com/tawesoft/go/tree/master/ximage/xcolor/LICENSE.txt

## About

Package xcolor implements Red, RedGreen, and RGB color models matching the core
image/color interface.

Note that there are good reasons these color types aren't in the core
image.color package. The native color types may have optimized fast-paths
for many use cases.

This package is a tradeoff of these optimizations against lower memory
usage. This package is intended to be used in computer graphics (e.g.
OpenGL) where images are uploaded to the GPU in a specific format (such as
GL_R, GL_RG, or GL_RGB) and we don't care about the performance of native
Go image manipulation.

OpenGL® and the oval logo are trademarks or registered trademarks of
Hewlett Packard Enterprise
in the United States and/or other countries worldwide.

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.