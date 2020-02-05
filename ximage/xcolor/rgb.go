// Based on https://golang.org/src/image/color/color.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xcolor

import (
    "image/color"
)

// RGB represents an 8-bit Red 8-bit Green 8-bit Blue color.
type RGB struct {
    R uint8
    G uint8
    B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
    r |= r << 8
    g = uint32(c.G)
    g |= g << 8
    b = uint32(c.B)
    b |= b << 8
	return r, g, b, 0xFFFF
}

var RGBModel color.Model = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
	if _, ok := c.(color.RGBA); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}
