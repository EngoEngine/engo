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

func (g *Game) Load() {
	eng.Files.Add("spineboy", "data/spineboy.png")
	eng.Files.Add("spineboy", "data/spineboy.json")
}

func (g *Game) Setup() {
	batch = eng.NewBatch()
	texture := eng.NewTexture(eng.Files.Image("spineboy"))
	regions = texture.Unpack(eng.Files.Json("spineboy"))
}

func (g *Game) Draw() {
	batch.Begin()
	batch.Draw(regions["head"], 480, 300, 32, 32, 1, 1, 0)
	batch.End()
}

func main() {
	eng.Run("Atlas", 1024, 640, false, new(Game))
}
