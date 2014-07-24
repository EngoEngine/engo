package main

import (
	"github.com/ajhager/engi"
	"math"
)

type Game struct {
	*engi.Stage
	bot  *engi.Sprite
	text *engi.Text
}

func NewGame() *Game {
	return &Game{Stage: engi.NewStage()}
}

func (game *Game) Preload() {
	game.Load("bot", "data/icon.png")
	//game.Load("bot", "data/engi.png")
	game.Load("font", "data/font.png")
}

func (game *Game) Setup() {
	game.SetBg(0x2d3638)

	texture := engi.NewTexture(engi.Files.Image("bot"))
	region := engi.NewRegion(texture, 0, 0, texture.Width(), texture.Height())

	bot := game.Sprite(region, game.Width()/2, game.Height()/1.75)
	bot.Pivot.Y = 1
	bot.Scale.SetTo(14)
	game.bot = bot

	font := engi.NewGridFont(engi.Files.Image("font"), 20, 20, "")
	text := game.Text(font, game.Width()/2, game.Height()/1.75, "ENGi")
	text.Scale.Set(3, 5)
	text.Pivot.Y = 0
	text.Tint = 0x6cb767
	game.text = text
}

var on bool

func (game *Game) Update() {
	if on {
		game.bot.Rotation = float32(math.Sin(float64(game.Time() * 200)))
	} else {
		game.bot.Rotation = 0
	}
}

func (game *Game) Mouse(x, y float32, action engi.Action) {
	switch action {
	case engi.PRESS:
		on = true
	case engi.RELEASE:
		on = false
	}
}

func main() {
	engi.Open("Hello", 1024, 640, true, NewGame())
}
