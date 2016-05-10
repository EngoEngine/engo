package main

import (
	"image/color"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/demos/demoutils"
)

type DefaultScene struct{}

var (
	worldWidth  int = 800
	worldHeight int = 800
)

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)
	w.AddSystem(&core.RenderSystem{})

	demoutils.NewBackground(w, worldWidth, worldHeight, color.RGBA{102, 153, 0, 255}, color.RGBA{102, 173, 0, 255})

	// We issue one camera zoom command at the start, but it takes a while to process because we set a duration
	engo.Mailbox.Dispatch(engo.CameraMessage{
		Axis:        engo.ZAxis,
		Value:       3, // so zooming out a lot
		Incremental: true,
		Duration:    time.Second * 5,
	})
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "IncrementalCamera Demo",
		Width:  worldWidth,
		Height: worldHeight,
	}
	engo.Run(opts, &DefaultScene{})
}
