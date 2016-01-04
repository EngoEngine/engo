package main

import (
	"log"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engi.Files.AddFromDir("data", false)

	log.Println("Preloaded")
}

func (game *GameWorld) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})

	// Create an entity part of the Render and Scale systems
	guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	// Retrieve a texture
	texture := engi.Files.Image("icon.png")

	// Create RenderComponent... Set scale to 8x, give lable "guy"
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{engi.Point{0, 0}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)

	w.AddEntity(guy)
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engi.RunOptions{
		Title:  "Hello Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Open(opts, &GameWorld{})
}
