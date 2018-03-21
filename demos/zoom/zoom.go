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
	zoomSpeed float32 = -0.125

	worldWidth  int = 800
	worldHeight int = 800
)

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseZoomer{zoomSpeed})

	// Create the background; this way we'll see when we actually zoom
	demoutils.NewBackground(w, worldWidth, worldHeight, color.RGBA{102, 153, 0, 255}, color.RGBA{102, 173, 0, 255})

	// Center camera if GlobalScale is Setup
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.XAxis,
		Value:       float32(worldWidth) / 2,
		Incremental: false})
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.YAxis,
		Value:       float32(worldHeight) / 2,
		Incremental: false})
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "Zoom Demo",
		Width:  worldWidth,
		Height: worldHeight,
	}
	engo.Run(opts, &DefaultScene{})
}
