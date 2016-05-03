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

type MyShape struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.RGBA{55, 55, 55, 255})
	w.AddSystem(&engo.RenderSystem{})

	// Adding KeyboardScroller so we can actually see the difference between background and HUD when scrolling
	w.AddSystem(engo.NewKeyboardScroller(scrollSpeed, engo.W, engo.D, engo.S, engo.A))
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})

	triangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle1.SpaceComponent = engo.SpaceComponent{Width: 200, Height: 200}
	triangle1.RenderComponent = engo.RenderComponent{Drawable: engo.Triangle{}, Color: color.RGBA{255, 0, 0, 255}}
	triangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&triangle1.BasicEntity, &triangle1.RenderComponent, &triangle1.SpaceComponent)
		}
	}

	rectangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	rectangle1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{200, 200}, Width: 200, Height: 200}
	rectangle1.RenderComponent = engo.RenderComponent{Drawable: engo.Rectangle{}, Color: color.RGBA{0, 255, 0, 255}}
	rectangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&rectangle1.BasicEntity, &rectangle1.RenderComponent, &rectangle1.SpaceComponent)
		}
	}

	circle1 := MyShape{BasicEntity: ecs.NewBasic()}
	circle1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{400, 400}, Width: 200, Height: 200}
	circle1.RenderComponent = engo.RenderComponent{Drawable: engo.Circle{}, Color: color.RGBA{0, 0, 255, 255}}
	circle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&circle1.BasicEntity, &circle1.RenderComponent, &circle1.SpaceComponent)
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
