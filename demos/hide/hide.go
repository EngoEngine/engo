package main

import (
	"math/rand"

	"github.com/paked/engi"
)

var World *GameWorld

type GameWorld struct{}

func (game *GameWorld) Preload() {
	engi.Files.AddFromDir("assets", false)
}

func (game *GameWorld) Setup(w *engi.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&HideSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "HideSystem"})
	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")
	collision := &engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := &engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	w.AddEntity(guy)
}

type HideSystem struct {
	*engi.System
}

func (HideSystem) Type() string {
	return "HideSystem"
}

func (s *HideSystem) New(*engi.World) {
	s.System = engi.NewSystem()
}

func (c *HideSystem) Update(e *engi.Entity, dt float32) {
	var render *engi.RenderComponent
	if !e.Component(&render) {
		return
	}
	if rand.Int()%10 == 0 {
		render.SetPriority(engi.Hidden)
	} else {
		render.SetPriority(engi.MiddleGround)
	}
}

func main() {
	World = &GameWorld{}
	engi.Open("Show and Hide Demo", 1024, 640, false, World)
}
