package main

import (
	"github.com/ajhager/eng"
)

var (
	batch   *eng.Batch
	camera  *eng.Camera
	regions []*eng.Region
	time    float32
)

type Demo struct {
	*eng.Game
}

func (d *Demo) Open() {
	eng.SetBgColor(eng.NewColorBytesA(20, 19, 22))
	texture := eng.NewTexture("data/test.png")
	regions = texture.Split(32, 32)
	batch = eng.NewBatch()
	camera = eng.NewCamera(480, 320)
}

func (d *Demo) Update(dt float32) {
	if eng.KeyPressed(eng.Esc) {
		eng.Exit()
	}
	time += dt * 4
	if time > camera.ViewportWidth/100*1.1 {
		time = 0
	}
}

func (d *Demo) Draw() {
	camera.Update()
	batch.SetProjection(camera.Combined)
	batch.Begin()
	batch.Draw(regions[3], time*100-camera.ViewportWidth/100, camera.Position.Y, 4, 30, 2, 2, -time*100, nil)
	batch.Draw(regions[2], camera.Position.X-16, camera.Position.Y-16, 16, 16, 1, 1, 0, nil)
	batch.End()
}

func (d *Demo) Resize(w, h int) {
	batch.Resize(w, h)
}

func main() {
	eng.Run(&Demo{})
}
