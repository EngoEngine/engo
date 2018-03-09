//+build demo

package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/engo/math"
)

type Guy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type DefaultScene struct{}

func (game *DefaultScene) Preload() {
	err := engo.Files.Load("icon.png")
	if err != nil {
		log.Println(err)
	}
}

func (game *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&RotationSystem{})
	w.AddSystem(&common.RenderSystem{})

	// Retrieve a texture
	texture, err := common.LoadedSprite("icon.png")
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
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *RotationSystem:
			sys.Add(&guy.BasicEntity, &guy.SpaceComponent)
		}
	}
}

type rotationEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type RotationSystem struct {
	entities []rotationEntity
}

func (r *RotationSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	r.entities = append(r.entities, rotationEntity{basic, space})
}

func (r *RotationSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range r.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		r.entities = append(r.entities[:delete], r.entities[delete+1:]...)
	}
}

func (r *RotationSystem) Update(dt float32) {
	// speed in radians per second
	var speed float32 = math.Pi
	// speed in degrees per second
	var speedDegrees float32 = speed * 180 / math.Pi

	for _, e := range r.entities {
		e.SpaceComponent.Rotation += speedDegrees * dt
		e.SpaceComponent.Rotation = math.Mod(e.SpaceComponent.Rotation, 360)
	}
}

func (*DefaultScene) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Rotation Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
