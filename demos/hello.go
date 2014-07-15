package main

import (
	"github.com/ajhager/eng"
)

var (
	batch *eng.Batch
)

type Game struct {
	*eng.Game
}

func (g *Game) Setup() {
	eng.SetBgColor(eng.NewColor(0.5, 0.75, 0.5))
	batch = eng.NewBatch()
}

func (g *Game) Draw() {
	batch.Begin()
	eng.DefaultFont().Print(batch, "Hello, world!", 0, 0)
	batch.End()
}

func main() {
	eng.Run("Hello", 1024, 640, false, new(Game))
}
