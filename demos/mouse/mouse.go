package main

import (
	"image"
	"image/color"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

var (
	boxWidth  float32 = 150
	boxHeight float32 = 150
)

type GameWorld struct{}

func (game *GameWorld) Preload() {}

// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *engi.RenderComponent {
	rect := image.Rect(0, 0, int(boxWidth), int(boxHeight))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{102, 153, 0, 255}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			img.Set(i, j, c1)
		}
	}
	bgTexture := engi.NewImageObject(img)
	fieldRender := engi.NewRenderComponent(engi.NewTexture(bgTexture), engi.Point{1, 1}, "Background1")
	fieldRender.SetPriority(engi.Background)
	return fieldRender
}

func (game *GameWorld) Setup(w *ecs.World) {
	engi.SetBg(0xFFFFFF)

	w.AddSystem(&engi.MouseSystem{})
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(engi.NewMouseZoomer(-0.125))

	w.AddEntity(game.CreateEntity())
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func (game *GameWorld) CreateEntity() *ecs.Entity {
	entity := ecs.NewEntity([]string{"MouseSystem", "RenderSystem", "ControlSystem"})

	entity.AddComponent(generateBackground())
	entity.AddComponent(&engi.MouseComponent{})
	entity.AddComponent(&engi.SpaceComponent{engi.Point{0, 0}, boxWidth, boxHeight})

	return entity
}

type ControlSystem struct {
	*ecs.System
}

func (ControlSystem) Type() string {
	return "ControlSystem"
}

func (c *ControlSystem) New(*ecs.World) {
	c.System = ecs.NewSystem()
}

func (c *ControlSystem) Update(entity *ecs.Entity, dt float32) {
	var mouse *engi.MouseComponent
	if !entity.Component(&mouse) {
		return
	}

	if mouse.Enter {
		engi.SetCursor(engi.Hand)
	} else if mouse.Leave {
		engi.SetCursor(engi.Arrow)
	}
}

func main() {
	opts := engi.RunOptions{
		Title:  "Mouse Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Open(opts, &GameWorld{})
}
