package main

import (
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
	engi.SetBg(0x2d3739)

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

	space := &engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	w.AddEntity(guy)
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }

type ControlSystem struct {
	*ecs.System
}

func (control *ControlSystem) New(*ecs.World) {
	control.System = ecs.NewSystem()
}

func (ControlSystem) Type() string {
	return "ControlSystem"
}

func (control *ControlSystem) Update(entity *ecs.Entity, dt float32) {
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
	*ecs.System

	world *ecs.World
}

func (RockSpawnSystem) Type() string {
	return "RockSpawnSystem"
}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.System = ecs.NewSystem()
	rock.world = w
}

func (rock *RockSpawnSystem) Update(entity *ecs.Entity, dt float32) {
	// 4% change of spawning a rock each frame
	if rand.Float32() < .96 {
		return
	}

	position := engi.Point{0, -32}
	position.X = rand.Float32() * (engi.Width())
	rock.world.AddEntity(NewRock(position))
}

func NewRock(position engi.Point) *ecs.Entity {
	rock := ecs.NewEntity([]string{"RenderSystem", "FallingSystem", "CollisionSystem", "SpeedSystem"})

	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "rock")
	space := &engi.SpaceComponent{position, texture.Width() * render.Scale().X, texture.Height() * render.Scale().Y}
	collision := &engi.CollisionComponent{Solid: true}

	rock.AddComponent(render)
	rock.AddComponent(space)
	rock.AddComponent(collision)

	return rock
}

type FallingSystem struct {
	*ecs.System
}

func (fs *FallingSystem) New(*ecs.World) {
	fs.System = ecs.NewSystem()
	//engi.Mailbox.Listen("CollisionMessage", fs)

}

func (FallingSystem) Type() string {
	return "FallingSystem"
}

func (fs *FallingSystem) Update(entity *ecs.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.Component(&space) {
		return
	}
	space.Position.Y += 200 * dt
}

type DeathSystem struct {
	*ecs.System
}

func (ds *DeathSystem) New(*ecs.World) {
	ds.System = ecs.NewSystem()
	// Subscribe to ScoreMessage
	engi.Mailbox.Listen("ScoreMessage", func(message engi.Message) {
		collision, isCollision := message.(engi.CollisionMessage)
		if isCollision {
			log.Println(collision, message)
			log.Println("DEAD")
		}
	})

}

func (DeathSystem) Type() string {
	return "DeathSystem"
}

func (fs *DeathSystem) Update(entity *ecs.Entity, dt float32) {

}

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
	engi.Open(opts, &Game{})
}
