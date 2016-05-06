package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
)

type DefaultScene struct{}

type Guy struct {
	ecs.BasicEntity

	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {
	// Load all files from the data directory. `false` means: do not do it recursively.
	engo.Files.Add("assets/icon.png")
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})

	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

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
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
