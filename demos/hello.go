package main

import (
	"github.com/paked/engi"
	"log"
)

var World *GameWorld

type GameWorld struct {
	engi.World
}

func (game *GameWorld) Preload() {
	game.New()
	engi.Files.Add(engi.NewResource("bot", "data/icon.png"),
		engi.NewResource("font", "data/font.png"),
		engi.NewResource("rock", "data/rock.png"),
		engi.NewResource("sheet", "data/sheet.png"),
		engi.NewResource("sample", "data/Hero.png"))

	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&engi.RenderSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ControlSystem", "CollisionSystem", "DeathSystem"})
	texture := engi.Files.Image("bot")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")
	collision := engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(&render)
	guy.AddComponent(&space)
	guy.AddComponent(&collision)

	game.AddEntity(guy)

}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
