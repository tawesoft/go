# ximage - extended image types

```shell script
go get "tawesoft.co.uk/go/"
```

```go
import "tawesoft.co.uk/go/ximage"
```

|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_ximage] ∙ [docs][docs_ximage] ∙ [src][src_ximage] | [BSD-3-Clause][copy_ximage] | ✔ yes |

[home_ximage]: https://tawesoft.co.uk/go/ximage
[src_ximage]:  https://github.com/tawesoft/go/tree/master/ximage
[docs_ximage]: https://www.tawesoft.co.uk/go/doc/ximage
[copy_ximage]: https://github.com/tawesoft/go/tree/master/ximage/LICENSE.txt

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

OpenGL® and the oval logo are trademarks or registered trademarks of
Hewlett Packard Enterprise
in the United States and/or other countries worldwide.

See also: ximage/xcolor (https://tawesoft.co.uk/go/ximage/xcolor)

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.