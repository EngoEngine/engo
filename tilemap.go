package engi

import ()

type Tilemap struct {
	Tiles [][]Tile
}

type Tile struct {
	Point
	Image *Texture
}

func NewTilemap(mapString [][]string) *Tilemap {
	tilemap := Tilemap{}
	position := Point{}
	tilesize := 16

	tilemap.Tiles = make([][]Tile, len(mapString))
	for i := range tilemap.Tiles {
		tilemap.Tiles[i] = make([]Tile, len(mapString[0]))
	}

	for y, slice := range mapString {
		for x, key := range slice {
			var image *Texture
			switch key {
			case "1":
				image = Files.Image("bot")
			case "2":
				image = Files.Image("rock")
			}
			tile := Tile{Point: Point{position.X + float32(x*tilesize), position.Y + float32(y*tilesize)}, Image: image}
			tilemap.Tiles[y][x] = tile
		}
	}

	return &tilemap
}

func CollideTilemap(e *Entity, et *Entity, t *Tilemap) {
	var eSpace *SpaceComponent
	var tSpace *SpaceComponent

	if !e.GetComponent(&eSpace) || !et.GetComponent(&tSpace) {
		return
	}

	for _, slice := range t.Tiles {
		for _, tile := range slice {
			aabb := AABB{Point{tile.X + tSpace.Position.X, tile.Y + tSpace.Position.Y}, Point{tile.X + tSpace.Position.X + 16, tile.Y + tSpace.Position.Y + 16}}
			if IsIntersecting(eSpace.AABB(), aabb) {
				mtd := MinimumTranslation(eSpace.AABB(), aabb)
				eSpace.Position.X += mtd.X
				eSpace.Position.Y += mtd.Y
				Mailbox.Dispatch(CollisionMessage{e})
			}
		}
	}
}
