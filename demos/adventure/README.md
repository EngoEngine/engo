# Adventure Demo

## What does it do?
* It applies various different things from multiple engo demos
* It demonstrates how one can load and render a TMX file created from the TileMap Editor
* It shows how to add a character entity
  * It demonstrates how to animate the character and move it around
* It shows how to access different layers in the TileMap
  * It demonstrates how to set different Z indices to create depth

## What are important aspects of the code?
These lines are key in this demo:

* `levelData := tmxResource.Level` to retrieve the level from the tmx resource
* `hero := scene.CreateHero( ...` to create a character instance
* `hero.ControlComponent = ControlComponent{ ...` to add a control component to the character
* `for _, tileLayer := range levelData.TileLayers { ...` to access different tile layer of the TMX map
* `for _, imageLayer := range levelData.ImageLayers { ...` to access different image layers of the TMX map
* `for _, objectLayer := range levelData.ObjectLayers { ...` to access different object layers of the TMX map
* `tile.RenderComponent.SetZIndex(3)` to add a Z index and 'depth' to one layer



# Add render and space components to each tile in each tile layer
```go
for _, tileLayer := range levelData.TileLayers {
    for _, tileElement := range tileLayer.Tiles {

        if tileElement.Image != nil {
            tile := &Tile{BasicEntity: ecs.NewBasic()}
            tile.RenderComponent = common.RenderComponent{
                Drawable: tileElement,
                Scale:    engo.Point{1, 1},
            }
            tile.SpaceComponent = common.SpaceComponent{
                Position: tileElement.Point,
                Width:    0,
                Height:   0,
            }

            if tileLayer.Name == "grass" {
                tile.RenderComponent.SetZIndex(0)
            }

            if tileLayer.Name == "trees" {
                tile.RenderComponent.SetZIndex(2)
            }

            tileComponents = append(tileComponents, tile)
        }
    }
}
```

# Add render and space components to each image in each image layer
```go
for _, imageLayer := range levelData.ImageLayers {
    for _, imageElement := range imageLayer.Images {

        if imageElement.Image != nil {
            tile := &Tile{BasicEntity: ecs.NewBasic()}
            tile.RenderComponent = common.RenderComponent{
                Drawable: imageElement,
                Scale:    engo.Point{1, 1},
            }
            tile.SpaceComponent = common.SpaceComponent{
                Position: imageElement.Point,
                Width:    0,
                Height:   0,
            }

            if imageLayer.Name == "clouds" {
                tile.RenderComponent.SetZIndex(3)
            }

            tileComponents = append(tileComponents, tile)
        }
    }
}
```

# Access object layers and do something with its regular and polyline objects
```go
// Access Object Layers
for _, objectLayer := range levelData.ObjectLayers {
    log.Println("This object layer is called " + objectLayer.Name)
    // Do something with every regular Object
    for _, object := range objectLayer.Objects {
        log.Println("This object is called " + object.Name)
    }

    // Do something with every polyline Object
    for _, polylineObject := range objectLayer.PolyObjects {
        log.Println("This object is called " + polylineObject.Name)
    }
}
```

# A Pause Menu

A pause system handles pausing

```go
type pauseEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*common.SpaceComponent
	*common.RenderComponent
	*ControlComponent
	*SpeedComponent
}
```

Notice the pauseEntity contains a lot of components! This is likely going to be the case!
Luckily, they're poointers so we can use nil!

```go
type PauseSystem struct {
	entities []pauseEntity
	world    *ecs.World
	paused   bool
}

func (p *PauseSystem) New(w *ecs.World) {
	p.world = w
}
```

It'll also need a New function so it can get access to the world.

```go
func (p *PauseSystem) Update(dt float32) {
	if engo.Input.Button(pauseButton).JustPressed() {
		if !p.paused {
			for _, system := range p.world.Systems() {
				switch sys := system.(type) {
				case *common.AnimationSystem:
					for _, ent := range p.entities {
						sys.Remove(*ent.BasicEntity)
					}
				case *SpeedSystem:
					for _, ent := range p.entities {
						sys.Remove(*ent.BasicEntity)
					}
				case *ControlSystem:
					for _, ent := range p.entities {
						sys.Remove(*ent.BasicEntity)
					}
				}
			}
		} else {
			for _, system := range p.world.Systems() {
				switch sys := system.(type) {
				case *common.AnimationSystem:
					for _, ent := range p.entities {
						if ent.AnimationComponent != nil {
							sys.Add(
								ent.BasicEntity,
								ent.AnimationComponent,
								ent.RenderComponent,
							)
						}
					}
				case *SpeedSystem:
					for _, ent := range p.entities {
						if ent.SpeedComponent != nil {
							sys.Add(
								ent.BasicEntity,
								ent.SpeedComponent,
								ent.SpaceComponent,
							)
						}
					}
				case *ControlSystem:
					for _, ent := range p.entities {
						if ent.ControlComponent != nil {
							sys.Add(
								ent.BasicEntity,
								ent.AnimationComponent,
								ent.ControlComponent,
								ent.SpaceComponent,
							)
						}
					}
				}
			}
		}
		p.paused = !p.paused
	}
}
```

Then in the update you can handle removing the entities from the systems that make
them move. This just removes those from the system. If you wanted to do something
like show a menu or a grey overlay, here's the place to do it.
