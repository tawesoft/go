# ximage/xcolor - extended color types

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

OpenGL® and the oval logo are trademarks or registered trademarks of Hewlett Packard Enterprise in
the United States and/or other countries worldwide.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [BSD-3-Clause][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/ximage/xcolor
[src_]:  https://github.com/tawesoft/go/tree/master/ximage/xcolor
[docs_]: https://godoc.org/tawesoft.co.uk/go/ximage/xcolor
[copy_]: https://github.com/tawesoft/go/tree/master/ximage/xcolor/COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/ximage/xcolor
```

## See Also:

* ximage (https://tawesoft.co.uk/go/ximage)