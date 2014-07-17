package main

import (
	"github.com/ajhager/eng"
)

var (
	batch *eng.Batch
	font  *eng.Font
)

type Game struct {
	*eng.Game
}

func (g *Game) Load() {
	eng.Files.Add("font", "data/font.png")
}

func (g *Game) Setup() {
	font = eng.NewGridFont(eng.Files.Image("font"), 20, 20, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
	eng.SetBgColor(eng.NewColor(0.5, 0.75, 0.5))
	batch = eng.NewBatch()
}

func (g *Game) Draw() {
	batch.Begin()
	font.Print(batch, "Hello, world!", eng.Width()/2-120, eng.Height()/2-20)
	batch.End()
}

func main() {
	eng.Run("Hello", 1024, 640, false, new(Game))
}
