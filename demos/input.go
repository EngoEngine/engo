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

func (g *Game) Setup() {
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

func (g *Game) Mouse(x, y float32, a eng.Action) {
	switch a {
	case eng.MOVE:
		mx = x
		my = y
	case eng.PRESS:
		eng.SetBgColor(eng.NewColorRand())
	case eng.RELEASE:
		eng.SetBgColor(eng.NewColor(0, 0, 0))
	}
}

func (g *Game) Scroll(amount float32) {
	mz += amount
}

func (g *Game) Key(key eng.Key, mod eng.Modifier, act eng.Action) {
	switch act {
	case eng.RELEASE:
		if key == eng.Escape {
			eng.Exit()
		}
	}
}

func (g *Game) Type(char rune) {
	letters = letters + string(char)
}

func main() {
	eng.Run("Input", 1024, 640, false, new(Game))
}
