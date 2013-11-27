// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"math"
	"math/rand"
)

// A type that satisfies the Blender interface takes a color,
// an index into a range of values, a maximum of the range,
// and returns a new interpolated color.
type Blender interface {
	Blend(*Color, int, int) *Color
}

// Color struct
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

// Copy returns a new color with the same components.
func (c *Color) Copy() *Color {
	return &Color{c.R, c.G, c.B, c.A}
}

// Blend satisfies the Blender interface by returning a constant copy.
func (c *Color) Blend(o *Color, i, t int) *Color {
	return c.Copy()
}

// Dodge = new / (white - old)
func (c *Color) Dodge(o *Color) *Color {
	return dodge(o, c)
}

// Multiply = old * new
func (c *Color) Multiply(o *Color) *Color {
	return multiply(o, c)
}

// BlendFunc is a function that takes a color and returns a new color.
type BlendFunc func(*Color) *Color

// BlendFunc satisfies the Blender interface by calling itself with
// the passed in color.
func (bf BlendFunc) Blend(o *Color, i, t int) *Color {
	return bf(o)
}

// BlendMultiply
func BlendMultiply(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return multiply(top, bot)
	}
}

// BlendDodge
func BlendDodge(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return dodge(top, bot)
	}
}

// ScaleFunc is a function that takes a color, an index into a range
// of values, a maximum of that range, and returns an interpolated
// new color.
type ScaleFunc func(*Color, int, int) *Color

// ScaleFunc satisfies the Blender interface by calling itself
// with the provided values.
func (sf ScaleFunc) Blend(c *Color, i, t int) *Color {
	return sf(c, i, t)
}

// DiscreteGradient
func DiscreteGradient(blenders ...Blender) ScaleFunc {
	return func(bot *Color, i, t int) *Color {
		return blenders[i%len(blenders)].Blend(bot, i, t)
	}
}

// LinearGradient
func LinearGradient(blenders ...Blender) ScaleFunc {
	return func(bot *Color, i, t int) *Color {
		if i == 0 {
			return blenders[0].Blend(bot, i, t)
		}

		if i == (t - 1) {
			return blenders[len(blenders)-1].Blend(bot, i, t)
		}

		a := (float32(i) / float32(t-1)) * float32(len(blenders)-1)
		b := int(math.Floor(float64(a)))
		return alpha(blenders[b+1].Blend(bot, i, t), blenders[b].Blend(bot, i, t), a-float32(b))
	}
}

// Float32 blending functions
func clampF(low, high, value float32) float32 {
	return float32(math.Min(float64(high), math.Max(float64(low), float64(value))))
}

func dodgeF(top, bot float32) float32 {
	if bot != 1 {
		return clampF(0, 1, top/(1-bot))
	}
	return 1
}

// Color blending functions
func alpha(top, bot *Color, a float32) *Color {
	a = clampF(0, 1, a)
	return NewColor(bot.R+(top.R-bot.R)*a, bot.G+(top.G-bot.G)*a, bot.B+(top.B-bot.B)*a, bot.A+(top.A-bot.A)*a)
}

func dodge(top, bot *Color) *Color {
	return NewColor(dodgeF(top.R, bot.R), dodgeF(top.G, bot.G), dodgeF(top.B, bot.B), dodgeF(top.A, bot.A))
}

func multiply(top, bot *Color) *Color {
	return NewColor(top.R*bot.R, top.G*bot.G, top.B*bot.B, top.A*bot.A)
}
