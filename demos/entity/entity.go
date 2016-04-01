package main

import (
	"image/color"
	"log"

	"github.com/engoengine/engo"
	"github.com/engoengine/engo/ecs"
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("data", false)

	log.Println("Preloaded")
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBg(color.White)

	w.AddSystem(&engo.RenderSystem{})

	// Create an entity part of the Render
	guy := ecs.NewEntity([]string{"RenderSystem"})
	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create RenderComponent... Set scale to 8x, give lable "guy"
	render := engo.NewRenderComponent(texture, engo.Point{8, 8}, "guy")

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engo.SpaceComponent{engo.Point{0, 0}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)

	err := w.AddEntity(guy)
	if err != nil {
		log.Println(err)
	}
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Hello Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &GameWorld{})
}
