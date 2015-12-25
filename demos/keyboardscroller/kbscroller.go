package main

import (
	"image"
	"image/color"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

type Game struct{}

var (
	scrollSpeed float32 = 700
	worldWidth  float32 = 800
	worldHeight float32 = 800
)

// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *ecs.Entity {
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
	field := ecs.NewEntity([]string{"RenderSystem"})
	fieldRender := engi.NewRenderComponent(engi.NewTexture(bgTexture), engi.Point{1, 1}, "Background1")
	fieldRender.SetPriority(engi.Background)
	fieldSpace := &engi.SpaceComponent{engi.Point{0, 0}, worldWidth, worldHeight}
	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}

func (game *Game) Preload() {}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }

// Setup is called before the main loop is started
func (game *Game) Setup(w *ecs.World) {
	engi.SetBg(0x222222)
	w.AddSystem(&engi.RenderSystem{})

	// The most important line in this whole demo:
	w.AddSystem(engi.NewKeyboardScroller(scrollSpeed, engi.W, engi.D, engi.S, engi.A))

	// Create the background; this way we'll see when we actually scroll
	w.AddEntity(generateBackground())
}

func main() {
	opts := engi.RunOptions{
		Title:  "KeyboardScroller Demo",
		Width:  int(worldWidth),
		Height: int(worldHeight),
	}
	engi.Open(opts, &Game{})
}
