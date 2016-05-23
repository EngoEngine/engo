package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

var (
	worldWidth  int = 800
	worldHeight int = 800
)

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.RGBA{55, 55, 55, 255})
	w.AddSystem(&common.RenderSystem{})
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "Text Demo",
		Width:  worldWidth,
		Height: worldHeight,
	}
	engo.Run(opts, &DefaultScene{})
}
