# Isometric Adventure Demo

## What does it do?
* It demonstrates how to load and use an isometric tmx resource
* Sets up a ControlSystem that moves your character on the map

## What are important aspects of the code?
These lines are key in this demo:

# Move only on the isometric TileMap

This checks if the isometric tile exists at the point and moves only if it does.

```go
var t *common.Tile
if e.SpaceComponent.Position.X >= 0 {
  t = levelData.GetTile(e.SpaceComponent.Center())
} else {
  t = levelData.GetTile(engo.Point{
    X: e.SpaceComponent.Position.X - float32(levelData.TileWidth),
    Y: e.SpaceComponent.Position.Y + float32(levelData.TileHeight),
  })
}

if t == nil {
  e.SpaceComponent.Position = prev
}
```
