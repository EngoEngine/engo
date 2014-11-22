//GAME IDEAS GUSY. GIVE ME SOME STUFF - QUICK as an example kinda trhing

package main

import (
	"github.com/paked/engi"
	"log"
)

var (
	World *Game
)

type Game struct {
	engi.World
}

func (game Game) Preload() {
	engi.Files.Add("guy", "data/icon.png")
	engi.Files.Add("rock", "data/rock.png")
	engi.Files.Add("font", "data/font.go")
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&engi.CollisionSystem{})
	game.AddSystem(&ControlSystem{})
	game.AddSystem(&RockSpawnSystem{})
	game.AddSystem(&FallingSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ControlSystem", "RockSpawnSystem", "CollisionSystem"})
	texture := engi.Files.Image("guy")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "guy")
	collisionMaster := engi.CollisionMasterComponent{}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}
	guy.AddComponent(&render)
	guy.AddComponent(&space)
	guy.AddComponent(&collisionMaster)

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
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if !hasSpace {
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
	if engi.Keys.KEY_SPACE.JustPressed() {
		World.AddEntity(NewRock())
	}
}

func NewRock() *engi.Entity {
	rock := engi.NewEntity([]string{"RenderSystem", "FallingSystem", "CollisionSystem"})
	texture := engi.Files.Image("rock")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "rock")
	space := engi.SpaceComponent{engi.Point{10, 10}, texture.Width() * render.Scale.X, texture.Height() * render.Scale.Y}
	rock.AddComponent(&render)
	rock.AddComponent(&space)
	return rock
}

type FallingSystem struct {
	*engi.System
}

func (falling *FallingSystem) New() {
	falling.System = &engi.System{}
	engi.TheWorld.Mailbox.Listen("CollisionMessage", falling)
}

func (falling *FallingSystem) Name() string {
	return "FallingSystem"
}

func (falling *FallingSystem) Recieve(message engi.Message) {
	log.Println("WHDAHDHADHAHDA")
	collisonMessage, isCollisionMesage := message.(engi.CollisionMessage)
	if !isCollisionMesage {
		return
	}

	log.Println("DESTROY ENTITY", collisonMessage)
}

func (falling *FallingSystem) Update(entity *engi.Entity, dt float32) {
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if !hasSpace {
		return
	}

	space.Position.Y += 200 * dt
}

func main() {
	log.Println("[Game] Says hello, written in github.com/paked/engi + Go")
	World = &Game{}
	engi.Open("Stream Game", 800, 800, false, World)
}
