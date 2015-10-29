package main

import (
	"image"
	"image/color"

	"github.com/paked/engi"
)

type Game struct {
	engi.World
}

var (
	scrollSpeed float32 = 700
	worldWidth  float32 = 800
	worldHeight float32 = 800
)

// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *engi.Entity {
	rect := image.Rect(0, 0, int(worldWidth), int(worldHeight))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{102, 153, 0, 255}
	c2 := color.RGBA{102, 173, 0, 255}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			if i%40 > 20 {
				if j%40 > 20 {
					img.Set(i, j, c1)
				} else {
					img.Set(i, j, c2)
				}
			} else {
				if j%40 > 20 {
					img.Set(i, j, c2)
				} else {
					img.Set(i, j, c1)
				}
			}
		}
	}
	bgTexture := engi.NewImageObject(img)
	field := engi.NewEntity([]string{"RenderSystem"})
	fieldRender := engi.NewRenderComponent(engi.NewRegion(engi.NewTexture(bgTexture), 0, 0, int(worldWidth), int(worldHeight)), engi.Point{1, 1}, "Background1")
	fieldRender.Priority = engi.Background
	fieldSpace := &engi.SpaceComponent{engi.Point{0, 0}, worldWidth, worldHeight}
	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}

// Setup is called before the main loop is started
func (game *Game) Setup() {
	engi.SetBg(0x222222)
	game.AddSystem(&engi.RenderSystem{})

	// The most important line in this whole demo:
	game.AddSystem(engi.NewKeyboardScroller(scrollSpeed, engi.W, engi.D, engi.S, engi.A))

	// Create the background; this way we'll see when we actually scroll
	game.AddEntity(generateBackground())
}

func main() {
	engi.SetFPSLimit(120)
	engi.Open("KeyboardScroller Demo", 400, 400, false, &Game{})
}
