package main

import (
	"image/color"
	"math/rand"

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
	// This could be done individually: engo.Files.Add("data/icon.png"), etc
	// Second value (false) says whether to check recursively or not
	engo.Files.AddFromDir("data", false)
}

func (*DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ScaleSystem{})

	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = engo.NewRenderComponent(texture, engo.Point{8, 8}, "guy")
	guy.SpaceComponent = engo.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale().X,
		Height:   texture.Height() * guy.RenderComponent.Scale().Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engo.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *ScaleSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent)
		}
	}
}

func (*DefaultScene) Hide()        {}
func (*DefaultScene) Show()        {}
func (*DefaultScene) Exit()        {}
func (*DefaultScene) Type() string { return "GameWorld" }

type scaleEntity struct {
	*ecs.BasicEntity
	*engo.RenderComponent
}

type ScaleSystem struct {
	entities []scaleEntity
}

func (s *ScaleSystem) Add(basic *ecs.BasicEntity, render *engo.RenderComponent) {
	s.entities = append(s.entities, scaleEntity{basic, render})
}

func (s *ScaleSystem) Remove(basic ecs.BasicEntity) {
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

func (s *ScaleSystem) Update(dt float32) {
	for _, e := range s.entities {
		var mod float32

		if rand.Int()%2 == 0 {
			mod = 0.1
		} else {
			mod = -0.1
		}

		if e.RenderComponent.Scale().X+mod >= 15 || e.RenderComponent.Scale().X+mod <= 1 {
			mod *= -1
		}

		newScale := e.RenderComponent.Scale()
		newScale.AddScalar(mod)
		e.RenderComponent.SetScale(newScale)
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Scale Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
