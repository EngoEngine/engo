package main

import (
	"log"

	"github.com/paked/engi"
)

var World *GameWorld

type GameWorld struct {
	engi.World
}

func (game *GameWorld) Preload() {
	game.New()
	engi.Files.Add(engi.NewResource("bot", "data/icon.png"))

	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&engi.RenderSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engi.Files.Image("bot")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{0, 0}, width, height}

	guy.AddComponent(&render)
	guy.AddComponent(&space)

	game.AddEntity(guy)
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
