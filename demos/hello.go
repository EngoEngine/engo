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

	entity := engi.NewEntity([]string{"RenderSystem"})
	component := RenderComponent{engi.Files.Image("bot")}
	entity.AddComponent(component)
	game.AddEntity(entity)

	// entityTwo := engi.NewEntity([]string{"RenderSystem"})
	// componentTwo := RenderComponent{engi.NewGridFont(engi.Files.Image("font"), 20, 20)}
	// entityTwo.AddComponent(componentTwo)
	// game.AddEntity(entityTwo)

	game.AddSystem(RenderSystem{})
	game.AddSystem(engi.TestSystem{})

	// game.bot = engi.Files.Image("bot")
	// game.font =
	log.Println("Setup")
}

type RenderSystem struct {
}

func (rs RenderSystem) Pre() {
	engi.Gl.Clear(engi.Gl.COLOR_BUFFER_BIT)
	World.batch.Begin()
}

func (rs RenderSystem) Post() {
	World.batch.End()
}
func (rs RenderSystem) Update(entity *engi.Entity, dt float32) {
	// log.Println(entity.GetComponent("RenderComponent"))

	component, ok := entity.GetComponent("RenderComponent").(RenderComponent)
	if ok {
		switch component.Display.(type) {
		case engi.Drawable:
			drawable := component.Display.(engi.Drawable)
			World.batch.Draw(drawable, 512, 320, 0.5, 0.5, 10, 10, 0, 0xffffff, 1)
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
	Display interface{}
}

func (rc RenderComponent) Name() string {
	return "RenderComponent"
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
