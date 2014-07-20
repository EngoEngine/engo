package main

import (
	"github.com/ajhager/eng"
)

type Game struct {
	*eng.Stage
}

func NewGame() *Game {
	return &Game{eng.NewStage()}
}

func (game *Game) Load() {
	eng.Files.Add("bot", "data/bot.png")
	eng.Files.Add("font", "data/font.png")
}

func (game *Game) Setup() {
	game.SetBg(eng.NewColor(45, 54, 56, 1))

	texture := eng.NewTexture(eng.Files.Image("bot"))
	regions := texture.Split(64, 64)
	font := eng.NewGridFont(eng.Files.Image("font"), 20, 20, "")

	bot := game.Sprite(regions[0], eng.Width()/2, eng.Height()/3)
	bot.Scale.SetTo(3)
	bot.Pivot.Y = 0

	text := game.Text(font, eng.Width()/2, eng.Height()/3, "ENG!")
	text.Scale.SetTo(1.5)
	text.Pivot.Set(0.5, 1)
	text.SetColor(eng.NewColor(108, 183, 103, 1))
}

func main() {
	eng.Open("Hello", 1024, 640, true, NewGame())
}
