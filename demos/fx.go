package main

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/eng/fx"
)

var batch *eng.Batch
var effect eng.Effect
var region *eng.Region

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	effect = fx.NewFilm(.3, .2, 4096, false)
	texture := eng.NewTexture("data/bot.png")
	batch.SetShader(effect.Shader())
	region = texture.Split(64, 64)[0]
}

func (g *Game) Draw() {
	batch.Begin()
	effect.Setup()
	batch.Draw(region, 200, 0, 0, 0, 10, 10, 0, nil)
	batch.End()
}

func main() {
	eng.Run("Fx", 1024, 640, false, new(Game))
}
