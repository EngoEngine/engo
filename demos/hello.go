package main

import (
	"github.com/ajhager/engi"
)

type Game struct {
	*engi.Game
	bot   engi.Drawable
	batch *engi.Batch
	font  *engi.Font
}

func (game *Game) Preload() {
	engi.Files.Add("bot", "data/icon.png")
	engi.Files.Add("font", "data/font.png")
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.bot = engi.Files.Image("bot")
	game.font = engi.NewGridFont(engi.Files.Image("font"), 20, 20)
}

func (game *Game) Render() {
	game.batch.Begin()
	game.font.Print(game.batch, "ENGI", 475, 200, 0xffffff)
	game.batch.Draw(game.bot, 512, 320, 0.5, 0.5, 10, 10, 0, 0xffffff, 1)
	game.batch.End()
}

func main() {
	engi.Open("Hello", 1024, 640, false, &Game{})
}
