package main

import (
	"github.com/ajhager/eng"
	"math/rand"
	"strconv"
)

var (
	batch   *eng.Batch
	font    *eng.Font
	regions []*eng.Region
	bots    []*Sprite
	on      bool
	num     int
)

type Sprite struct {
	X, Y   float32
	DX, DY float32
	Image  *eng.Region
}

type Game struct {
	*eng.Game
}

func (g *Game) Load() {
	eng.Files.Add("bot", "data/bot.png")
	eng.Files.Add("font", "data/font.png")
}

func (g *Game) Setup() {
	font = eng.NewGridFont(eng.Files.Image("font"), 20, 20, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
	batch = eng.NewBatch()
	texture := eng.NewTexture(eng.Files.Image("bot"))
	regions = texture.Split(64, 64)
}

func (g *Game) Update(dt float32) {
	if on {
		for i := 0; i < 10; i++ {
			bots = append(bots, &Sprite{0, 0, rand.Float32() * 500, (rand.Float32() * 500) - 250, regions[0]})
		}
		num += 10
	}

	minX := float32(0)
	maxX := float32(eng.Width()) - regions[0].Width()
	minY := float32(0)
	maxY := float32(eng.Height()) - regions[0].Height()

	for _, bot := range bots {
		bot.X += bot.DX * dt
		bot.Y += bot.DY * dt
		bot.DY += 750 * dt

		if bot.X < minX {
			bot.DX *= -1
			bot.X = minX
		} else if bot.X > maxX {
			bot.DX *= -1
			bot.X = maxX
		}

		if bot.Y < minY {
			bot.DY = 0
			bot.Y = minY
		} else if bot.Y > maxY {
			bot.DY *= -.85
			bot.Y = maxY
			if rand.Float32() > 0.5 {
				bot.DY -= rand.Float32() * 200
			}
		}
	}
}

func (g *Game) Draw() {
	n := strconv.FormatInt(int64(num), 10)
	batch.Begin()
	for _, bot := range bots {
		batch.Draw(bot.Image, bot.X, bot.Y, 0.5, 0.5, 0.75, 0.75, 0)
	}
	font.Print(batch, n, 0, 0)
	batch.End()
}

func (g *Game) Mouse(x, y float32, a eng.Action) {
	switch a {
	case eng.MOVE:
	case eng.PRESS:
		on = true
	case eng.RELEASE:
		on = false
	}
}

func main() {
	eng.Run("Hello", 1024, 640, true, new(Game))
}
