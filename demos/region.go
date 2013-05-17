package main

import (
	"github.com/ajhager/eng"
)

var (
	batch   *eng.Batch
	regions []*eng.Region
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	texture := eng.NewTexture("data/bot.png")
	regions = texture.Split(64, 64)
}

func (g *Game) Draw() {
	batch.Begin()
	batch.Draw(regions[0], 480, 300, 32, 32, 1, 1, 0, nil)
	batch.End()
}

func main() {
	eng.Run("Region", 1024, 640, false, new(Game))
}
