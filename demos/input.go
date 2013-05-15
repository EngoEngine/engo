package main

import (
	"fmt"
	"github.com/ajhager/eng"
	"math"
	"math/rand"
)

var (
	batch   *eng.Batch
	mx, my  int
	mz      float32
	color   *eng.Color
	letters string
)

type Game struct {
	*eng.Game
}

func (g *Game) Init(config *eng.Config) {
	config.Title = "Hello"
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	color = eng.NewColor(1, 1, 1, 1)
}

func (g *Game) Update(dt float32) {
	if math.Abs(float64(mz)) > .1 {
		mz -= float32(math.Copysign(float64(dt)*100, float64(mz)))
	}
}

func (g *Game) Draw() {
	batch.Begin()
	eng.DefaultFont().Print(batch, fmt.Sprintf("%v %v", mx, my), float32(mx-48), float32(my-16)+mz, color)
	eng.DefaultFont().Print(batch, letters, 0, 320, nil)
	batch.End()
}

func (g *Game) MouseMove(x, y int) {
	mx = x
	my = y
}

func (g *Game) MouseDown(x, y, b int) {
	switch b {
	default:
	case eng.MouseLeft:
		color.R = .25
	case eng.MouseRight:
		color.G = .25
	case eng.MouseMiddle:
		color.B = .25
	}
}

func (g *Game) MouseUp(x, y, b int) {
	color.R = 1
	color.G = 1
	color.B = 1
	color.A = 1
}

func (g *Game) MouseScroll(x, y, amount int) {
	mz += float32(amount) * 3
}

func (g *Game) KeyType(k int) {
	letters = letters + string(rune(k))
}

func (g *Game) KeyDown(k int) {
	if k == eng.Space {
		eng.SetBgColor(eng.NewColor(rand.Float32(), rand.Float32(), rand.Float32(), 1))
	}
}

func (g *Game) KeyUp(k int) {
	if k == eng.Esc {
		eng.Exit()
	}
}

func main() {
	eng.Run(new(Game))
}
