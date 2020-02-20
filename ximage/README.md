# ximage - extended image types

## About

Package ximage implements Red, RG, and RGB images matching the core
image interface.

Note that there are good reasons these image types aren't in the core image
package. The native image types may have optimized fast-paths for many use
cases.

This package is a tradeoff of these optimizations against lower memory
usage. This package is intended to be used in computer graphics (e.g.
OpenGL) where images are uploaded to the GPU in a specific format (such as
GL_R, GL_RG, or GL_RGB) and we don't care too much about the performance of
native Go image manipulation.

OpenGL® and the oval logo are trademarks or registered trademarks of Hewlett Packard Enterprise in
the United States and/or other countries worldwide.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [BSD-3-Clause][copy_] | ✔ yes |

[home_]: https://tawesoft.co.uk/go/ximage
[src_]:  https://github.com/tawesoft/go/tree/master/ximage
[docs_]: https://godoc.org/tawesoft.co.uk/go/ximage
[copy_]: https://github.com/tawesoft/go/tree/master/ximage/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/ximage
```

## See Also:

* ximage/xcolor (https://tawesoft.co.uk/go/ximage/xcolor)