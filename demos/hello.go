package main

import (
	"log"
	"math/rand"

	"github.com/paked/engi"
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
	game.AddSystem(&ScaleSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ScaleSystem"})
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

type ScaleSystem struct {
	*engi.System
}

func (s *ScaleSystem) Name() string {
	return "ScaleSystem"
}

func (s *ScaleSystem) New() {
	s.System = &engi.System{}
}

func (c *ScaleSystem) Update(e *engi.Entity, dt float32) {
	var render *engi.RenderComponent
	if !e.GetComponent(&render) {
		return
	}
	var mod float32

	if rand.Int()%2 == 0 {
		mod = 0.1
	} else {
		mod = -0.1
	}

	if render.Scale.X+mod >= 15 || render.Scale.X+mod <= 1 {
		mod *= -1
	}

	render.Scale.AddScalar(mod)
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
