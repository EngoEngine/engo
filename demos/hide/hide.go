package main

import (
	"image/color"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
)

type DefaultScene struct{}

type Rock struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

func (*DefaultScene) Preload() {
	engo.Files.AddFromDir("assets", false)
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&HideSystem{})

	// Retrieve a texture
	texture := engo.Files.Image("rock.png")

	// Create an entity
	rock := Rock{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	rock.RenderComponent = engo.NewRenderComponent(texture, engo.Point{8, 8}, "rock")
	rock.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * rock.RenderComponent.Scale().X,
		Height:   texture.Height() * rock.RenderComponent.Scale().Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *HideSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type hideEntity struct {
	*ecs.BasicEntity
	*engo.RenderComponent
}

type HideSystem struct {
	entities []hideEntity
}

func (h *HideSystem) Add(basic *ecs.BasicEntity, render *engo.RenderComponent) {
	h.entities = append(h.entities, hideEntity{basic, render})
}

func (h *HideSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range h.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		h.entities = append(h.entities[:delete], h.entities[delete+1:]...)
	}
}

func (h *HideSystem) Update(dt float32) {
	for _, e := range h.entities {
		if rand.Int()%10 == 0 {
			e.RenderComponent.Hidden = true
		} else {
			e.RenderComponent.Hidden = false
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Show and Hide Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
