//+build demo

package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/engo/demos/demoutils"
)

type DefaultScene struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)
	w.AddSystem(&common.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(common.NewKeyboardScroller(scrollSpeed, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis))
	w.AddSystem(&common.MouseZoomer{zoomSpeed})

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
	hudBg.RenderComponent.SetShader(common.HUDShader)
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:          "HUD Demo",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &DefaultScene{})
}
