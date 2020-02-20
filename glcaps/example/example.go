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
