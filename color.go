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

// Add = old + new
func (c *Color) Add(o *Color) *Color {
	return add(o, c)
}

// AddAlpha = old + alpha*new
func (c *Color) AddAlpha(o *Color, a float32) *Color {
	return addAlpha(o, c, a)
}

// Alpha = (1-alpha)*old + alpha*(new-old)
func (c *Color) Alpha(o *Color, a float32) *Color {
	return alpha(o, c, a)
}

// Blend satisfies the Blender interface by returning a constant copy.
func (c *Color) Blend(o *Color, i, t int) *Color {
	return c.Copy()
}

// Burn = old + new - white
func (c *Color) Burn(o *Color) *Color {
	return burn(o, c)
}

// Dodge = new / (white - old)
func (c *Color) Dodge(o *Color) *Color {
	return dodge(o, c)
}

// Multiply = old * new
func (c *Color) Multiply(o *Color) *Color {
	return multiply(o, c)
}

// Overlay = new.x <= 0.5 ? 2*new*old : white - 2*(white-new)*(white-old)
func (c *Color) Overlay(o *Color) *Color {
	return overlay(o, c)
}

// Screen = white - (white - old) * (white - new)
func (c *Color) Screen(o *Color) *Color {
	return screen(o, c)
}

// Darken = MIN(old, new)
func (c *Color) Darken(o *Color) *Color {
	return darken(o, c)
}

// Lighten = MIN(old, new)
func (c *Color) Lighten(o *Color) *Color {
	return lighten(o, c)
}

// Scale = old * s
func (c *Color) Scale(s float32) *Color {
	return scale(c, s)
}

// RandScale scales the color a random amount.
func (c *Color) RandScale() *Color {
	return scale(c, rand.Float32())
}

// BlendFunc is a function that takes a color and returns a new color.
type BlendFunc func(*Color) *Color

// BlendFunc satisfies the Blender interface by calling itself with
// the passed in color.
func (bf BlendFunc) Blend(o *Color, i, t int) *Color {
	return bf(o)
}

// BlendAdd
func BlendAdd(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return add(top, bot)
	}
}

// BlendAddAlpha
func BlendAddAlpha(top *Color, a float32) BlendFunc {
	return func(bot *Color) *Color {
		return addAlpha(top, bot, a)
	}
}

// BlendAlpha
func BlendAlpha(top *Color, a float32) BlendFunc {
	return func(bot *Color) *Color {
		return addAlpha(top, bot, a)
	}
}

// BlendBurn
func BlendBurn(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return burn(top, bot)
	}
}

// BlendDarken
func BlendDarken(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return darken(top, bot)
	}
}

// BlendDodge
func BlendDodge(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return dodge(top, bot)
	}
}

// BlendLighten
func BlendLighten(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return lighten(top, bot)
	}
}

// BlendMultiply
func BlendMultiply(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return multiply(top, bot)
	}
}

// BlendOverlay
func BlendOverlay(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return overlay(top, bot)
	}
}

// BlendScreen
func BlendScreen(top *Color) BlendFunc {
	return func(bot *Color) *Color {
		return screen(top, bot)
	}
}

// BlendRandScale
func BlendRandScale() BlendFunc {
	return func(bot *Color) *Color {
		return bot.RandScale()
	}
}

