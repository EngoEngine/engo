package main

import (
	"github.com/paked/engi"
	"log"
)

var World *GameWorld

type GameWorld struct {
	batch *engi.Batch
	engi.World
}

func (game *GameWorld) Preload() {
	engi.Files.Add("bot", "data/icon.png")
	engi.Files.Add("font", "data/font.png")
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&RenderSystem{})
	game.AddSystem(&MovingSystem{})
	game.AddSystem(&engi.CollisionSystem{})

	entity := engi.NewEntity([]string{"RenderSystem", "MovingSystem", "CollisionSystem"})
	texture := engi.Files.Image("bot")
	render := engi.NewRenderComponent(texture, engi.Point{1, 1}, "bot")
	space := engi.SpaceComponent{Position: engi.Point{10, 10}, Width: texture.Width() * render.Scale.X, Height: texture.Height() * render.Scale.Y}
	collisionMaster := engi.CollisionMasterComponent{}
	entity.AddComponent(&render)
	entity.AddComponent(&space)
	entity.AddComponent(&collisionMaster)
	game.AddEntity(entity)

	entity3 := engi.NewEntity([]string{"RenderSystem", "CollisionSystem"})
	render3 := engi.NewRenderComponent(texture, engi.Point{10, 10}, "bigbot")
	space3 := engi.SpaceComponent{Position: engi.Point{100, 100}, Width: texture.Width() * render3.Scale.X, Height: texture.Height() * render3.Scale.Y}
	entity3.AddComponent(&render3)
	entity3.AddComponent(&space3)
	game.AddEntity(entity3)

	entityTwo := engi.NewEntity([]string{"RenderSystem"})
	componentTwo := engi.NewRenderComponent(engi.NewGridFont(engi.Files.Image("font"), 20, 20), engi.Point{1, 1}, "wut.")
	space2 := engi.SpaceComponent{Position: engi.Point{500, 100}, Width: 100, Height: 100}
	entityTwo.AddComponent(&componentTwo)
	entityTwo.AddComponent(&space2)
	game.AddEntity(entityTwo)
	log.Println("Setup")
}

type MovingSystem struct {
	*engi.System
}

func (ms *MovingSystem) New() {
	ms.System = &engi.System{}
}

var vel float32

func (ms *MovingSystem) Update(entity *engi.Entity, dt float32) {
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if hasSpace {
		vel = 200 * dt
		if engi.Keys.KEY_D.Down() {
			space.Position.X += vel
		}

		if engi.Keys.KEY_A.Down() {
			space.Position.X -= vel
		}

		if engi.Keys.KEY_W.Down() {
			space.Position.Y -= vel
		}

		if engi.Keys.KEY_S.Down() {
			space.Position.Y += vel
		}
	}
}

func (ms MovingSystem) Name() string {
	return "MovingSystem"
}

type RenderSystem struct {
	*engi.System
}

func (rs *RenderSystem) New() {
	rs.System = &engi.System{}
}

func (rs RenderSystem) Pre() {
	engi.Gl.Clear(engi.Gl.COLOR_BUFFER_BIT)
	World.batch.Begin()
}

func (rs RenderSystem) Post() {
	World.batch.End()
}

func (rs *RenderSystem) Update(entity *engi.Entity, dt float32) {
	render, hasRender := entity.GetComponent("RenderComponent").(*engi.RenderComponent)
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if hasRender && hasSpace {
		switch render.Display.(type) {
		case engi.Drawable:
			drawable := render.Display.(engi.Drawable)
			World.batch.Draw(drawable, space.Position.X, space.Position.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, 0xffffff, 1)
		case *engi.Font:
			font := render.Display.(*engi.Font)
			font.Print(World.batch, render.Label, space.Position.X, space.Position.Y, 0xffffff)
		}
	}
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
