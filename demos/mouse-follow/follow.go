//+build demo

package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	engo.Files.Load("icon.png")
}

func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&FollowSystem{})

	// Retrieve a texture
	texture, err := common.LoadedSprite("icon.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy := Player{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *FollowSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type followEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}

type FollowSystem struct {
	entities []followEntity
}

func (s *FollowSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent) {
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
