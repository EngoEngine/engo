package main

import (
	"github.com/ajhager/eng"
	gl "github.com/chsc/gogl/gl33"
)

var (
	batch  *eng.Batch
	canvas *eng.Canvas
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
	gl.ClearColor(.8, .1, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	eng.DefaultFont().Print(batch, "canvas", x, y, nil)
	batch.End()
	canvas.End()

	region := canvas.Region()

	batch.Begin()
	batch.Draw(region, -200, 0, 512, 320, .5, .5, 0, nil)
	batch.Draw(region, 100, 200, 512, 320, .5, .5, 0, nil)
	batch.Draw(region, 200, -100, 512, 320, .5, .5, 0, nil)
	batch.End()
}

func main() {
	eng.Run("Canvas", 1024, 640, false, new(Game))
}
