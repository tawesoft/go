# glcaps - read and check OpenGL capabilities

```shell script
go get -u "tawesoft.co.uk/go"
```

```go
import "tawesoft.co.uk/go/glcaps"
```




## FROZEN - PLEASE MIGRATE

These packages are moving to https://github.com/tawesoft/golib.

This is to increase security against possible supply chain attacks such as our domain name expiring in the future and being registered by someone else.

Please migrate to https://github.com/tawesoft/golib (when available) instead.

Most programs relying on a package in this monorepo, such as the dialog or lxstrconv packages, will continue to work for the foreseeable future.

Rarely used packages have been hidden for now - they are in the git commit history at https://github.com/tawesoft/go if you need to resurrect one.



|  Links  | License | Stable? |
|:-------:|:-------:|:-------:|
| [home][home_glcaps] ∙ [docs][docs_glcaps] ∙ [src][src_glcaps] | [MIT][copy_glcaps] | ✔ yes |

[home_glcaps]: https://tawesoft.co.uk/go/glcaps
[src_glcaps]:  https://github.com/tawesoft/go/tree/master/glcaps
[docs_glcaps]: https://www.tawesoft.co.uk/go/doc/glcaps
[copy_glcaps]: https://github.com/tawesoft/go/tree/master/glcaps/LICENSE.txt

## About

Package glcaps provides a nice interface to declare OpenGL capabilities you
care about, including minimum required extensions or capabilities. Glcaps has
no dependencies and is agnostic to the exact OpenGL binding used.

OpenGL® and the oval logo are trademarks or registered trademarks of Hewlett
Packard Enterprise in the United States and/or other countries worldwide.


## Examples


Example using glcaps with an OpenGL binding and a struct with tags.

See https://godoc.org/tawesoft.co.uk/go/glcaps#Parse for a description of the
struct annotation syntax.
```go
package main

// Example output:
//
//     glcaps error: FluxCapacitor is required
//     glcaps error: Frobbinators is 150 but must be < 100
//     Supports.TextureCompressionBPTC: true
//     Supports.FluxCapacitor: false
//     Supports.BigTextures: true
//     MaxTextureUnits: 192
//     Frobbinators: 150
//     380 extensions supported
//
//     Info struct { Version string "glcaps:\"GetString GL_VERSION\"";
//     GLSLVersion string "glcaps:\"GetString GL_SHADING_LANGUAGE_VERSION\"";
//     Vendor string "glcaps:\"GetString GL_VENDOR\"";
//     Renderer string "glcaps:\"GetString GL_RENDERER\"" }{
//         Version:"4.6.0 NVIDIA 4??.??",
//         GLSLVersion:"4.60 NVIDIA",
//         Vendor:"NVIDIA Corporation",
//         Renderer:"GeForce GTX ????/PCIe/SSE2",
//     }

import (
    "fmt"

    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "tawesoft.co.uk/go/glcaps"
)

func start() func() {
    var err = gl.Init()
    if err != nil { panic(err.Error()) }

    err = glfw.Init()
    if err != nil { panic(err.Error()) }

    glfw.WindowHint(glfw.Visible, glfw.False)

    window, err := glfw.CreateWindow(640, 480, "Example", nil, nil)
    if err != nil { panic(err.Error()) }

    window.MakeContextCurrent()

    return glfw.Terminate
}

func main() {
    var closer = start()
    defer closer()

    type Caps struct {
        Info struct {
            Version                 string `glcaps:"GetString GL_VERSION"`
            GLSLVersion             string `glcaps:"GetString GL_SHADING_LANGUAGE_VERSION"`
            Vendor                  string `glcaps:"GetString GL_VENDOR"`
            Renderer                string `glcaps:"GetString GL_RENDERER"`
        }

        Supports struct {
            NPOTTextures            bool `glcaps:"ext GL_ARB_texture_non_power_of_two; required"`
            BPTextureCompression    bool `glcaps:"ext GL_ARB_texture_compression_bptc; required"`
            BigTextures             bool `glcaps:"gte GetIntegerv GL_MAX_TEXTURE_SIZE 8192"`
            AnisotropicFiltering    bool `glcaps:"and ext GL_EXT_texture_filter_anisotropic gte GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
            FluxCapacitor           bool `glcaps:"and ext FLUX1 ext FLUX2; required"`
        }

        MaxTextureUnits             int     `glcaps:"GetIntegerv GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS"`
        MaxTextureSize              int     `glcaps:"GetIntegerv GL_MAX_TEXTURE_SIZE; gte 8192"`
        MaxAnisotropy               float32 `glcaps:"if ext GL_EXT_texture_filter_anisotropic GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
        Frobbinators                int     `glcaps:"150; gte 10 lt 100 neq 13"`
    }

    var Binding = glcaps.Binding{
        GetIntegerv: gl.GetIntegerv,
        GetFloatv:   gl.GetFloatv,
        GetString:   func(name uint32) string {
            return gl.GoStr(gl.GetString(name))
        },
        GetStringi:  func(name uint32, index uint32) string {
            return gl.GoStr(gl.GetStringi(name, index))
        },
    }

    var MyCaps Caps

    var extensions, errors = glcaps.Parse(&Binding, &MyCaps)
    for _, i := range errors {
        fmt.Printf("glcaps error: %s\n", i.Message)
    }

    fmt.Printf("Supports.TextureCompressionBPTC: %t\n", MyCaps.Supports.BPTextureCompression)
    fmt.Printf("Supports.FluxCapacitor: %t\n", MyCaps.Supports.FluxCapacitor)
    fmt.Printf("Supports.BigTextures: %t\n", MyCaps.Supports.BigTextures)
    fmt.Printf("MaxTextureUnits: %d\n", MyCaps.MaxTextureUnits)
    fmt.Printf("Frobbinators: %d\n", MyCaps.Frobbinators)
    fmt.Printf("%d extensions supported\n", len(extensions))

    fmt.Printf("\nInfo %#v\n", MyCaps.Info)
}
```

## Getting Help

This package is part of [tawesoft.co.uk/go](https://www.tawesoft.co.uk/go),
a monorepo for small Go modules maintained by Tawesoft®.
Check out that URL for more information about other Go modules from
Tawesoft plus community and commercial support options.