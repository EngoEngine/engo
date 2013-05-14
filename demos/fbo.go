package main

import (
	"github.com/ajhager/eng"
	gl "github.com/chsc/gogl/gl33"
)

var (
	batch  *eng.Batch
	canvas *eng.Canvas
	region *eng.Region
)

type Hello struct {
	*eng.Game
}

func (g *Hello) Open() {
	batch = eng.NewBatch()
	canvas = eng.NewCanvas(eng.Width(), eng.Height())
	region = eng.NewRegion(canvas.Texture(), 0, 0, eng.Width(), eng.Height())
	region.Flip(false, true)
}

func (g *Hello) Draw() {
	canvas.Begin()
	batch.Begin()
	gl.ClearColor(.8, .1, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	eng.DefaultFont().Print(batch, "Hello, world!", 430, 280, nil)
	batch.End()
	canvas.End()

	batch.Begin()
	batch.Draw(region, 0, 0, 512, 320, .5, .5, 0, nil)
	batch.End()
}

func main() {
	eng.Run(new(Hello))
}
