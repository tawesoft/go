// tawesoft.co.uk/go/glcaps
// 
// Copyright © 2020 Tawesoft Ltd <open-source@tawesoft.co.uk>
// Copyright © 2020 Ben Golightly <ben@tawesoft.co.uk>
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction,  including without limitation the rights
// to use,  copy, modify,  merge,  publish, distribute, sublicense,  and/or sell
// copies  of  the  Software,  and  to  permit persons  to whom  the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED  "AS IS",  WITHOUT WARRANTY OF ANY KIND,  EXPRESS OR
// IMPLIED,  INCLUDING  BUT  NOT LIMITED TO THE WARRANTIES  OF  MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE  AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS  OR COPYRIGHT HOLDERS  BE LIABLE  FOR ANY  CLAIM,  DAMAGES  OR  OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package glcaps provides a nice interface to declare OpenGL capabilities you
// care about, including minimum required extensions or capabilities. Glcaps has
// no dependencies and is agnostic to the exact OpenGL binding used.
// 
// OpenGL® and the oval logo are trademarks or registered trademarks of Hewlett
// Packard Enterprise in the United States and/or other countries worldwide.
// 
// Example
// 
// Usage involves defining an OpenGL binding and parsing into an annotated struct.
// 
// See https://godoc.org/tawesoft.co.uk/go/glcaps#Parse for a description of the
// struct annotation syntax.
// 
//     package main
// 
//     // Example output:
//     //
//     //     glcaps error: FluxCapacitor is required
//     //     glcaps error: Frobbinators is 150 but must be < 100
//     //     Supports.TextureCompressionBPTC: true
//     //     Supports.FluxCapacitor: false
//     //     Supports.BigTextures: true
//     //     MaxTextureUnits: 192
//     //     Frobbinators: 150
//     //     380 extensions supported
//     //
//     //     Info struct { Version string "glcaps:\"GetString GL_VERSION\"";
//     //     GLSLVersion string "glcaps:\"GetString GL_SHADING_LANGUAGE_VERSION\"";
//     //     Vendor string "glcaps:\"GetString GL_VENDOR\"";
//     //     Renderer string "glcaps:\"GetString GL_RENDERER\"" }{
//     //         Version:"4.6.0 NVIDIA 4??.??",
//     //         GLSLVersion:"4.60 NVIDIA",
//     //         Vendor:"NVIDIA Corporation",
//     //         Renderer:"GeForce GTX ????/PCIe/SSE2",
//     //     }
// 
//     import (
//         "fmt"
// 
//         "github.com/go-gl/gl/v3.3-core/gl"
//         "github.com/go-gl/glfw/v3.2/glfw"
//         "tawesoft.co.uk/go/glcaps"
//     )
// 
//     func start() func() {
//         var err = gl.Init()
//         if err != nil { panic(err.Error()) }
// 
//         err = glfw.Init()
//         if err != nil { panic(err.Error()) }
// 
//         glfw.WindowHint(glfw.Visible, glfw.False)
// 
//         window, err := glfw.CreateWindow(640, 480, "Example", nil, nil)
//         if err != nil { panic(err.Error()) }
// 
//         window.MakeContextCurrent()
// 
//         return glfw.Terminate
//     }
// 
//     func main() {
//         var closer = start()
//         defer closer()
// 
//         type Caps struct {
//             Info struct {
//                 Version                 string `glcaps:"GetString GL_VERSION"`
//                 GLSLVersion             string `glcaps:"GetString GL_SHADING_LANGUAGE_VERSION"`
//                 Vendor                  string `glcaps:"GetString GL_VENDOR"`
//                 Renderer                string `glcaps:"GetString GL_RENDERER"`
//             }
// 
//             Supports struct {
//                 NPOTTextures            bool `glcaps:"ext GL_ARB_texture_non_power_of_two; required"`
//                 BPTextureCompression    bool `glcaps:"ext GL_ARB_texture_compression_bptc; required"`
//                 BigTextures             bool `glcaps:"gte GetIntegerv GL_MAX_TEXTURE_SIZE 8192"`
//                 AnisotropicFiltering    bool `glcaps:"and ext GL_EXT_texture_filter_anisotropic gte GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
//                 FluxCapacitor           bool `glcaps:"and ext FLUX1 ext FLUX2; required"`
//             }
// 
//             MaxTextureUnits             int     `glcaps:"GetIntegerv GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS"`
//             MaxTextureSize              int     `glcaps:"GetIntegerv GL_MAX_TEXTURE_SIZE; gte 8192"`
//             MaxAnisotropy               float32 `glcaps:"if ext GL_EXT_texture_filter_anisotropic GetFloatv GL_MAX_TEXTURE_MAX_ANISOTROPY 1.0"`
//             Frobbinators                int     `glcaps:"150; gte 10 lt 100 neq 13"`
//         }
// 
//         var Binding = glcaps.Binding{
//             GetIntegerv: gl.GetIntegerv,
//             GetFloatv:   gl.GetFloatv,
//             GetString:   func(name uint32) string {
//                 return gl.GoStr(gl.GetString(name))
//             },
//             GetStringi:  func(name uint32, index uint32) string {
//                 return gl.GoStr(gl.GetStringi(name, index))
//             },
//         }
// 
//         var MyCaps Caps
// 
//         var extensions, errors = glcaps.Parse(&Binding, &MyCaps)
//         for _, i := range errors {
//             fmt.Printf("glcaps error: %s\n", i.Message)
//         }
// 
//         fmt.Printf("Supports.TextureCompressionBPTC: %t\n", MyCaps.Supports.BPTextureCompression)
//         fmt.Printf("Supports.FluxCapacitor: %t\n", MyCaps.Supports.FluxCapacitor)
//         fmt.Printf("Supports.BigTextures: %t\n", MyCaps.Supports.BigTextures)
//         fmt.Printf("MaxTextureUnits: %d\n", MyCaps.MaxTextureUnits)
//         fmt.Printf("Frobbinators: %d\n", MyCaps.Frobbinators)
//         fmt.Printf("%d extensions supported\n", len(extensions))
// 
//         fmt.Printf("\nInfo %#v\n", MyCaps.Info)
//     }
//
// Package Information
//
// License: MIT (see LICENSE.txt)
//
// Stable: yes
//
// For more information, documentation, source code, examples, support, links,
// etc. please see https://www.tawesoft.co.uk/go and 
// https://www.tawesoft.co.uk/go/glcaps
package glcaps // import "tawesoft.co.uk/go/glcaps"

// Code generated by internal. DO NOT EDIT.
// Instead, edit DESC.txt and run mkdocs.sh.