package main

import (
	"github.com/ajhager/eng"
)

var img = [][]byte{
	{8, 8, 8, 8, 8, 8, 8, 8},
	{8, 6, 6, 6, 6, 6, 6, 8},
	{8, 1, 1, 1, 1, 1, 1, 8},
	{8, 1, 1, 1, 1, 1, 1, 8},
	{8, 1, 1, 1, 1, 1, 1, 8},
	{8, 1, 1, 8, 8, 1, 1, 8},
	{8, 1, 1, 8, 8, 1, 1, 8},
	{8, 8, 8, 8, 8, 8, 8, 8},
}

type Demo struct {
	*eng.Game
	regions []*eng.Region
}

func (d *Demo) Open() {
	eng.SetBgColor(eng.NewColorBytesA(153, 119, 119))
	texture := eng.NewTexture("pal.png")
	d.regions = texture.Split(1, 1)
}

func (d *Demo) Update(dt float64) {
	if eng.KeyPressed(eng.Esc) {
		eng.Exit()
	}
}

func (d *Demo) Draw() {
	s := float32(8)
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			eng.Draw(d.regions[img[7-y][x]], float32(10+x)*s, float32(10+y)*s, 0, 0, s, s, 0, nil)
		}
	}
}

func main() {
	eng.Run(&Demo{})
}
