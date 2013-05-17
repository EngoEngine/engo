package main

import (
	"github.com/ajhager/eng"
)

var batch *eng.Batch

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
}

func (g *Game) Draw() {
	batch.Begin()
	eng.DefaultFont().Print(batch, "Hello, world!", 430, 280, nil)
	batch.End()
}

func main() {
	eng.Run("Hello", 1024, 640, false, new(Game))
}
