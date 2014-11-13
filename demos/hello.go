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
	game.bot = engi.Files.Image("bot")
	game.font = engi.NewGridFont(engi.Files.Image("font"), 20, 20)
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
	World.font.Print(World.batch, "ENGI", 475, 200, 0xffffff)
	World.batch.Draw(World.bot, 512, 320, 0.5, 0.5, 10, 10, 0, 0xffffff, 1)
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}

type RenderComponent struct {
}

func main() {
	World = &GameWorld{}
	entity := engi.Entity{}
	World.AddEntity(&entity)

	World.AddSystem(RenderSystem{})

	engi.Open("Hello", 1024, 640, false, World)
}
