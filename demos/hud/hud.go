package main

import (
	"image"
	"image/color"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

type Game struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700
	worldWidth  float32 = 800
	worldHeight float32 = 800

	hudBackgroundPriority = engi.PriorityLevel(engi.HUDGround)
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

// generateHUDBackground creates a violet HUD on the left side of the screen - might be inefficient
func generateHUDBackground(width, height float32) *ecs.Entity {
	rect := image.Rect(0, 0, int(width), int(height))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{255, 0, 255, 180}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			img.Set(i, j, c1)
		}
	}
	bgTexture := engi.NewImageObject(img)
	field := ecs.NewEntity([]string{"RenderSystem"})
	fieldRender := engi.NewRenderComponent(engi.NewTexture(bgTexture), engi.Point{1, 1}, "HUDBackground1")
	fieldRender.SetPriority(hudBackgroundPriority)
	fieldSpace := &engi.SpaceComponent{engi.Point{-1, -1}, width, height}
	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}

func (game *Game) Preload() {}

// Setup is called before the main loop is started
func (game *Game) Setup(w *ecs.World) {
	engi.SetBg(0x222222)
	w.AddSystem(&engi.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(engi.NewKeyboardScroller(scrollSpeed, engi.W, engi.D, engi.S, engi.A))
	w.AddSystem(engi.NewMouseZoomer(zoomSpeed))

	// Create background, so we can see difference between this and HUD
	w.AddEntity(generateBackground())

	// Creating the HUD
	hudWidth := float32(200)         // Can be anything you want
	hudHeight := engi.WindowHeight() // Can be anything you want

	// Generate something that uses the PriorityLevel HUDGround or up
	hudBg := generateHUDBackground(hudWidth, hudHeight)
	w.AddEntity(hudBg)
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }

func main() {
	opts := engi.RunOptions{
		Title:  "HUD Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Open(opts, &Game{})
}
