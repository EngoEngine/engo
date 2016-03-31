package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

type Game struct{}

func (game *Game) Preload() {
	// Add all the files in the data directory non recursively
	engi.Files.AddFromDir("data", false)
}

func (game *Game) Setup(w *ecs.World) {
	engi.SetBg(color.White)

	// Add all of the systems
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.CollisionSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&FallingSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&RockSpawnSystem{})

	// Create new entity subscribed to all the systems!
	guy := ecs.NewEntity([]string{"RenderSystem", "ControlSystem", "RockSpawnSystem", "CollisionSystem", "DeathSystem"})
	texture := engi.Files.Image("icon.png")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "guy")
	// Tell the collision system that this player is solid
	collision := &engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{
		Position: engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2},
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
func (*Game) Type() string { return "Game" }

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }
func (*ControlSystem) Pre()         {}
func (*ControlSystem) Post()        {}

func (control *ControlSystem) New(*ecs.World) {}

func (control *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.Component(&space) {
		return
	}

	speed := 400 * dt

	if engi.Keys.Get(engi.A).Down() {
		space.Position.X -= speed
	}

	if engi.Keys.Get(engi.D).Down() {
		space.Position.X += speed
	}

	if engi.Keys.Get(engi.W).Down() {
		space.Position.Y -= speed
	}

	if engi.Keys.Get(engi.S).Down() {
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

	position := engi.Point{0, -32}
	position.X = rand.Float32() * (engi.Width())
	err := rock.world.AddEntity(NewRock(position))
	if err != nil {
		log.Println(err)
	}
}

func NewRock(position engi.Point) *ecs.Entity {
	rock := ecs.NewEntity([]string{"RenderSystem", "FallingSystem", "CollisionSystem"})

	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "rock")
	space := &engi.SpaceComponent{
		Position: position,
		Width:    texture.Width() * render.Scale().X,
		Height:   texture.Height() * render.Scale().Y,
	}
	collision := &engi.CollisionComponent{Solid: true}

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
	var space *engi.SpaceComponent
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
	engi.Mailbox.Listen("CollisionMessage", func(message engi.Message) {
		collision, isCollision := message.(engi.CollisionMessage)
		if isCollision {
			log.Println(collision, message)
			log.Println("DEAD")
		}
	})
}

func (fs *DeathSystem) UpdateEntity(entity *ecs.Entity, dt float32) {}

func (fs *DeathSystem) Receive(message engi.Message) {
	collision, isCollision := message.(engi.CollisionMessage)
	if isCollision {
		log.Println(collision, message)
		log.Println("DEAD")
	}
}

func main() {
	opts := engi.RunOptions{
		Title:  "Falling Demo",
		Width:  1024,
		Height: 640,
	}
	engi.Run(opts, &Game{})
}
