package main

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/eng/scene"
)

var (
	stage  *scene.Stage
	region *eng.Region
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	stage = scene.NewStage(1280, 800, true)
	stage.Camera().Position.X += stage.GutterWidth()
	stage.Camera().Position.Y -= stage.GutterHeight()

	region = eng.NewRegion(eng.BlankTexture(), 0, 0, 1, 1)
}

func (g *Game) Update(dt float32) {
	stage.Update()
}

func (g *Game) Draw() {
	batch := stage.Batch()
	batch.Begin()
	batch.Draw(region, 0, 0, 0, 0, 1280, 800, 0, eng.DarkerSea)
	batch.Draw(region, 640, 400, .5, .5, 960, 640, 0, eng.DarkSea)
	batch.End()
}

func main() {
	eng.Run("stage", 960, 640, false, new(Game))
}
