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

	// Adding camera controllers so we can verify it doesn't break when we move
	w.AddSystem(engo.NewKeyboardScroller(scrollSpeed, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis))
	w.AddSystem(&engo.MouseZoomer{zoomSpeed})
	w.AddSystem(&engo.MouseRotator{RotationSpeed: 0.125})

	triangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle1.SpaceComponent = engo.SpaceComponent{Width: 100, Height: 100}
	triangle1.RenderComponent = engo.RenderComponent{Drawable: engo.Triangle{}, Color: color.RGBA{255, 0, 0, 255}}
	triangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&triangle1.BasicEntity, &triangle1.RenderComponent, &triangle1.SpaceComponent)
		}
	}

	rectangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	rectangle1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{100, 100}, Width: 100, Height: 100}
	rectangle1.RenderComponent = engo.RenderComponent{Drawable: engo.Rectangle{}, Color: color.RGBA{0, 255, 0, 255}}
	rectangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&rectangle1.BasicEntity, &rectangle1.RenderComponent, &rectangle1.SpaceComponent)
		}
	}

	circle1 := MyShape{BasicEntity: ecs.NewBasic()}
	circle1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{200, 200}, Width: 100, Height: 100}
	circle1.RenderComponent = engo.RenderComponent{Drawable: engo.Circle{}, Color: color.RGBA{0, 0, 255, 255}}
	circle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&circle1.BasicEntity, &circle1.RenderComponent, &circle1.SpaceComponent)
		}
	}

	triangle2 := MyShape{BasicEntity: ecs.NewBasic()}
	triangle2.SpaceComponent = engo.SpaceComponent{Position: engo.Point{300, 300}, Width: 100, Height: 100}
	triangle2.RenderComponent = engo.RenderComponent{Drawable: engo.Triangle{TriangleType: engo.TriangleRight}, Color: color.RGBA{255, 255, 0, 255}}
	triangle2.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&triangle2.BasicEntity, &triangle2.RenderComponent, &triangle2.SpaceComponent)
		}
	}

	line1 := MyShape{BasicEntity: ecs.NewBasic()}
	line1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{400, 400}, Width: 1, Height: 100}
	line1.RenderComponent = engo.RenderComponent{Drawable: engo.Rectangle{}, Color: color.RGBA{0, 255, 255, 255}}
	line1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&line1.BasicEntity, &line1.RenderComponent, &line1.SpaceComponent)
		}
	}

	complexTriangle1 := MyShape{BasicEntity: ecs.NewBasic()}
	complexTriangle1.SpaceComponent = engo.SpaceComponent{Position: engo.Point{500, 500}, Width: 100, Height: 100}
	complexTriangle1.RenderComponent = engo.RenderComponent{Drawable: engo.ComplexTriangles{
		Points: []engo.Point{
			{0.0, 0.0}, {1.0, 0.25}, {0.5, 0.5},
			{0.5, 0.5}, {1.0, 0.75}, {0.0, 1.0},
			{0.0, 0.0}, {0.5, 0.50}, {0.0, 1.0},
		}}, Color: color.RGBA{255, 0, 255, 255}}
	complexTriangle1.RenderComponent.SetShader(engo.LegacyShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&complexTriangle1.BasicEntity, &complexTriangle1.RenderComponent, &complexTriangle1.SpaceComponent)
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
