package main

import (
	"fmt"
	"github.com/ajhager/eng"
	"math"
)

var (
	batch      *eng.Batch
	mx, my, mz float32
	color      *eng.Color
	letters    string
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	color = eng.NewColor(1, 1, 1)
}

func (g *Game) Update(dt float32) {
	if math.Abs(float64(mz)) > .1 {
		mz -= float32(math.Copysign(float64(dt)*100, float64(mz)))
	}
}

func (g *Game) Draw() {
	batch.Begin()
	eng.DefaultFont().Print(batch, fmt.Sprintf("%.0f %.0f", mx, my), mx-48, my-16+mz, color)
	eng.DefaultFont().Print(batch, letters, 0, 320, nil)
	batch.End()
}

func (g *Game) MouseMove(x, y float32) {
	mx = x
	my = y
}

func (g *Game) MouseDown(x, y float32, b eng.MouseButton) {
	switch b {
	default:
	case eng.MouseButtonLeft:
		color.R = .25
	case eng.MouseButtonRight:
		color.G = .25
	case eng.MouseButtonMiddle:
		color.B = .25
	}
}

func (g *Game) MouseUp(x, y float32, b eng.MouseButton) {
	color.R = 1
	color.G = 1
	color.B = 1
	color.A = 1
}

func (g *Game) MouseScroll(x, y, amount float32) {
	mz += float32(amount) * 3
}

func (g *Game) KeyType(k rune) {
	letters = letters + string(k)
}

func (g *Game) KeyDown(k eng.Key) {
	if k == eng.Space {
		eng.SetBgColor(eng.NewColorRand())
	}
}

func (g *Game) KeyUp(k eng.Key) {
	if k == eng.Escape {
		eng.Exit()
	}
}

func main() {
	eng.Run("Input", 1024, 640, false, new(Game))
}
