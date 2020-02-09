// Based on https://golang.org/src/image/image.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ximage

import (
    "image"
    "image/color"
    "tawesoft.co.uk/go/ximage/xcolor"
)

// RG is an in-memory image whose At method returns color.RG values.
type RG struct {
	// Pix holds the image's pixels, as Red values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *RG) ColorModel() color.Model { return xcolor.RGModel }

func (p *RG) Bounds() image.Rectangle { return p.Rect }

func (p *RG) At(x, y int) color.Color {
	return p.RGAt(x, y)
}

func (p *RG) RGAt(x, y int) xcolor.RG {
	if !(image.Point{x, y}.In(p.Rect)) {
		return xcolor.RG{}
	}
	i := p.PixOffset(x, y)
	return xcolor.RG{R: p.Pix[i], G: p.Pix[i+1]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RG) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

func (p *RG) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
    i := p.PixOffset(x, y)
    rgba := xcolor.RGModel.Convert(c).(color.RGBA)
    p.Pix[i] = rgba.R
    p.Pix[i+1] = rgba.G
}

func (p *RG) SetRG(x, y int, c xcolor.RG) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
    p.Pix[i] = c.R
    p.Pix[i+1] = c.G
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RG) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RG{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RG{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RG) Opaque() bool {
	return true
}

// NewRG returns a new RG image with the given bounds.
func NewRG(r image.Rectangle) *RG {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 2*w*h)
	return &RG{pix, 2 * w, r}
}
