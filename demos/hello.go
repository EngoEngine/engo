package main

import (
	"github.com/ajhager/engi"
)

type Game struct {
	*engi.Game
	bot   *engi.Region
	batch *engi.Batch
}

func (game *Game) Preload() {
	engi.Files.Add("bot", "data/icon.png")
	game.batch = engi.NewBatch(1024, 640)
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	texture := engi.NewTexture(engi.Files.Image("bot"))
	game.bot = engi.NewRegion(texture, 0, 0, texture.Width(), texture.Height())
}

func (game *Game) Render() {
	game.batch.Begin()
	game.batch.Draw(game.bot, 512, 320, 0.5, 0.5, 20, 20, 0, 0xffffff, 1)
	game.batch.End()
}

func main() {
	engi.Open("Hello", 1024, 640, false, &Game{})
}
