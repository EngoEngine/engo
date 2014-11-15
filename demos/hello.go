package main

import (
	"github.com/paked/engi"
	"log"
)

var World *GameWorld

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
	texture := engi.Files.Image("bot")
	render := NewRenderComponent(texture, engi.Point{0, 0}, engi.Point{1, 1}, "bot")
	space := SpaceComponent{Position: engi.Point{10, 10}, Width: texture.Width(), Height: texture.Height()}
	entity.AddComponent(&render)
	entity.AddComponent(&space)
	game.AddEntity(entity)

	// entityTwo := engi.NewEntity([]string{"RenderSystem"})
	// componentTwo := NewRenderComponent(engi.NewGridFont(engi.Files.Image("font"), 20, 20), engi.Point{200, 400}, engi.Point{1, 1}, "YOLO MATE WASSUP")

	// entityTwo.AddComponent(&componentTwo)
	// game.AddEntity(entityTwo)
	// log.Println("Setup")
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

var pos float32

func (rs *RenderSystem) Update(entity *engi.Entity, dt float32) {
	// component, ok := entity.GetComponent("RenderComponent").(RenderComponent)
	// space, ok := entity.GetComponent("SpaceComponent").(SpaceComponent)
	// if ok {
	// 	switch component.Display.(type) {
	// 	case engi.Drawable:
	// 		drawable := component.Display.(engi.Drawable)
	// 		World.batch.Draw(drawable, component.Position.X, component.Position.Y, 0, 0, component.Scale.X, component.Scale.Y, 0, 0xffffff, 1)
	// 	case *engi.Font:
	// 		font := component.Display.(*engi.Font)
	// 		font.Print(World.batch, component.Label, component.Position.X, component.Position.Y, 0xffffff)
	// 	}
	// }

	render, hasRender := entity.GetComponent("RenderComponent").(*RenderComponent)
	space, hasSpace := entity.GetComponent("SpaceComponent").(*SpaceComponent)
	if hasRender && hasSpace {
		log.Println(space.Position.X, space.Position.Y)
		// pos += 1
		// pos += 1 * dt
		// log.Println(dt)
		// space.Position.Y += .30 * dt
		// log.Println(pos, "pos")
		switch render.Display.(type) {
		case engi.Drawable:
			drawable := render.Display.(engi.Drawable)
			World.batch.Draw(drawable, space.Position.X, space.Position.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, 0xffffff, 1)
		}

		// space.Position.Y += 1 * dt
		log.Println(space.Position)
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
	Label    string
}

func NewRenderComponent(display interface{}, position, scale engi.Point, label string) RenderComponent {
	return RenderComponent{Display: display, Position: position, Scale: scale, Label: label}
}

func (rc RenderComponent) Name() string {
	return "RenderComponent"
}

type SpaceComponent struct {
	Position engi.Point
	Width    float32
	Height   float32
}

func (sc SpaceComponent) Name() string {
	return "SpaceComponent"
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
