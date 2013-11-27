// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"math/rand"
)

type Color struct {
	R, G, B, A float32
}

// NewColor constructs a color using 32bit floating point values in
// the range 0.0 to 1.0.
func NewColor(r, g, b, a float32) *Color {
	return &Color{r, g, b, a}
}

// NewColorBytes constructs a color using 8bit integers in the range
// 0 to 255.
func NewColorBytes(r, g, b, a byte) *Color {
	color := new(Color)
	color.R = float32(r) / 255.0
	color.G = float32(g) / 255.0
	color.B = float32(b) / 255.0
	color.A = float32(a) / 255.0
	return color
}

// NewColorHex contructs a color from a uint32, ie. 0xFFFFFF.
func NewColorHex(n uint32) *Color {
	return NewColorBytes(uint8((n>>16)&0xFF), uint8((n>>8)&0xFF), uint8(n&0xFF), 255)
}

// NewColorRand constructs a random color.
func NewColorRand() *Color {
	return &Color{rand.Float32(), rand.Float32(), rand.Float32(), 1}
}
