package main

import (
	"math/rand"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	engi.Files.AddFromDir("assets", false)
}

func (game *GameWorld) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&HideSystem{})

	guy := ecs.NewEntity([]string{"RenderSystem", "HideSystem"})
	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")
	collision := &engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	w.AddEntity(guy)
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Type() string { return "GameWorld" }

type HideSystem struct {
	*ecs.System
}

func (HideSystem) Type() string {
	return "HideSystem"
}

func (s *HideSystem) New(*ecs.World) {
	s.System = ecs.NewSystem()
}

func (c *HideSystem) Update(e *ecs.Entity, dt float32) {
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
	opts := engi.RunOptions{
		Title:  "Show and Hide Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Open(opts, &GameWorld{})
}
