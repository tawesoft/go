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

// Red is an in-memory image whose At method returns color.Red values.
type Red struct {
	// Pix holds the image's pixels, as Red values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *Red) ColorModel() color.Model { return xcolor.RedModel }

func (p *Red) Bounds() image.Rectangle { return p.Rect }

func (p *Red) At(x, y int) color.Color {
	return p.RedAt(x, y)
}

func (p *Red) RedAt(x, y int) xcolor.Red {
	if !(image.Point{x, y}.In(p.Rect)) {
		return xcolor.Red{}
	}
	i := p.PixOffset(x, y)
	return xcolor.Red{R: p.Pix[i]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Red) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

func (p *Red) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = xcolor.RedModel.Convert(c).(color.RGBA).R
}

func (p *Red) SetRed(x, y int, c xcolor.Red) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = c.R
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Red) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Red{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Red{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *Red) Opaque() bool {
	return true
}

// NewRed returns a new Red image with the given bounds.
func NewRed(r image.Rectangle) *Red {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Red{pix, 1 * w, r}
}
