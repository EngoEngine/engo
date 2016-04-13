package main

import (
	"image/color"
	"log"
	"math/rand"

	"engo.io/engo"
	"engo.io/ecs"
)

type GameWorld struct{}

func (game *GameWorld) Preload() {
	engo.Files.AddFromDir("assets", false)
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&HideSystem{})

	guy := ecs.NewEntity("RenderSystem", "HideSystem")
	texture := engo.Files.Image("rock.png")
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

type HideSystem struct {
	ecs.LinearSystem
}

func (*HideSystem) Type() string { return "HideSystem" }

func (s *HideSystem) New(*ecs.World) {}

func (c *HideSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var render *engo.RenderComponent
	if !entity.Component(&render) {
		return
	}
	if rand.Int()%10 == 0 {
		render.SetPriority(engo.Hidden)
	} else {
		render.SetPriority(engo.MiddleGround)
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Show and Hide Demo",
		Width:  1024,
		Height: 640,
		
	}
	engo.Run(opts, &GameWorld{})
}
