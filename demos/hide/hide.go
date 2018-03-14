//+build demo

package main

import (
	"image/color"
	"log"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

type Rock struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*DefaultScene) Preload() {
	err := engo.Files.Load("rock.png")
	if err != nil {
		log.Println(err)
	}
}

func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&HideSystem{})

	// Retrieve a texture
	texture, err := common.LoadedSprite("rock.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	rock := Rock{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	rock.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	rock.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * rock.RenderComponent.Scale.X,
		Height:   texture.Height() * rock.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *HideSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

type hideEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
}

type HideSystem struct {
	entities []hideEntity
}

func (h *HideSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent) {
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
