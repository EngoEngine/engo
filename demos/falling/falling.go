package main

import (
	"github.com/paked/engi"
	"log"
	"math/rand"
)

var (
	W Game
)

type Game struct {
	engi.World
}

func (game Game) Preload() {
	// Add all the files in the data directory non recursively
	engi.Files.AddFromDir("data", false)
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)

	// Add all of the systems
	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&engi.CollisionSystem{})
	game.AddSystem(&DeathSystem{})
	game.AddSystem(&FallingSystem{})
	game.AddSystem(&ControlSystem{})
	game.AddSystem(&RockSpawnSystem{})

	// Create new entity subscribed to all the systems!
	guy := engi.NewEntity([]string{"RenderSystem", "ControlSystem", "RockSpawnSystem", "CollisionSystem", "DeathSystem"})
	texture := engi.Files.Image("icon.png")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "guy")
	// Tell the collision system that this player is solid
	collision := engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(&render)
	guy.AddComponent(&space)
	guy.AddComponent(&collision)

	game.AddEntity(guy)
}

type ControlSystem struct {
	*engi.System
}

func (control *ControlSystem) New() {
	control.System = &engi.System{}
}

func (control ControlSystem) Name() string {
	return "ControlSystem"
}

func (control *ControlSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.GetComponent(&space) {
		return
	}

	speed := 400 * dt

	if engi.Keys.KEY_A.Down() {
		space.Position.X -= speed
	}

	if engi.Keys.KEY_D.Down() {
		space.Position.X += speed
	}

	if engi.Keys.KEY_W.Down() {
		space.Position.Y -= speed
	}

	if engi.Keys.KEY_S.Down() {
		space.Position.Y += speed
	}
}

type RockSpawnSystem struct {
	*engi.System
}

func (rock RockSpawnSystem) Name() string {
	return "RockSpawnSystem"
}

func (rock *RockSpawnSystem) New() {
	rock.System = &engi.System{}
}

func (rock *RockSpawnSystem) Update(entity *engi.Entity, dt float32) {
	// 4% change of spawning a rock each frame
	if rand.Float32() < .96 {
		return
	}

	position := engi.Point{0, -32}
	position.X = rand.Float32() * (engi.Width())
	W.AddEntity(NewRock(position))
}

func NewRock(position engi.Point) *engi.Entity {
	rock := engi.NewEntity([]string{"RenderSystem", "FallingSystem", "CollisionSystem", "SpeedSystem"})

	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "rock")
	space := engi.SpaceComponent{position, texture.Width() * render.Scale.X, texture.Height() * render.Scale.Y}
	collision := engi.CollisionComponent{Solid: true}

	rock.AddComponent(&render)
	rock.AddComponent(&space)
	rock.AddComponent(&collision)

	return rock
}

type FallingSystem struct {
	*engi.System
}

func (fs *FallingSystem) New() {
	fs.System = &engi.System{}
	//engi.Mailbox.Listen("CollisionMessage", fs)

}

func (fs FallingSystem) Name() string {
	return "FallingSystem"
}

func (fs FallingSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.GetComponent(&space) {
		return
	}
	space.Position.Y += 200 * dt
}

type DeathSystem struct {
	*engi.System
}

func (ds *DeathSystem) New() {
	ds.System = &engi.System{}
	// Subscribe to ScoreMessage
	engi.Mailbox.Listen("ScoreMessage", func(message engi.Message) {
		collision, isCollision := message.(engi.CollisionMessage)
		if isCollision {
			log.Println(collision, message)
			log.Println("DEAD")
		}
	})

}

func (ds DeathSystem) Name() string {
	return "DeathSystem"
}

func (fs DeathSystem) Update(entity *engi.Entity, dt float32) {

}

func (fs DeathSystem) Receive(message engi.Message) {
	collision, isCollision := message.(engi.CollisionMessage)
	if isCollision {
		log.Println(collision, message)
		log.Println("DEAD")
	}
}

func main() {
	log.Println("[Game] Says hello, written in github.com/paked/engi + Go")
	W = Game{}
	engi.Open("Stream Game", 800, 800, false, &W)
}
