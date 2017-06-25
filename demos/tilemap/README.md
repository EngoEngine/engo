# TileMap Demo

## What does it do?
It demonstrates how one can load and render a TMX file created from the TileMap Editor  

## What are important aspects of the code?
These lines are key in this demo:

* 'resource, err := engo.Files.Resource("example.tmx")' to retrieve the resource
* 'tmxResource := resource.(common.TMXResource)'to cast the resource to a tmx resource
* 'levelData := tmxResource.Level' to retrieve the level from the tmx resource

## Add render and space components to each tile in each tile layer
```go
    tileComponents := make([]*Tile, 0)

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

                tileComponents = append(tileComponents, tile)
            }
        }
    }
```

## Add render and space components to each image in each image layer
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

                tileComponents = append(tileComponents, tile)
            }
        }
    }
```

## Access object layers and do something with its regular and polyline objects
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