// BlendScale
func BlendScale(s float32) BlendFunc {
	return func(bot *Color) *Color {
		return scale(bot, s)
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

func overlayF(top, bot float32) float32 {
	if bot < 0.5 {
		return 2 * top * bot
	} else {
		return 1 - 2*(1-top)*(1-bot)
	}
}

// Color blending functions
func add(top, bot *Color) *Color {
	return NewColor(
		clampF(0, 1, bot.R+top.R),
		clampF(0, 1, bot.G+top.G),
		clampF(0, 1, bot.B+top.B),
		clampF(0, 1, bot.A+top.A),
	)
}

func addAlpha(top, bot *Color, a float32) *Color {
	return NewColor(
		clampF(0, 1, bot.R*a+top.R),
		clampF(0, 1, bot.G*a+top.G),
		clampF(0, 1, bot.B*a+top.B),
		clampF(0, 1, bot.A*a+top.A),
	)
}

func alpha(top, bot *Color, a float32) *Color {
	a = clampF(0, 1, a)
	return NewColor(bot.R+(top.R-bot.R)*a, bot.G+(top.G-bot.G)*a, bot.B+(top.B-bot.B)*a, bot.A+(top.A-bot.A)*a)
}

func burn(top, bot *Color) *Color {
	return NewColor(
		clampF(0, 1, bot.R+top.R-1),
		clampF(0, 1, bot.G+top.G-1),
		clampF(0, 1, bot.B+top.B-1),
		clampF(0, 1, bot.A+top.A-1),
	)
}

func darken(top, bot *Color) *Color {
	return NewColor(
		float32(math.Min(float64(top.R), float64(bot.R))),
		float32(math.Min(float64(top.G), float64(bot.G))),
		float32(math.Min(float64(top.B), float64(bot.B))),
		float32(math.Min(float64(top.A), float64(bot.A))),
	)
}

func dodge(top, bot *Color) *Color {
	return NewColor(dodgeF(top.R, bot.R), dodgeF(top.G, bot.G), dodgeF(top.B, bot.B), dodgeF(top.A, bot.A))
}

func lighten(top, bot *Color) *Color {
	return NewColor(
		float32(math.Max(float64(top.R), float64(bot.R))),
		float32(math.Max(float64(top.G), float64(bot.G))),
		float32(math.Max(float64(top.B), float64(bot.B))),
		float32(math.Max(float64(top.A), float64(bot.A))),
	)
}

func multiply(top, bot *Color) *Color {
	return NewColor(top.R*bot.R, top.G*bot.G, top.B*bot.B, top.A*bot.A)
}

func overlay(top, bot *Color) *Color {
	return NewColor(overlayF(top.R, bot.R), overlayF(top.G, bot.G), overlayF(top.B, bot.B), overlayF(top.A, bot.A))
}

func scale(bot *Color, s float32) *Color {
	return NewColor(bot.R*s, bot.G*s, bot.B*s, bot.A*s)
}

func screen(top, bot *Color) *Color {
	return NewColor(1-(1-top.R)*(1-bot.R), 1-(1-top.G)*(1-bot.G), 1-(1-top.B)*(1-bot.B), 1-(1-top.A)*(1-bot.A))
}

// Predefined Colors
var (
	Black        = NewColorBytes(0, 0, 0, 255)
	DarkestGrey  = NewColorBytes(31, 31, 31, 255)
	DarkerGrey   = NewColorBytes(63, 63, 63, 255)
	DarkGrey     = NewColorBytes(95, 95, 95, 255)
	Grey         = NewColorBytes(127, 127, 127, 255)
	LightGrey    = NewColorBytes(159, 159, 159, 255)
	LighterGrey  = NewColorBytes(191, 191, 191, 255)
	LightestGrey = NewColorBytes(223, 223, 223, 255)
	White        = NewColorBytes(255, 255, 255, 255)

	DarkestSepia  = NewColorBytes(31, 24, 15, 255)
	DarkerSepia   = NewColorBytes(63, 50, 31, 255)
	DarkSepia     = NewColorBytes(94, 75, 47, 255)
	Sepia         = NewColorBytes(127, 101, 63, 255)
	LightSepia    = NewColorBytes(158, 134, 100, 255)
	LighterSepia  = NewColorBytes(191, 171, 143, 255)
	LightestSepia = NewColorBytes(222, 211, 195, 255)

	DesaturatedRed        = NewColorBytes(127, 63, 63, 255)
	DesaturatedFlame      = NewColorBytes(127, 79, 63, 255)
	DesaturatedOrange     = NewColorBytes(127, 95, 63, 255)
	DesaturatedAmber      = NewColorBytes(127, 111, 63, 255)
	DesaturatedYellow     = NewColorBytes(127, 127, 63, 255)
	DesaturatedLime       = NewColorBytes(111, 127, 63, 255)
	DesaturatedChartreuse = NewColorBytes(95, 127, 63, 255)
	DesaturatedGreen      = NewColorBytes(63, 127, 63, 255)
	DesaturatedSea        = NewColorBytes(63, 127, 95, 255)
	DesaturatedTurquoise  = NewColorBytes(63, 127, 111, 255)
	DesaturatedCyan       = NewColorBytes(63, 127, 127, 255)
	DesaturatedSky        = NewColorBytes(63, 111, 127, 255)
	DesaturatedAzure      = NewColorBytes(63, 95, 127, 255)
	DesaturatedBlue       = NewColorBytes(63, 63, 127, 255)
	DesaturatedHan        = NewColorBytes(79, 63, 127, 255)
	DesaturatedViolet     = NewColorBytes(95, 63, 127, 255)
	DesaturatedPurple     = NewColorBytes(111, 63, 127, 255)
	DesaturatedFuchsia    = NewColorBytes(127, 63, 127, 255)
	DesaturatedMagenta    = NewColorBytes(127, 63, 111, 255)
	DesaturatedPink       = NewColorBytes(127, 63, 95, 255)
	DesaturatedCrimson    = NewColorBytes(127, 63, 79, 255)

	LightestRed        = NewColorBytes(255, 191, 191, 255)
	LightestFlame      = NewColorBytes(255, 207, 191, 255)
	LightestOrange     = NewColorBytes(255, 223, 191, 255)
	LightestAmber      = NewColorBytes(255, 239, 191, 255)
	LightestYellow     = NewColorBytes(255, 255, 191, 255)
	LightestLime       = NewColorBytes(239, 255, 191, 255)
	LightestChartreuse = NewColorBytes(223, 255, 191, 255)
	LightestGreen      = NewColorBytes(191, 255, 191, 255)
	LightestSea        = NewColorBytes(191, 255, 223, 255)
	LightestTurquoise  = NewColorBytes(191, 255, 239, 255)
	LightestCyan       = NewColorBytes(191, 255, 255, 255)
	LightestSky        = NewColorBytes(191, 239, 255, 255)
	LightestAzure      = NewColorBytes(191, 223, 255, 255)
	LightestBlue       = NewColorBytes(191, 191, 255, 255)
	LightestHan        = NewColorBytes(207, 191, 255, 255)
	LightestViolet     = NewColorBytes(223, 191, 255, 255)
	LightestPurple     = NewColorBytes(239, 191, 255, 255)
	LightestFuchsia    = NewColorBytes(255, 191, 255, 255)
	LightestMagenta    = NewColorBytes(255, 191, 239, 255)
	LightestPink       = NewColorBytes(255, 191, 223, 255)
	LightestCrimson    = NewColorBytes(255, 191, 207, 255)

	LighterRed        = NewColorBytes(255, 127, 127, 255)
	LighterFlame      = NewColorBytes(255, 159, 127, 255)
	LighterOrange     = NewColorBytes(255, 191, 127, 255)
	LighterAmber      = NewColorBytes(255, 223, 127, 255)
	LighterYellow     = NewColorBytes(255, 255, 127, 255)
	LighterLime       = NewColorBytes(223, 255, 127, 255)
	LighterChartreuse = NewColorBytes(191, 255, 127, 255)
	LighterGreen      = NewColorBytes(127, 255, 127, 255)
	LighterSea        = NewColorBytes(127, 255, 191, 255)
	LighterTurquoise  = NewColorBytes(127, 255, 223, 255)
	LighterCyan       = NewColorBytes(127, 255, 255, 255)
	LighterSky        = NewColorBytes(127, 223, 255, 255)
	LighterAzure      = NewColorBytes(127, 191, 255, 255)
	LighterBlue       = NewColorBytes(127, 127, 255, 255)
	LighterHan        = NewColorBytes(159, 127, 255, 255)
	LighterViolet     = NewColorBytes(191, 127, 255, 255)
	LighterPurple     = NewColorBytes(223, 127, 255, 255)
	LighterFuchsia    = NewColorBytes(255, 127, 255, 255)
	LighterMagenta    = NewColorBytes(255, 127, 223, 255)
	LighterPink       = NewColorBytes(255, 127, 191, 255)
	LighterCrimson    = NewColorBytes(255, 127, 159, 255)

	LightRed        = NewColorBytes(255, 63, 63, 255)
	LightFlame      = NewColorBytes(255, 111, 63, 255)
	LightOrange     = NewColorBytes(255, 159, 63, 255)
	LightAmber      = NewColorBytes(255, 207, 63, 255)
	LightYellow     = NewColorBytes(255, 255, 63, 255)
	LightLime       = NewColorBytes(207, 255, 63, 255)
	LightChartreuse = NewColorBytes(159, 255, 63, 255)
	LightGreen      = NewColorBytes(63, 255, 63, 255)
	LightSea        = NewColorBytes(63, 255, 159, 255)
	LightTurquoise  = NewColorBytes(63, 255, 207, 255)
	LightCyan       = NewColorBytes(63, 255, 255, 255)
	LightSky        = NewColorBytes(63, 207, 255, 255)
	LightAzure      = NewColorBytes(63, 159, 255, 255)
	LightBlue       = NewColorBytes(63, 63, 255, 255)
	LightHan        = NewColorBytes(111, 63, 255, 255)
	LightViolet     = NewColorBytes(159, 63, 255, 255)
	LightPurple     = NewColorBytes(207, 63, 255, 255)
	LightFuchsia    = NewColorBytes(255, 63, 255, 255)
	LightMagenta    = NewColorBytes(255, 63, 207, 255)
	LightPink       = NewColorBytes(255, 63, 159, 255)
	LightCrimson    = NewColorBytes(255, 63, 111, 255)

	Red        = NewColorBytes(255, 0, 0, 255)
	Flame      = NewColorBytes(255, 63, 0, 255)
	Orange     = NewColorBytes(255, 127, 0, 255)
	Amber      = NewColorBytes(255, 191, 0, 255)
	Yellow     = NewColorBytes(255, 255, 0, 255)
	Lime       = NewColorBytes(191, 255, 0, 255)
	Chartreuse = NewColorBytes(127, 255, 0, 255)
	Green      = NewColorBytes(0, 255, 0, 255)
	Sea        = NewColorBytes(0, 255, 127, 255)
	Turquoise  = NewColorBytes(0, 255, 191, 255)
	Cyan       = NewColorBytes(0, 255, 255, 255)
	Sky        = NewColorBytes(0, 191, 255, 255)
	Azure      = NewColorBytes(0, 127, 255, 255)
	Blue       = NewColorBytes(0, 0, 255, 255)
	Han        = NewColorBytes(63, 0, 255, 255)
	Violet     = NewColorBytes(127, 0, 255, 255)
	Purple     = NewColorBytes(191, 0, 255, 255)
	Fuchsia    = NewColorBytes(255, 0, 255, 255)
	Magenta    = NewColorBytes(255, 0, 191, 255)
	Pink       = NewColorBytes(255, 0, 127, 255)
	Crimson    = NewColorBytes(255, 0, 63, 255)

	DarkRed        = NewColorBytes(191, 0, 0, 255)
	DarkFlame      = NewColorBytes(191, 47, 0, 255)
	DarkOrange     = NewColorBytes(191, 95, 0, 255)
	DarkAmber      = NewColorBytes(191, 143, 0, 255)
	DarkYellow     = NewColorBytes(191, 191, 0, 255)
	DarkLime       = NewColorBytes(143, 191, 0, 255)
	DarkChartreuse = NewColorBytes(95, 191, 0, 255)
	DarkGreen      = NewColorBytes(0, 191, 0, 255)
	DarkSea        = NewColorBytes(0, 191, 95, 255)
	DarkTurquoise  = NewColorBytes(0, 191, 143, 255)
	DarkCyan       = NewColorBytes(0, 191, 191, 255)
	DarkSky        = NewColorBytes(0, 143, 191, 255)
	DarkAzure      = NewColorBytes(0, 95, 191, 255)
	DarkBlue       = NewColorBytes(0, 0, 191, 255)
	DarkHan        = NewColorBytes(47, 0, 191, 255)
	DarkViolet     = NewColorBytes(95, 0, 191, 255)
	DarkPurple     = NewColorBytes(143, 0, 191, 255)
	DarkFuchsia    = NewColorBytes(191, 0, 191, 255)
	DarkMagenta    = NewColorBytes(191, 0, 143, 255)
	DarkPink       = NewColorBytes(191, 0, 95, 255)
	DarkCrimson    = NewColorBytes(191, 0, 47, 255)

	DarkerRed        = NewColorBytes(127, 0, 0, 255)
	DarkerFlame      = NewColorBytes(127, 31, 0, 255)
	DarkerOrange     = NewColorBytes(127, 63, 0, 255)
	DarkerAmber      = NewColorBytes(127, 95, 0, 255)
	DarkerYellow     = NewColorBytes(127, 127, 0, 255)
	DarkerLime       = NewColorBytes(95, 127, 0, 255)
	DarkerChartreuse = NewColorBytes(63, 127, 0, 255)
	DarkerGreen      = NewColorBytes(0, 127, 0, 255)
	DarkerSea        = NewColorBytes(0, 127, 63, 255)
	DarkerTurquoise  = NewColorBytes(0, 127, 95, 255)
	DarkerCyan       = NewColorBytes(0, 127, 127, 255)
	DarkerSky        = NewColorBytes(0, 95, 127, 255)
	DarkerAzure      = NewColorBytes(0, 63, 127, 255)
	DarkerBlue       = NewColorBytes(0, 0, 127, 255)
	DarkerHan        = NewColorBytes(31, 0, 127, 255)
	DarkerViolet     = NewColorBytes(63, 0, 127, 255)
	DarkerPurple     = NewColorBytes(95, 0, 127, 255)
	DarkerFuchsia    = NewColorBytes(127, 0, 127, 255)
	DarkerMagenta    = NewColorBytes(127, 0, 95, 255)
	DarkerPink       = NewColorBytes(127, 0, 63, 255)
	DarkerCrimson    = NewColorBytes(127, 0, 31, 255)

	DarkestRed        = NewColorBytes(63, 0, 0, 255)
	DarkestFlame      = NewColorBytes(63, 15, 0, 255)
	DarkestOrange     = NewColorBytes(63, 31, 0, 255)
	DarkestAmber      = NewColorBytes(63, 47, 0, 255)
	DarkestYellow     = NewColorBytes(63, 63, 0, 255)
	DarkestLime       = NewColorBytes(47, 63, 0, 255)
	DarkestChartreuse = NewColorBytes(31, 63, 0, 255)
	DarkestGreen      = NewColorBytes(0, 63, 0, 255)
	DarkestSea        = NewColorBytes(0, 63, 31, 255)
	DarkestTurquoise  = NewColorBytes(0, 63, 47, 255)
	DarkestCyan       = NewColorBytes(0, 63, 63, 255)
	DarkestSky        = NewColorBytes(0, 47, 63, 255)
	DarkestAzure      = NewColorBytes(0, 31, 63, 255)
	DarkestBlue       = NewColorBytes(0, 0, 63, 255)
	DarkestHan        = NewColorBytes(15, 0, 63, 255)
	DarkestViolet     = NewColorBytes(31, 0, 63, 255)
	DarkestPurple     = NewColorBytes(47, 0, 63, 255)
	DarkestFuchsia    = NewColorBytes(63, 0, 63, 255)
	DarkestMagenta    = NewColorBytes(63, 0, 47, 255)
	DarkestPink       = NewColorBytes(63, 0, 31, 255)
	DarkestCrimson    = NewColorBytes(63, 0, 15, 255)

	Brass  = NewColorBytes(191, 151, 96, 255)
	Copper = NewColorBytes(197, 136, 124, 255)
	Gold   = NewColorBytes(229, 191, 0, 255)
	Silver = NewColorBytes(203, 203, 203, 255)

	Celadon = NewColorBytes(172, 255, 175, 255)
	Peach   = NewColorBytes(255, 159, 127, 255)
)
