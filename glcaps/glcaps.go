// Package glcaps provides a nice interface to declare OpenGL capabilities you care about, including minimum required
// extensions or capabilities. Glcaps has no dependencies and is agnostic to the exact OpenGL binding used.
//
// Example:
//
//     package main
//
//    import (
//        "fmt"
//
//        "github.com/go-gl/gl/v3.3-core/gl"
//        "github.com/go-gl/glfw/v3.2/glfw"
//        "tawesoft.co.uk/go/glcaps"
//    )
//
//    func start() func() {
//        var err = gl.Init()
//        if err != nil { panic(err.Error()) }
//
//        err = glfw.Init()
//        if err != nil { panic(err.Error()) }
//
//        glfw.WindowHint(glfw.Visible, glfw.False)
//
//        window, err := glfw.CreateWindow(640, 480, "Example", nil, nil)
//        if err != nil { panic(err.Error()) }
//
//        window.MakeContextCurrent()
//
//        return glfw.Terminate
//    }
//
//    func main() {
//        var closer = start()
//        defer closer()
//
//        type Caps struct {
//            Supports struct {
//                NPOTTextures           bool `glcaps:"ext GL_ARB_texture_non_power_of_two"; required"`
//                BPTextureCompression   bool `glcaps:"ext GL_ARB_texture_compression_bptc; required"`
//                BigTextures            bool `glcaps:"gte GetIntegerv GL_MAX_TEXTURE_SIZE 8192"`
//                AnisotropicFiltering   bool `glcaps:"and ext GL_EXT_texture_filter_anisotropic gte GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
//                FluxCapacitor          bool `glcaps:"and ext FLUX1 ext FLUX2; required"`
//            }
//
//            MaxTextureUnits            int     `glcaps:"GetIntegerv GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS"`
//            MaxTextureSize             int     `glcaps:"GetIntegerv GL_MAX_TEXTURE_SIZE; gte 8192"`
//            MaxAnisotropy              float32 `glcaps:"if ext GL_EXT_texture_filter_anisotropic GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
//            Frobbinators               int    `glcaps:"150; gte 10 lt 100 neq 13"`
//        }
//
//        var Binding = glcaps.Binding{
//            GetIntegerv: gl.GetIntegerv,
//            GetFloatv:   gl.GetFloatv,
//            GetStringi: func(name uint32, index uint32) string {
//                return gl.GoStr(gl.GetStringi(name, index))
//            },
//        }
//
//        var MyCaps Caps
//
//        var extensions, errors = glcaps.Parse(&Binding, &MyCaps)
//        for _, i := range errors {
//            fmt.Printf("glcaps error: %s\n", i.Message)
//        }
//
//        fmt.Printf("Supports.TextureCompressionBPTC: %t\n", MyCaps.Supports.BPTextureCompression)
//        fmt.Printf("Supports.FluxCapacitor: %t\n", MyCaps.Supports.FluxCapacitor)
//        fmt.Printf("Supports.BigTextures: %t\n", MyCaps.Supports.BigTextures)
//        fmt.Printf("MaxTextureUnits: %d\n", MyCaps.MaxTextureUnits)
//        fmt.Printf("Frobbinators: %d\n", MyCaps.Frobbinators)
//        fmt.Printf("%d extensions supported\n", len(extensions))
//    }
//
package glcaps

import (
    "sort"
)

// Extensions is an ordered list of supported OpenGL extensions.
type Extensions []string

// HaveExtension returns true iff the ordered list of supported OpenGL extensions contains a given extension.
func (extensions Extensions) Contains(key string) bool {
    var index = sort.SearchStrings(extensions, key)
    return (index < len(extensions)) && (extensions[index] == key)
}

// Binding implements a binding between this package and a specific OpenGL implementation (e.g. a specific `go-gl`
// module).
type Binding struct {
    GetIntegerv func(name uint32, data *int32)
    GetFloatv   func(name uint32, data *float32)
    GetStringi  func(name uint32, index uint32) string // required to return a Go string, not a C string!
}

// QueryExtensions returns all extensions supported by the current OpenGL context as a sorted list of strings. It is an
// error to call this method if a current OpenGL context does not exist.
func (b *Binding) QueryExtensions() Extensions {
    var numExtensions int32
    b.GetIntegerv(glconstants["GL_NUM_EXTENSIONS"], &numExtensions)
    if numExtensions <= 0 {
        panic("failed to query OpenGL extensions (is the OpenGL context current?)")
    }

    var xs = make([]string, 0, numExtensions)
    
    for i := uint32(0); i < uint32(numExtensions); i++ {
        x := b.GetStringi(glconstants["GL_EXTENSIONS"], i)
        if len(x) == 0 { continue }
        xs = append(xs, x)
    }

    sort.Strings(xs)
    return xs
}

// Error implements an error result type for reporting a capability that doesn't meet a requirement.
type Error struct {
    Field       string // the name of the field in the struct that failed
    Tag         string // the original tag string
    Requirement requirement // the requirement that failed
    Message     string // a human-readable message
}

type Errors []Error

func (es *Errors) append(e ... Error) {
    if *es == nil && (len(e) > 0) {
        *es = make([]Error, 0)
    }
    
    *es = append(*es, e...)
}




