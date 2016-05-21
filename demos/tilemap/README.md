# TileMap Demo

## What does it do?
It demonstrates how one can load and render a TMX file created from the TileMap Editor  

## What are important aspects of the code?
These lines are key in this demo:

* 'resource, err := engo.Files.Resource("example.tmx")' to retrieve the resource
* 'tmxResource := resource.(common.TMXResource)'to cast the resource to a tmx resource
* 'levelData := tmxResource.Level' to retrieve the level from the tmx resource

# Add render and space components to each tile
```
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
```
