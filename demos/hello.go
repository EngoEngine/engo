package main

import (
	"github.com/ajhager/eng"
)

var (
	regions []*eng.Region
	stage   *eng.Stage
)

type Game struct {
	*eng.Game
}

func (g *Game) Load() {
	eng.Files.Add("bot", "data/bot.png")
	eng.Files.Add("font", "data/font.png")
}

func (g *Game) Setup() {
	eng.SetBgColor(eng.NewColorHex(0x2d3638))

	texture := eng.NewTexture(eng.Files.Image("bot"))
	regions = texture.Split(64, 64)

	font := eng.NewGridFont(eng.Files.Image("font"), 20, 20, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")

	stage = eng.NewStage()
	bot := eng.NewSprite(regions[0], eng.Width()/2, eng.Height()/3)
	bot.Scale.SetTo(3)
	bot.Pivot.Set(0.5, 0)
	stage.Add(bot)
	text := eng.NewText(font, eng.Width()/2, eng.Height()/3, "ENG!")
	text.Scale.SetTo(1.5)
	text.Pivot.Set(0.5, 1)
	text.SetColor(eng.NewColorHexA(0x6cb767, 1))
	stage.Add(text)
}

func (g *Game) Draw() {
	stage.Draw()
}

func main() {
	eng.Run("Hello", 1024, 640, true, new(Game))
}
