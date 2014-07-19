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
	eng.SetBgColor(eng.NewColorHex(0x38345c))
	texture := eng.NewTexture(eng.Files.Image("spineboy"))
	regions = texture.Unpack(eng.Files.Json("spineboy"))
}

func (g *Game) Draw() {
	batch.Begin()
	batch.Draw(regions["head"], 680, 200, 0.5, 0.5, 1, 1, 0)
	batch.Draw(regions["left-foot"], 930, 515, 0.5, 0.5, 1, 1, 15)
	batch.Draw(regions["eyes"], 540, 400, 0.5, 0.5, 1, 1, 0)
	batch.End()
}

func main() {
	eng.Run("Atlas", 1024, 640, true, new(Game))
}
