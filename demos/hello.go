package main

import (
	"github.com/ajhager/engi"
	"math"
)

type Game struct {
	*engi.Stage
	bot *engi.Sprite
	//text *engi.Text
}

func NewGame() *Game {
	return &Game{Stage: engi.NewStage()}
}

func (game *Game) Preload() {
	game.Load("bot", "data/icon.png")
	game.Load("font", "data/font.png")
}

func (game *Game) Setup() {
	game.SetBg(0x2d3638)

	texture := engi.NewTexture(engi.Files.Image("bot"))
	region := engi.NewRegion(texture, 0, 0, texture.Width(), texture.Height())

	bot := engi.NewSprite(region, game.Width()/2, game.Height()/1.75)
	game.AddChild(bot)
	bot.Anchor.Set(0.5, 1)
	bot.Scale.SetTo(14)
	game.bot = bot

	bot2 := engi.NewSprite(region, 0, 5)
	bot.AddChild(bot2)
	bot2.Anchor.Set(0.5, 0)
	bot2.Scale.SetTo(0.33)

	/*
		font := engi.NewGridFont(engi.Files.Image("font"), 20, 20, "")
		text := engi.NewText(font, game.Width()/2, game.Height()/1.75, "ENGi")
		game.AddChild(text.Sprite)
		text.Anchor.Set(0.5, 0)
		text.Scale.Set(3, 4)
		text.SetTint(0x6cb767)
		game.text = text
	*/
}

var on bool

func (game *Game) Update() {
	//game.text.SetText(strconv.FormatInt(int64(game.Fps()), 10))
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
