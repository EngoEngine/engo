package main

import (
	"github.com/ajhager/engi"
	"math/rand"
)

var (
	bots   []*Bot
	on     bool
	num    int
	region *engi.Region
)

type Bot struct {
	*engi.Sprite
	DX, DY float32
}

type Game struct {
	*engi.Stage
}

func NewGame() *Game {
	return &Game{Stage: engi.NewStage()}
}

func (game *Game) Preload() {
	game.Load("bot", "data/icon.png")
}

func (game *Game) Setup() {
	game.SetBg(0x2d3638)
	texture := engi.NewTexture(engi.Files.Image("bot"))
	region = engi.NewRegion(texture, 0, 0, texture.Width(), texture.Height())
}

var time float32

func (game *Game) Update() {
	time += game.Delta()
	if time > 1 {
		println(num)
		println(int(game.Fps()))
		time = 0
	}

	if on {
		for i := 0; i < 10; i++ {
			bot := &Bot{game.Sprite(region, 0, 0), rand.Float32() * 500, rand.Float32()*500 - 250}
			bot.Alpha = 0.75
			bot.Scale.SetTo(1.5)
			bots = append(bots, bot)
		}
		num += 10
	}

	minX := float32(0)
	maxX := game.Width()
	minY := float32(0)
	maxY := game.Height()

	dt := game.Delta()

	for _, bot := range bots {
		bot.Position.X += bot.DX * dt
		bot.Position.Y += bot.DY * dt
		bot.DY += 750 * dt

		if bot.Position.X < minX {
			bot.DX *= -1
			bot.Position.X = minX
		} else if bot.Position.X > maxX {
			bot.DX *= -1
			bot.Position.X = maxX
		}

		if bot.Position.Y < minY {
			bot.DY = 0
			bot.Position.Y = minY
		} else if bot.Position.Y > maxY {
			bot.DY *= -.85
			bot.Position.Y = maxY
			if rand.Float32() > 0.5 {
				bot.DY -= rand.Float32() * 200
			}
		}
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
	engi.Open("Botmark", 800, 600, false, NewGame())
}
