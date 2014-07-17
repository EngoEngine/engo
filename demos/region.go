package main

import (
	"github.com/ajhager/eng"
)

var (
	batch   *eng.Batch
	regions []*eng.Region
	angle   float32
)

type Game struct {
	*eng.Game
}

func (g *Game) Load() {
	eng.Files.Add("bot", "data/bot.png")
}

func (g *Game) Setup() {
	eng.SetBgColor(eng.NewColor(0.1, 0.6, 0.9))
	batch = eng.NewBatch()
	texture := eng.NewTexture(eng.Files.Image("bot"))
	regions = texture.Split(64, 64)
}

func (g *Game) Update(dt float32) {
	angle += 15 * dt
}

func (g *Game) Draw() {
	batch.Begin()
	batch.Draw(regions[0], eng.Width()/2, eng.Height()/2, 0.5, 0.5, 5, 5, angle)
	batch.End()
}

func main() {
	eng.Run("Region", 1024, 640, false, new(Game))
}
