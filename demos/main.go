package main

import (
	"github.com/ajhager/eng"
)

type Demo struct {
	*eng.Game
	regions []*eng.Region
	index   int
}

func (d *Demo) Open() {
	texture := eng.NewTexture("test.png")
	d.regions = texture.Split(32, 32)
	d.index = 0
}

func (d *Demo) Update(dt float64) {
	if eng.KeyPressed(eng.Esc) {
		eng.Exit()
	}
}

func (d *Demo) MouseDown(x, y, b int) {
	d.index += 1
	if d.index >= len(d.regions) {
		d.index = 0
	}
}

func (d *Demo) Draw() {
	eng.SetColor(eng.White)
	eng.Draw(d.regions[d.index], float32(eng.MouseX()), float32(eng.MouseY()))
	eng.SetColor(eng.DesaturatedSky)
	eng.Print("Hello, world!", float32(eng.Width())/2-6.5*16, float32(eng.Height())/2)
}

func main() {
	eng.Run(&Demo{})
}
