package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

type DefaultScene struct{}

var (
	zoomSpeed   float32 = -0.125
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

type MyTriangle struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.RGBA{5, 5, 5, 255})
	w.AddSystem(&engo.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(engo.NewKeyboardScroller(scrollSpeed, engo.W, engo.D, engo.S, engo.A))
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	triangle1 := MyTriangle{BasicEntity: ecs.NewBasic()}
	triangle1.SpaceComponent = engo.SpaceComponent{Width: 400, Height: 200}
	triangle1.RenderComponent = engo.RenderComponent{Drawable: engo.Triangle{}, Color: color.RGBA{255, 0, 0, 255}}
	triangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&triangle1.BasicEntity, &triangle1.RenderComponent, &triangle1.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "Triangle Demo",
		Width:  worldWidth,
		Height: worldHeight,
	}
	engo.Run(opts, &DefaultScene{})
}
