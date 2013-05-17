package main

import (
	"github.com/ajhager/eng"
)

var (
	batch   *eng.Batch
	regions map[string]*eng.Region
)

type Game struct {
	*eng.Game
}

func (g *Game) Open() {
	batch = eng.NewBatch()
	texture := eng.NewTexture("data/spineboy.png")
	regions = texture.Unpack("data/spineboy.json")
}

func (g *Game) Draw() {
	batch.Begin()
	batch.Draw(regions["head"], 480, 300, 32, 32, 1, 1, 0, nil)
	batch.End()
}

func main() {
	eng.Run("Atlas", 1024, 640, false, new(Game))
}
