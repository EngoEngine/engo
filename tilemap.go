package engi

import ()

type Tilemap struct {
	Tiles [][]Tile
}

type Tile struct {
	Point
	Image Drawable
}

func NewTilemap() *Tilemap {
	tilemap := Tilemap{}
	// size := Point{10, 10}
	mapString := [][]string{{"1", "1", "1"}, {"1", "0", "1"}, {"1", "1", "1"}}
	position := Point{}
	tilesize := 16

	tilemap.Tiles = make([][]Tile, len(mapString))
	for i := range tilemap.Tiles {
		tilemap.Tiles[i] = make([]Tile, len(mapString[0]))
	}

	for y, slice := range mapString {
		for x, _ := range slice {
			tile := Tile{Point: Point{position.X + float32(x*tilesize), position.Y + float32(y*tilesize)}, Image: Files.Image("bot")}
			tilemap.Tiles[y][x] = tile
		}
	}

	return &tilemap
}
