package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/demos/demoutils"
)

type DefaultScene struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700

<<<<<<< HEAD
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
	field := ecs.NewEntity("RenderSystem")
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1})
	fieldSpace := &engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    worldWidth,
		Height:   worldHeight,
	}

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
	field := ecs.NewEntity("RenderSystem")
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1})
	fieldRender.SetShader(engo.HUDShader)
	fieldRender.SetZIndex(1) // A value larger than 0 (default), to ensure being drawn on top of the background
	fieldSpace := &engo.SpaceComponent{
		Position: engo.Point{-1, -1},
		Width:    width,
		Height:   height,
	}

	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}
=======
	worldWidth  int = 800
	worldHeight int = 800
)
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)
	w.AddSystem(&engo.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(engo.NewKeyboardScroller(scrollSpeed, engo.W, engo.D, engo.S, engo.A))
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	// Create background, so we can see difference between this and HUD
	demoutils.NewBackground(w, worldWidth, worldHeight, color.RGBA{102, 153, 0, 255}, color.RGBA{102, 173, 0, 255})

	// Define parameters for the hud
	hudWidth := 200                       // Can be anything you want
	hudHeight := int(engo.WindowHeight()) // Can be anything you want

	// Generate something that uses the PriorityLevel HUDGround or up. We're giving the same color twice,
	// so it'll create one solid color.
	hudBg := demoutils.NewBackground(w, hudWidth, hudHeight, color.RGBA{255, 0, 255, 180}, color.RGBA{255, 0, 255, 180})

	// These adjustments are needed to transform it into a HUD:
	hudBg.RenderComponent.SetZIndex(1) // something bigger than default (0), so it'll be on top of the regular background
	hudBg.RenderComponent.SetShader(engo.HUDShader)
}

<<<<<<< HEAD
func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Exit()        {}
func (*Game) Type() string { return "Game" }
=======
func (*DefaultScene) Type() string { return "Game" }
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c

func main() {
	opts := engo.RunOptions{
		Title:  "HUD Demo",
<<<<<<< HEAD
		Width:  1024,
		Height: 640,
=======
		Width:  worldWidth,
		Height: worldHeight,
>>>>>>> 28393c45ef7ce198babe3c6854931398faaba25c
	}
	engo.Run(opts, &DefaultScene{})
}
