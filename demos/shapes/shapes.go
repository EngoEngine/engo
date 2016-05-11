package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/core"
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
	core.RenderComponent
	core.SpaceComponent
}

func (*DefaultScene) Preload() {}

// Setup is called before the main loop is started
func (*DefaultScene) Setup(w *ecs.World) {
	core.SetBackground(color.RGBA{55, 55, 55, 255})
	w.AddSystem(&core.RenderSystem{})

	// Adding camera controllers so we can verify it doesn't break when we move
	w.AddSystem(core.NewKeyboardScroller(scrollSpeed, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis))
	w.AddSystem(&core.MouseZoomer{zoomSpeed})
	w.AddSystem(&core.MouseRotator{RotationSpeed: 0.125})

	triangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle1.SpaceComponent = core.SpaceComponent{Width: 100, Height: 100}
	triangle1.RenderComponent = core.RenderComponent{Drawable: core.Triangle{}, Color: color.RGBA{255, 0, 0, 255}}
	triangle1.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&triangle1.BasicEntity, &triangle1.RenderComponent, &triangle1.SpaceComponent)
		}
	}

	rectangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	rectangle1.SpaceComponent = core.SpaceComponent{Position: engo.Point{100, 100}, Width: 100, Height: 100}
	rectangle1.RenderComponent = core.RenderComponent{Drawable: core.Rectangle{}, Color: color.RGBA{0, 255, 0, 255}}
	rectangle1.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&rectangle1.BasicEntity, &rectangle1.RenderComponent, &rectangle1.SpaceComponent)
		}
	}

	circle1 := MyShape{BasicEntity: ecs.NewBasic()}
	circle1.SpaceComponent = core.SpaceComponent{Position: engo.Point{200, 200}, Width: 100, Height: 100}
	circle1.RenderComponent = core.RenderComponent{Drawable: core.Circle{}, Color: color.RGBA{0, 0, 255, 255}}
	circle1.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&circle1.BasicEntity, &circle1.RenderComponent, &circle1.SpaceComponent)
		}
	}

	triangle2 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle2.SpaceComponent = core.SpaceComponent{Position: engo.Point{300, 300}, Width: 100, Height: 100}
	triangle2.RenderComponent = core.RenderComponent{Drawable: core.Triangle{TriangleType: core.TriangleRight}, Color: color.RGBA{255, 255, 0, 255}}
	triangle2.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&triangle2.BasicEntity, &triangle2.RenderComponent, &triangle2.SpaceComponent)
		}
	}

	line1 := MyShape{BasicEntity: ecs.NewBasic()}
	line1.SpaceComponent = core.SpaceComponent{Position: engo.Point{400, 400}, Width: 1, Height: 100}
	line1.RenderComponent = core.RenderComponent{Drawable: core.Rectangle{}, Color: color.RGBA{0, 255, 255, 255}}
	line1.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&line1.BasicEntity, &line1.RenderComponent, &line1.SpaceComponent)
		}
	}

	complexTriangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	complexTriangle1.SpaceComponent = core.SpaceComponent{Position: engo.Point{500, 500}, Width: 100, Height: 100}
	complexTriangle1.RenderComponent = core.RenderComponent{Drawable: core.ComplexTriangles{
		Points: []engo.Point{
			{0.0, 0.0}, {1.0, 0.25}, {0.5, 0.5},
			{0.5, 0.5}, {1.0, 0.75}, {0.0, 1.0},
			{0.0, 0.0}, {0.5, 0.50}, {0.0, 1.0},
		}}, Color: color.RGBA{255, 0, 255, 255}}
	complexTriangle1.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&complexTriangle1.BasicEntity, &complexTriangle1.RenderComponent, &complexTriangle1.SpaceComponent)
		}
	}

	triangle3 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle3.SpaceComponent = core.SpaceComponent{Position: engo.Point{23, 123}, Width: 50, Height: 50}
	triangle3.RenderComponent = core.RenderComponent{Drawable: core.Triangle{BorderWidth: 1, BorderColor: color.White}, Color: color.RGBA{255, 0, 0, 255}}
	triangle3.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&triangle3.BasicEntity, &triangle3.RenderComponent, &triangle3.SpaceComponent)
		}
	}

	rectangle2 := MyShape{BasicEntity: ecs.NewBasic()}
	rectangle2.SpaceComponent = core.SpaceComponent{Position: engo.Point{123, 223}, Width: 50, Height: 50}
	rectangle2.RenderComponent = core.RenderComponent{Drawable: core.Rectangle{BorderWidth: 1, BorderColor: color.White}, Color: color.RGBA{0, 255, 0, 255}}
	rectangle2.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&rectangle2.BasicEntity, &rectangle2.RenderComponent, &rectangle2.SpaceComponent)
		}
	}

	circle2 := MyShape{BasicEntity: ecs.NewBasic()}
	circle2.SpaceComponent = core.SpaceComponent{Position: engo.Point{223, 323}, Width: 50, Height: 50}
	circle2.RenderComponent = core.RenderComponent{Drawable: core.Circle{BorderWidth: 1, BorderColor: color.White}, Color: color.RGBA{0, 0, 255, 255}}
	circle2.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&circle2.BasicEntity, &circle2.RenderComponent, &circle2.SpaceComponent)
		}
	}

	triangle4 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle4.SpaceComponent = core.SpaceComponent{Position: engo.Point{323, 423}, Width: 50, Height: 50}
	triangle4.RenderComponent = core.RenderComponent{Drawable: core.Triangle{TriangleType: core.TriangleRight, BorderWidth: 1, BorderColor: color.White}, Color: color.RGBA{255, 255, 0, 255}}
	triangle4.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&triangle4.BasicEntity, &triangle4.RenderComponent, &triangle4.SpaceComponent)
		}
	}

	complexTriangle2 := MyShape{BasicEntity: ecs.NewBasic()}
	complexTriangle2.SpaceComponent = core.SpaceComponent{Position: engo.Point{523, 623}, Width: 50, Height: 50}
	complexTriangle2.RenderComponent = core.RenderComponent{Drawable: core.ComplexTriangles{
		BorderWidth: 1, BorderColor: color.White,
		Points: []engo.Point{
			{0.0, 0.0}, {1.0, 0.25}, {0.5, 0.5},
			{0.5, 0.5}, {1.0, 0.75}, {0.0, 1.0},
			{0.0, 0.0}, {0.5, 0.50}, {0.0, 1.0},
		}}, Color: color.RGBA{255, 0, 255, 255}}
	complexTriangle2.RenderComponent.SetShader(core.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.RenderSystem:
			sys.Add(&complexTriangle2.BasicEntity, &complexTriangle2.RenderComponent, &complexTriangle2.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:          "Shapes Demo",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
		MSAA:           4, // This one is not mandatory, but makes the shapes look so much better when rotating the camera
	}
	engo.Run(opts, &DefaultScene{})
}
