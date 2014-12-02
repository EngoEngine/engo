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
	engi.Files.Add("rock", "data/rock.png")
	engi.Files.Add("sheet", "data/sheet.png")
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&MovingSystem{})
	game.AddSystem(&engi.CollisionSystem{})

	entity := engi.NewEntity([]string{"RenderSystem", "MovingSystem", "CollisionSystem"})
	texture := engi.Files.Image("bot")
	render := engi.NewRenderComponent(texture, engi.Point{1, 1}, "bot")
	space := engi.SpaceComponent{Position: engi.Point{10, 10}, Width: texture.Width() * render.Scale.X, Height: texture.Height() * render.Scale.Y}
	entity.AddComponent(&render)
	entity.AddComponent(&space)
	entity.AddComponent(&engi.CollisionMasterComponent{})
	game.AddEntity(entity)

	text := engi.NewEntity([]string{"RenderSystem"})
	textTexture := engi.NewText("Hello World", engi.NewGridFont(engi.Files.Image("font"), 20, 20))
	textRender := engi.NewRenderComponent(textTexture, engi.Point{1, 1}, "yolo?")
	textSpace := engi.SpaceComponent{engi.Point{100, 100}, textTexture.Width(), textTexture.Height()}

	text.AddComponent(&textRender)
	text.AddComponent(&textSpace)
	game.AddEntity(text)

	gameMap := engi.NewEntity([]string{"RenderSystem", "CollisionSystem"})
	tilemap := engi.NewTilemap([][]string{{"1", "2", "1"}, {"0", "0", "2"}, {"0", "0", "1"}, {"2", "1", "2"}}, engi.Files.Image("sheet"))
	mapRender := engi.NewRenderComponent(tilemap, engi.Point{1, 1}, "map")
	mapSpace := engi.SpaceComponent{engi.Point{100, 100}, 0, 0}
	gameMap.AddComponent(&mapRender)
	gameMap.AddComponent(&mapSpace)

	game.AddEntity(gameMap)

	engi.Cam.FollowEntity(entity)
}

type MovingSystem struct {
	*engi.System
}

func (ms *MovingSystem) New() {
	ms.System = &engi.System{}
}

var vel float32

func (ms *MovingSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.GetComponent(&space) {
		return
	}

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

	if engi.Keys.KEY_SPACE.JustPressed() {
		entity.Exists = false
	}
}

func (ms MovingSystem) Name() string {
	return "MovingSystem"
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
