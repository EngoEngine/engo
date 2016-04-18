package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/luxengine/math"
)

var globalGuy *ecs.Entity

type GameWorld struct{}

func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("data", false)

	log.Println("Preloaded")
}

func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&RotationSystem{})
	w.AddSystem(&engo.RenderSystem{})

	// Create an entity part of the Render
	globalGuy = ecs.NewEntity("RenderSystem")
	// Retrieve a texture
	texture := engo.Files.Image("icon.png")

	// Create RenderComponent... Set scale to 8x, give lable "guy"
	render := engo.NewRenderComponent(texture, engo.Point{8, 8})

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engo.SpaceComponent{
		Position: engo.Point{400, 400},
		Width:    width,
		Height:   height,
	}

	globalGuy.AddComponent(render)
	globalGuy.AddComponent(space)

	err := w.AddEntity(globalGuy)
	if err != nil {
		log.Println(err)
	}
}

type RotationSystem struct{}

func (*RotationSystem) Type() string             { return "RotationSytem" }
func (*RotationSystem) Priority() int            { return 0 }
func (*RotationSystem) New(*ecs.World)           {}
func (*RotationSystem) AddEntity(*ecs.Entity)    {}
func (*RotationSystem) RemoveEntity(*ecs.Entity) {}

func (*RotationSystem) Update(dt float32) {
	var space *engo.SpaceComponent
	if !globalGuy.Component(&space) {
		return
	}

	// speed in radians per second
	var speed float32 = math.Pi
	// speed in degrees per second
	var speedDegrees float32 = speed * (360 / (2 * math.Pi))

	space.Rotation += speedDegrees * dt
	space.Rotation = math.Mod(space.Rotation, 360)
}

func (*GameWorld) Hide()        {}
func (*GameWorld) Show()        {}
func (*GameWorld) Exit()        {}
func (*GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Rotation Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &GameWorld{})
}
