// Based on https://golang.org/src/image/color/color.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xcolor

import (
    "image/color"
)

// Red represents an 8-bit Red 8-bit Green color.
type RG struct {
    R uint8
    G uint8
}

func (c RG) RGBA() (r, g, b, a uint32) {
    r = uint32(c.R)
    r |= r << 8
    g = uint32(c.G)
    g |= g << 8
    return r, g, 0, 0xFFFF
}

var RGModel color.Model = color.ModelFunc(func (c color.Color) color.Color {
    if _, ok := c.(color.RGBA); ok {
        return c
    }
    r, g, _, _ := c.RGBA()
    return RG{uint8(r >> 8), uint8(g >> 8)}
})
