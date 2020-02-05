// Based on https://golang.org/src/image/color/color.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xcolor

import (
    "image/color"
)

// Red represents an 8-bit red color.
type Red struct {
	R uint8
}

func (c Red) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	return r, 0, 0, 0xFFFF
}


var RedModel color.Model = color.ModelFunc(redModel)

func redModel(c color.Color) color.Color {
	if _, ok := c.(color.RGBA); ok {
		return c
	}
	r, _, _, _ := c.RGBA()
	return Red{uint8(r >> 8)}
}
