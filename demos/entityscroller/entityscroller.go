package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type GameWorld struct{}

type Character struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type ControlSystem struct {
	entity *Character
}

func (c *ControlSystem) Add(char *Character) {
	c.entity = char
}

func (c *ControlSystem) Remove(ecs.BasicEntity) {
	c.entity = nil
}

func (c *ControlSystem) Update(dt float32) {
	if engo.Input.Button("moveup").Down() {
		c.entity.SpaceComponent.Position.Y -= 5
	}
	if engo.Input.Button("movedown").Down() {
		c.entity.SpaceComponent.Position.Y += 5
	}
	if engo.Input.Button("moveleft").Down() {
		c.entity.SpaceComponent.Position.X -= 5
	}
	if engo.Input.Button("moveright").Down() {
		c.entity.SpaceComponent.Position.X += 5
	}
}

func (game *GameWorld) Preload() {
	// A tmx file can be generated from the Tiled Map Editor.
	// The engo tmx loader only accepts tmx files that are base64 encoded and compressed with zlib.
	// When you add tilesets to the Tiled Editor, the location where you added them from is where the engo loader will look for them
	// Tileset from : http://opengameart.org

	if err := engo.Files.LoadMany("example.tmx", "icon.png"); err != nil {
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

	character := Character{BasicEntity: ecs.NewBasic()}
	characterTexture, err := common.PreloadedSpriteSingle("icon.png")
	if err != nil {
		panic(err)
	}
	character.RenderComponent = common.RenderComponent{
		Drawable: characterTexture,
		Scale:    engo.Point{5, 5},
	}
	character.RenderComponent.SetZIndex(1)
	character.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{engo.CanvasWidth() / 2, engo.CanvasHeight() / 2},
		Width:    characterTexture.Width() * 5,
		Height:   characterTexture.Height() * 5,
	}
	// Add each of the tiles entities and its components to the render system along with the character
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&character.BasicEntity, &character.RenderComponent, &character.SpaceComponent)
			for _, v := range tileComponents {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}

		}
	}

	w.AddSystem(&ControlSystem{&character})

	// Add the EntityScroller system which contains the space component of the character and is bounded to the tmx level dimensions
	w.AddSystem(&common.EntityScroller{SpaceComponent: &character.SpaceComponent, TrackingBounds: levelData.Bounds()})
	engo.Input.RegisterButton("moveup", engo.ArrowUp)
	engo.Input.RegisterButton("moveleft", engo.ArrowLeft)
	engo.Input.RegisterButton("moveright", engo.ArrowRight)
	engo.Input.RegisterButton("movedown", engo.ArrowDown)
}

func (game *GameWorld) Exit()        {}
func (game *GameWorld) Hide()        {}
func (game *GameWorld) Show()        {}
func (game *GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:         "EntityScroller Demo",
		Width:         500,
		Height:        500,
		ScaleOnResize: false,
	}
	engo.Run(opts, &GameWorld{})
}
