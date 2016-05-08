package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

type DefaultScene struct{}

type Player struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {
	engo.Files.Add("data/icon.png")
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&FollowSystem{})

	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create an entity
	guy := Player{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = engo.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *FollowSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type followEntity struct {
	*ecs.BasicEntity
	*engo.RenderComponent
	*engo.SpaceComponent
}

type FollowSystem struct {
	entities []followEntity
}

func (s *FollowSystem) Add(basic *ecs.BasicEntity, render *engo.RenderComponent, space *engo.SpaceComponent) {
	s.entities = append(s.entities, followEntity{basic, render, space})
}

func (s *FollowSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}

	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *FollowSystem) Update(dt float32) {
	for _, e := range s.entities {
		e.SpaceComponent.Position.X += engo.Input.Axis(engo.DefaultMouseXAxis).Value()
		e.SpaceComponent.Position.Y += engo.Input.Axis(engo.DefaultMouseYAxis).Value()
	}
}

func main() {
	opts := engo.RunOptions{
		Title:          "Follow Demo",
		Width:          1024,
		Height:         640,
		StandardInputs: true,
	}

	engo.Run(opts, &DefaultScene{})
}
