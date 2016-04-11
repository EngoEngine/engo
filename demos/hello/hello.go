package main

import (
	"image/color"
	"log"
	"fmt"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
)

type GameWorld struct{}

func (game *GameWorld) Preload() {

	// This could be done individually: engo.Files.Add("data/icon.png"), etc
	// Second value (false) says whether to check recursively or not
	engo.Files.AddFromDir("data", false)

	log.Println("Preloaded")
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ScaleSystem{})

	guy := ecs.NewEntity("RenderSystem", "ScaleSystem")
	texture := engo.Files.Image("icon.png")
	render := engo.NewRenderComponent(texture, engo.Point{8, 8}, "guy")
	collision := &engo.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engo.SpaceComponent{engo.Point{(engo.Width() - width) / 2, (engo.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	err := w.AddEntity(guy)
	if err != nil {
		log.Println(err)
	}
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Exit() 		{}
func (*GameWorld) Type() string { return "GameWorld" }

type ScaleSystem struct {
	ecs.LinearSystem
}

func (*ScaleSystem) Type() string { return "ScaleSystem" }

func (s *ScaleSystem) New(*ecs.World) {}

func (c *ScaleSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var render *engo.RenderComponent
	if !entity.Component(&render) {
		return
	}
	var mod float32

	if rand.Int()%2 == 0 {
		mod = 0.1
	} else {
		mod = -0.1
	}

	if render.Scale().X+mod >= 15 || render.Scale().X+mod <= 1 {
		mod *= -1
	}

	newScale := render.Scale()
	newScale.AddScalar(mod)
	render.SetScale(newScale)
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello Demo",
		Width:  1024,
		Height: 640,
		DefaultCloseAction: true,
	}
	engo.Run(opts, &GameWorld{})
}
