package main

import (
	"github.com/ajhager/eng"
)

var (
	batch   *eng.Batch
	regions []*eng.Region
	time    float32
	font    *eng.Font
)

type Demo struct {
	*eng.Game
}

func (d *Demo) Open() {
	eng.SetBgColor(eng.NewColorBytesA(20, 19, 22))
	texture := eng.NewTexture("test.png")
	regions = texture.Split(32, 32)
	batch = eng.NewBatch()
	font = eng.DefaultFont
}

func (d *Demo) Update(dt float32) {
	if eng.KeyPressed(eng.Esc) {
		eng.Exit()
	}
	time += dt * 4
	if time > 12 {
		time = 0
	}
}

func (d *Demo) Draw() {
	batch.Begin()
	batch.Draw(regions[2], time*100-100, 300, 4, 30, 3, 3, -time*50, nil)
	font.Print(batch, "I'M AN APPLE!!!", 400, 150, eng.LightChartreuse)
	batch.End()
}

func main() {
	eng.Run(&Demo{})
}
