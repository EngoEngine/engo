//go:build demo
// +build demo

package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type DefaultScene struct{}

type Guy struct {
	ecs.BasicEntity
	common.MouseComponent
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	err := engo.Files.Load("guy.png")
	if err != nil {
		log.Println(err)
	}
}

func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&ControlSystem{})

	// These are not required, but allow you to move / rotate and still see that it works
	w.AddSystem(&common.MouseZoomer{-0.125})
	w.AddSystem(common.NewKeyboardScroller(500, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis))
	w.AddSystem(&common.MouseRotator{RotationSpeed: 0.125})

	// Retrieve a texture
	texture, err := common.LoadedSprite("guy.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{200, 200},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
		Rotation: 90,
	}
	// guy.MouseComponent doesn't have to be set, because its default values will do

	// Add our guy to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent, &guy.SpaceComponent, &guy.RenderComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type controlEntity struct {
	*ecs.BasicEntity
	*common.MouseComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, mouse *common.MouseComponent) {
	c.entities = append(c.entities, controlEntity{basic, mouse})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ControlSystem) Update(float32) {
	for _, e := range c.entities {
		if e.MouseComponent.Enter {
			engo.SetCursor(engo.CursorHand)
		} else if e.MouseComponent.Leave {
			engo.SetCursor(engo.CursorNone)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Mouse Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
