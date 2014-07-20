package main

import (
	"github.com/ajhager/engi"
)

type Game struct {
	*engi.Stage
}

func NewGame() *Game {
	return &Game{engi.NewStage()}
}

func (game *Game) Load() {
	engi.Files.Add("bot", "data/bot.png")
	engi.Files.Add("font", "data/font.png")
}

func (game *Game) Setup() {
	game.SetBg(0x2d3638)

	texture := engi.NewTexture(engi.Files.Image("bot"))
	regions := texture.Split(64, 64)
	font := engi.NewGridFont(engi.Files.Image("font"), 20, 20, "")

	bot := game.Sprite(regions[0], game.Width()/2, game.Height()/3)
	bot.Scale.SetTo(3)
	bot.Pivot.Y = 0

	text := game.Text(font, game.Width()/2, game.Height()/3, "ENGi")
	text.Scale.SetTo(1.5)
	text.Pivot.Set(0.5, 1)
	text.Tint = 0x6cb767
}

func main() {
	engi.Open("Hello", 1024, 640, true, NewGame())
}
