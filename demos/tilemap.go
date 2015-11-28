package main

import (
	"log"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

var World *GameWorld

type GameWorld struct {
	ecs.World
}

func (game *GameWorld) Preload() {
	game.New()
	engi.Files.Add("data/sheet.png")

	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engi.SetBg(0x2d3739)

	game.AddSystem(&engi.RenderSystem{})

	gameMap := ecs.NewEntity([]string{"RenderSystem"})
	tilemap := engi.NewTilemap(
		[][]string{
			{"0", "2", "0"},
			{"4", "5", "1"},
			{"2", "3", "4"},
			{"5", "1", "2"}},
		engi.Files.Image("sheet"), 16)

	mapRender := engi.NewRenderComponent(tilemap, engi.Point{1, 1}, "map")
	mapSpace := &engi.SpaceComponent{engi.Point{100, 100}, 0, 0}
	gameMap.AddComponent(mapRender)
	gameMap.AddComponent(mapSpace)

	game.AddEntity(gameMap)
}

func main() {
	World = &GameWorld{}
	engi.Open("Hello", 1024, 640, false, World)
}
