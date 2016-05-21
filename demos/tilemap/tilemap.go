package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type GameWorld struct{}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (game *GameWorld) Preload() {
	// A tmx file can be generated from the Tiled Map Editor.
	// The engo tmx loader only accepts tmx files that are base64 encoded and compressed with zlib.
	// When you add tilesets to the Tiled Editor, the location where you added them from is where the engo loader will look for them
	// Tileset from : http://opengameart.org

	if err := engo.Files.LoadMany("example.tmx"); err != nil {
		panic(err)
	}
}

func (game *GameWorld) Setup(w *ecs.World) {
	common.SetBackground(color.RGBA{0x00, 0x00, 0x00, 0x00})

	w.AddSystem(&common.RenderSystem{})

	resource, err := engo.Files.Resource("example.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	// Create render and space components for each of the tiles
	tileComponents := make([]*Tile, 0)
	for _, v := range levelData.Tiles {
		if v.Image != nil {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent = common.RenderComponent{
				Drawable: v,
				Scale:    engo.Point{1, 1},
			}
			tile.SpaceComponent = common.SpaceComponent{
				Position: v.Point,
				Width:    0,
				Height:   0,
			}
			tileComponents = append(tileComponents, tile)
		}
	}
	// Do the same the levels images
	for _, v := range levelData.Images {
		if v.Image != nil {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent = common.RenderComponent{
				Drawable: v,
				Scale:    engo.Point{1, 1},
			}
			tile.SpaceComponent = common.SpaceComponent{
				Position: v.Point,
				Width:    0,
				Height:   0,
			}
			tileComponents = append(tileComponents, tile)
		}
	}

	// Add each of the tiles entities and its components to the render system
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tileComponents {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}

		}
	}

}

func (game *GameWorld) Exit()        {}
func (game *GameWorld) Hide()        {}
func (game *GameWorld) Show()        {}
func (game *GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:         "TileMap Demo",
		Width:         800,
		Height:        800,
		ScaleOnResize: false,
	}
	engo.Run(opts, &GameWorld{})
}
