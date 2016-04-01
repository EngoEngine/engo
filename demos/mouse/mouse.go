package main

import (
	"image"
	"image/color"
	"log"

	"github.com/engoengine/engo"
	"github.com/engoengine/ecs"
)

var (
	boxWidth  float32 = 150
	boxHeight float32 = 150
)

type GameWorld struct{}

func (game *GameWorld) Preload() {}

// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *engo.RenderComponent {
	rect := image.Rect(0, 0, int(boxWidth), int(boxHeight))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{102, 153, 0, 255}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			img.Set(i, j, c1)
		}
	}
	bgTexture := engo.NewImageObject(img)
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1}, "Background1")
	fieldRender.SetPriority(engo.Background)
	return fieldRender
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBg(color.White)

	w.AddSystem(&engo.MouseSystem{})
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&engo.MouseZoomer{-0.125})

	err := w.AddEntity(game.CreateEntity())
	if err != nil {
		log.Println(err)
	}
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func (game *GameWorld) CreateEntity() *ecs.Entity {
	entity := ecs.NewEntity([]string{"MouseSystem", "RenderSystem", "ControlSystem"})

	entity.AddComponent(generateBackground())
	entity.AddComponent(&engo.MouseComponent{})
	entity.AddComponent(&engo.SpaceComponent{engo.Point{0, 0}, boxWidth, boxHeight})

	return entity
}

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }

func (c *ControlSystem) New(*ecs.World) {}

func (c *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var mouse *engo.MouseComponent
	if !entity.Component(&mouse) {
		return
	}

	if mouse.Enter {
		engo.SetCursor(engo.Hand)
	} else if mouse.Leave {
		engo.SetCursor(nil)
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Mouse Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &GameWorld{})
}
