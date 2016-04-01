package main

import (
	"image"
	"image/color"
	"log"

	"github.com/engoengine/engo"
	"github.com/engoengine/ecs"
)

type Game struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700
	worldWidth  float32 = 800
	worldHeight float32 = 800

	hudBackgroundPriority = engo.PriorityLevel(engo.HUDGround)
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
	bgTexture := engo.NewImageObject(img)
	field := ecs.NewEntity([]string{"RenderSystem"})
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1}, "Background1")
	fieldRender.SetPriority(engo.Background)
	fieldSpace := &engo.SpaceComponent{engo.Point{0, 0}, worldWidth, worldHeight}
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
	bgTexture := engo.NewImageObject(img)
	field := ecs.NewEntity([]string{"RenderSystem"})
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1}, "HUDBackground1")
	fieldRender.SetPriority(hudBackgroundPriority)
	fieldSpace := &engo.SpaceComponent{engo.Point{-1, -1}, width, height}
	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}

func (game *Game) Preload() {}

// Setup is called before the main loop is started
func (game *Game) Setup(w *ecs.World) {
	engo.SetBg(color.White)
	w.AddSystem(&engo.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(engo.NewKeyboardScroller(scrollSpeed, engo.W, engo.D, engo.S, engo.A))
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	// Create background, so we can see difference between this and HUD
	err := w.AddEntity(generateBackground())
	if err != nil {
		log.Println(err)
	}

	// Creating the HUD
	hudWidth := float32(200)         // Can be anything you want
	hudHeight := engo.WindowHeight() // Can be anything you want

	// Generate something that uses the PriorityLevel HUDGround or up
	hudBg := generateHUDBackground(hudWidth, hudHeight)
	err = w.AddEntity(hudBg)
	if err != nil {
		log.Println(err)
	}
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "HUD Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &Game{})
}
