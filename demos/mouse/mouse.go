package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/demos/demoutils"
)

type DefaultScene struct{}

type Guy struct {
	ecs.BasicEntity
	engo.MouseComponent
	engo.RenderComponent
	engo.SpaceComponent
}

type Tracker struct {
	*demoutils.Background
	engo.MouseComponent
}

var tracker Tracker

func (*DefaultScene) Preload() {
	// Load all files from the data directory. `false` means: do not do it recursively.
	engo.Files.AddFromDir("data", false)
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.Black)

	w.AddSystem(&engo.MouseSystem{})
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&engo.MouseZoomer{-0.125})
	w.AddSystem(engo.NewKeyboardScroller(500, engo.W, engo.D, engo.S, engo.A))
	w.AddSystem(&engo.MouseRotator{RotationSpeed: 0.125})

	bg := demoutils.NewBackground(w, 1000, 1000, color.RGBA{100, 255, 100, 255}, color.RGBA{100, 200, 100, 255})
	bg.RenderComponent.SetZIndex(-1)
	bg.SpaceComponent.Position = engo.Point{250, 250}

	tracker = Tracker{Background: demoutils.NewBackground(w, 10, 10, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 255, 255})}
	tracker.RenderComponent.SetZIndex(1)
	tracker.MouseComponent.Track = true

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.MouseSystem:
			sys.Add(&tracker.BasicEntity, &tracker.MouseComponent, nil, nil)
		}
	}

	demoutils.NewBackground(w, 10, 10, color.RGBA{0, 255, 255, 255}, color.RGBA{0, 255, 255, 255})

	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = engo.NewRenderComponent(texture, engo.Point{8, 8})
	guy.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{bg.SpaceComponent.Position.X + 150, bg.SpaceComponent.Position.Y + 150},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
		Rotation: 180,
	}
	// guy.MouseComponent doesn't have to be set, because its default values will do

	// Add our guy to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *engo.MouseSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent, &guy.SpaceComponent, &guy.RenderComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type controlEntity struct {
	*ecs.BasicEntity
	*engo.MouseComponent
}

type ControlSystem struct {
	entities []controlEntity
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, mouse *engo.MouseComponent) {
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

func (c *ControlSystem) Update(dt float32) {
	tracker.Position.X, tracker.Position.Y = tracker.MouseComponent.MouseX, tracker.MouseComponent.MouseY

	for _, e := range c.entities {
		if e.MouseComponent.Enter {
			engo.SetCursor(engo.Hand)
		} else if e.MouseComponent.Leave {
			engo.SetCursor(nil)
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
