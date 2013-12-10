package main

import (
	"github.com/ajhager/eng"
)

var (
	batch  *eng.Batch
	canvas *eng.Canvas
	red    = eng.NewColorBytesA(255, 128, 0, 128)
	green  = eng.NewColorBytesA(128, 255, 0, 128)
	blue   = eng.NewColorBytesA(128, 0, 255, 128)
	white  = eng.NewColor(1, 1, 1)
	black  = eng.NewColor(0, 0, 0)
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	canvas = eng.NewCanvas(eng.Width(), eng.Height())
}

func (g *Game) Draw() {
	x := float32(canvas.Width()/2 - 50)
	y := float32(canvas.Height() / 2)

	canvas.Begin()
	batch.Begin()
	eng.Clear(white)
	eng.DefaultFont().Print(batch, "canvas", x, y, black)
	batch.End()
	canvas.End()

	region := canvas.Region()

	batch.Begin()
	batch.Draw(region, -200, 0, 512, 320, .5, .5, 0, blue)
	batch.Draw(region, 100, 200, 512, 320, .5, .5, 0, red)
	batch.Draw(region, 200, -100, 512, 320, .5, .5, 0, green)
	batch.End()
}

func main() {
	eng.Run("Canvas", 1024, 640, false, new(Game))
}
