package eng

import (
	"math/rand"
)

type Color struct {
	R, G, B, A float32
}

func NewColor(r, g, b, a float32) *Color {
	return &Color{r, g, b, a}
}

func NewColorBytes(r, g, b, a byte) *Color {
	color := new(Color)
	color.R = float32(r) / 255.0
	color.G = float32(g) / 255.0
	color.B = float32(b) / 255.0
	color.A = float32(a) / 255.0
	return color
}

func NewColorBytesA(r, g, b byte) *Color {
	color := new(Color)
	color.R = float32(r) / 255.0
	color.G = float32(g) / 255.0
	color.B = float32(b) / 255.0
	color.A = float32(1)
	return color
}

func NewColorRand() *Color {
	return &Color{rand.Float32(), rand.Float32(), rand.Float32(), 1}
}

func NewColorRandA() *Color {
	return &Color{rand.Float32(), rand.Float32(), rand.Float32(), rand.Float32()}
}
