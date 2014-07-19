package main

import (
	"github.com/ajhager/eng"
)

var (
	batch *eng.Batch
	font  *eng.Font
	blue  = eng.NewColorHex(0x33b5e5)
	grey  = eng.NewColorHex(0xf4f4f4)
)

type Game struct {
	*eng.Game
}

func (g *Game) Load() {
	eng.Files.Add("font", "data/font.png")
}

func (g *Game) Setup() {
	font = eng.NewGridFont(eng.Files.Image("font"), 20, 20, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
	eng.SetBgColor(blue)
	batch = eng.NewBatch()
}

func (g *Game) Draw() {
	batch.Begin()
	batch.SetColor(grey)
	font.Print(batch, "hello world", eng.Width()/2-120, eng.Height()/2-20)
	batch.End()
}

func main() {
	eng.Run("Hello", 1024, 640, true, new(Game))
}
