package main

import (
	"github.com/ajhager/eng"
)

var (
	batch *eng.Batch
	ttf   *eng.Font
	color = eng.NewColorBytes(200, 128, 64)
)

type Game struct {
	*eng.Game
}

func (g *Game) Setup() {
	batch = eng.NewBatch()
	ttf = eng.NewTrueTypeFont("data/DroidSansMono.ttf", 75, " !\"#$%&'()*+,-./0123456789:;█<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
}

func (g *Game) Draw() {
	batch.Begin()
	ttf.Print(batch, "Hello, True Type Fonts!", 0, float32(eng.Height()-150)/2, color)
	batch.End()
}

func main() {
	eng.Run("ttf", 1024, 640, false, new(Game))
}
