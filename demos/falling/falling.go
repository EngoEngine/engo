package main

import (
	"image/color"
	"log"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
)

type Game struct{}

func (game *Game) Preload() {
	// Add all the files in the data directory non recursively
	engo.Files.AddFromDir("data", false)
}

func (game *Game) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	// Add all of the systems
	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.CollisionSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&FallingSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&RockSpawnSystem{})

	// Create new entity subscribed to all the systems!
	guy := ecs.NewEntity("RenderSystem", "ControlSystem", "RockSpawnSystem", "CollisionSystem", "DeathSystem")
	texture := engo.Files.Image("icon.png")
	render := engo.NewRenderComponent(texture, engo.Point{4, 4})
	// Tell the collision system that this player is solid
	collision := &engo.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := &engo.SpaceComponent{
		Position: engo.Point{(engo.Width() - width) / 2, (engo.Height() - height) / 2},
		Width:    width,
		Height:   height,
	}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	err := w.AddEntity(guy)
	if err != nil {
		log.Println(err)
	}
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Exit()        {}
func (*Game) Type() string { return "Game" }

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }
func (*ControlSystem) Pre()         {}
func (*ControlSystem) Post()        {}

func (control *ControlSystem) New(*ecs.World) {}

func (control *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var space *engo.SpaceComponent
	if !entity.Component(&space) {
		return
	}

	speed := 400 * dt

	if engo.Keys.Get(engo.A).Down() {
		space.Position.X -= speed
	}

	if engo.Keys.Get(engo.D).Down() {
		space.Position.X += speed
	}

	if engo.Keys.Get(engo.W).Down() {
		space.Position.Y -= speed
	}

	if engo.Keys.Get(engo.S).Down() {
		space.Position.Y += speed
	}
}

type RockSpawnSystem struct {
	ecs.LinearSystem

	world *ecs.World
}

func (*RockSpawnSystem) Type() string { return "RockSpawnSystem" }
func (*RockSpawnSystem) Pre()         {}
func (*RockSpawnSystem) Post()        {}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w
}

func (rock *RockSpawnSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	// 4% change of spawning a rock each frame
	if rand.Float32() < .96 {
		return
	}

	position := engo.Point{0, -32}
	position.X = rand.Float32() * (engo.Width())
	err := rock.world.AddEntity(NewRock(position))
	if err != nil {
		log.Println(err)
	}
}

func NewRock(position engo.Point) *ecs.Entity {
	rock := ecs.NewEntity("RenderSystem", "FallingSystem", "CollisionSystem")

	texture := engo.Files.Image("rock.png")
	render := engo.NewRenderComponent(texture, engo.Point{4, 4})
	space := &engo.SpaceComponent{
		Position: position,
		Width:    texture.Width() * render.Scale.X,
		Height:   texture.Height() * render.Scale.Y,
	}
	collision := &engo.CollisionComponent{Solid: true}

	rock.AddComponent(render)
	rock.AddComponent(space)
	rock.AddComponent(collision)

	return rock
}

type FallingSystem struct {
	ecs.LinearSystem
}

func (*FallingSystem) Type() string { return "FallingSystem" }

func (fs *FallingSystem) New(*ecs.World) {}

func (fs *FallingSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var space *engo.SpaceComponent
	if !entity.Component(&space) {
		return
	}
	space.Position.Y += 200 * dt
}

type DeathSystem struct {
	ecs.LinearSystem
}

func (*DeathSystem) Type() string { return "DeathSystem" }

func (ds *DeathSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		collision, isCollision := message.(engo.CollisionMessage)

		if isCollision {
			log.Println(collision, message)
			log.Println("DEAD")
		}
	})
}

func (fs *DeathSystem) UpdateEntity(entity *ecs.Entity, dt float32) {}

func (fs *DeathSystem) Receive(message engo.Message) {
	collision, isCollision := message.(engo.CollisionMessage)
	if isCollision {
		log.Println(collision, message)
		log.Println("DEAD")
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Falling Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &Game{})
}
