package main

import (
	"github.com/paked/engi"
	"log"
)

var (
	bot   engi.Drawable
	batch *engi.Batch
	World *GameWorld
)

type GameWorld struct {
	bot   engi.Drawable
	batch *engi.Batch
	font  *engi.Font
	engi.World
}

func (game *GameWorld) Preload() {
	engi.Files.Add("bot", "data/icon.png")
	engi.Files.Add("font", "data/font.png")
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&RenderSystem{})

	entity := engi.NewEntity([]string{"RenderSystem"})
	component := NewRenderComponent(engi.Files.Image("bot"), engi.Point{0, 0}, engi.Point{10, 10})
	entity.AddComponent(component)
	game.AddEntity(entity)

	entityTwo := engi.NewEntity([]string{"RenderSystem"})
	componentTwo := NewRenderComponent(engi.Files.Image("font"), engi.Point{100, 100}, engi.Point{1, 1})

	entityTwo.AddComponent(componentTwo)
	game.AddEntity(entityTwo)
	log.Println("Setup")
}

type RenderSystem struct {
	*engi.System
}

func (rs *RenderSystem) New() {
	rs.System = &engi.System{}
}

func (rs RenderSystem) Pre() {
	engi.Gl.Clear(engi.Gl.COLOR_BUFFER_BIT)
	World.batch.Begin()
}

func (rs RenderSystem) Post() {
	World.batch.End()
}

func (rs RenderSystem) Update(entity *engi.Entity, dt float32) {
	component, ok := entity.GetComponent("RenderComponent").(RenderComponent)
	if ok {
		switch component.Display.(type) {
		case engi.Drawable:
			drawable := component.Display.(engi.Drawable)
			// World.batch.Draw(drawable, 512, 320, 0.5, 0.5, 10, 10, 0, 0xffffff, 1)
			World.batch.Draw(drawable, component.Position.X, component.Position.Y, 0, 0, component.Scale.X, component.Scale.Y, 0, 0xffffff, 1)
		case engi.Font:
			font := component.Display.(engi.Font)
			font.Print(batch, "Hello", 0, 0, 0x000)
		}
	}
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}

type RenderComponent struct {
	Display  interface{}
	Position engi.Point
	Scale    engi.Point
}

func NewRenderComponent(display interface{}, position, scale engi.Point) RenderComponent {
	return RenderComponent{Display: display, Position: position, Scale: scale}
}

func (rc RenderComponent) Name() string {
	return "RenderComponent"
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
