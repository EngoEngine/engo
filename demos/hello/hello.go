package main

import (
	"log"
	"math/rand"

	"github.com/paked/engi"
)

var World *GameWorld

type GameWorld struct{}

func (game *GameWorld) Preload() {

	// This could be done individually: engi.Files.Add("data/icon.png"), etc
	// Second value (false) says whether to check recursively or not
	engi.Files.AddFromDir("data", false)

	log.Println("Preloaded")
}

func (game *GameWorld) Setup(w *engi.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&ScaleSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engi.Files.Image("icon.png")
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

type ScaleSystem struct {
	*engi.System
}

func (ScaleSystem) Type() string {
	return "ScaleSystem"
}

func (s *ScaleSystem) New() {
	s.System = engi.NewSystem()
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
	engi.Open("Hello Demo", 1024, 640, false, World)
}
