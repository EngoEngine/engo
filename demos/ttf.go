package main

import (
	"github.com/ajhager/eng"
)

var (
	batch *eng.Batch
	ttf   *eng.Font
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	ttf = eng.NewTrueTypeFont("data/DroidSansMono.ttf", 75, " !\"#$%&'()*+,-./0123456789:;█<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
}

func (g *Game) Draw() {
	batch.Begin()
	ttf.Print(batch, "Hello, True Type Fonts!", 0, float32(eng.Height()-75)/2, eng.DarkSky)
	batch.End()
}

func main() {
	eng.Run("Hello", 1024, 640, false, new(Game))
}
