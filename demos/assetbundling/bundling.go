//+build demo

package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/demos/assetbundling/assets"
)

const tilemapURL string = "example.tmx"

type Tile struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

const GameSceneType string = "GameScene"

type GameScene struct {
	files []string
}

func (s *GameScene) Preload() {
	s.files = []string{
		tilemapURL,
		"grass-tiles-2-small.png",
		"tree2-final.png",
	}
	for _, file := range s.files {
		data, err := assets.Asset(file)
		if err != nil {
			log.Fatalf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Unable to load asset with URL: %v\n At %v", file, s.Type())
		}
	}
}

func (s *GameScene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w := u.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})

	resource, err := engo.Files.Resource(tilemapURL)
	if err != nil {
		panic(err)
	}
	levelData := resource.(common.TMXResource).Level

	tileComponents := []*Tile{}

	for idx, layer := range levelData.TileLayers {
		for _, tileElement := range layer.Tiles {
			if tileElement.Image == nil {
				log.Printf("Tile is lacking image at point: %v", tileElement.Point)
			}
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent = common.RenderComponent{
				Drawable:    tileElement.Image,
				Scale:       engo.Point{X: 1, Y: 1},
				StartZIndex: float32(idx),
			}
			tile.SpaceComponent = common.SpaceComponent{
				Position: tileElement.Point,
			}
			tileComponents = append(tileComponents, tile)
		}
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tileComponents {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

func (s *GameScene) Type() string {
	return GameSceneType
}

func main() {
	opts := engo.RunOptions{
		Title:  "Assets included!",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &GameScene{})
}
