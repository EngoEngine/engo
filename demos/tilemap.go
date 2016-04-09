package main

import (
	"log"

	"engo.io/engo"
	"github.com/engoengine/ecs"
)

var World *GameWorld

type GameWorld struct {
	ecs.World
}

func (game *GameWorld) Preload() {
	game.New()
	engo.Files.Add("data/sheet.png")

	log.Println("Preloaded")
}

func (game *GameWorld) Setup() {
	engo.SetBackground(0x2d3739)

	game.AddSystem(&engo.RenderSystem{})

	gameMap := ecs.NewEntity([]string{"RenderSystem"})
	tilemap := engo.NewTilemap(
		[][]string{
			{"0", "2", "0"},
			{"4", "5", "1"},
			{"2", "3", "4"},
			{"5", "1", "2"}},
		engo.Files.Image("sheet"), 16)

	mapRender := engo.NewRenderComponent(tilemap, engo.Point{1, 1}, "map")
	mapSpace := &engo.SpaceComponent{engo.Point{100, 100}, 0, 0}
	gameMap.AddComponent(mapRender)
	gameMap.AddComponent(mapSpace)

	err := game.AddEntity(gameMap)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	World = &GameWorld{}
	engo.Run("Hello", 1024, 640, false, World)
}
