/*
Package xcolor implements Red, RedGreen, and RGB color models matching the core
image.color interface.

Note that there are good reasons these color types aren't in the core
image.color package. The native color types may have optimized fast-paths
for many use cases.

This package is a tradeoff of these optimizations against lower memory
usage. This package is intended to be used in computer graphics (e.g.
OpenGL) where images are uploaded to the GPU in a specific format (such as
GL_R, GL_RG, or GL_RGB) and we don't care about the performance of native
Go image manipulation.

Home page https://tawesoft.co.uk/go

For source code see https://github.com/tawesoft/go/tree/master/ximage/xcolor

For documentation see https://godoc.org/tawesoft.co.uk/go/ximage/xcolor

See also https://godoc.org/tawesoft.co.uk/go/ximage
*/
package xcolor // import "tawesoft.co.uk/go/ximage/xcolor"
