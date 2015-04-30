package engi

type Tile struct {
	Point
	Image *Region
}

type Tilemap struct {
	Tiles    []*Tile
	TileSize int
}

func NewTilemap(texture *Texture, mapping []uint32, tileSize int, mapHeight, mapWidth int) *Tilemap {

	// create tile map by splitting up the texture into pieces
	setWidth := int(texture.Width()) / tileSize
	setHeight := int(texture.Height()) / tileSize
	totalTiles := setWidth * setHeight

	tileset := make([]*Tile, totalTiles)

	for y := 0; y < setHeight; y++ {
		for x := 0; x < setWidth; x++ {
			t := &Tile{}
			t.Image = NewRegion(texture, x*tileSize, y*tileSize, tileSize, tileSize)
			idx := x + y*setWidth
			tileset[idx] = t
		}
	}

	tilemap := make([]*Tile, mapHeight*mapWidth)

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			idx := x + y*mapWidth
			t := &Tile{}
			if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
				t.Image = tileset[tileIdx].Image
				t.Point = Point{float32(x * tileSize), float32(y * tileSize)}
			}
			tilemap[idx] = t
		}
	}

	return &Tilemap{tilemap, tileSize}
}

func getRegionOfSpriteSheet(texture *Texture, tilesize int, index int) *Region {
	width := texture.Width()
	widthInSprites := width / float32(tilesize)

	pointer := Point{}
	step := 0
	for step != (index) {
		step += 1
		if pointer.X < (widthInSprites - 1) {
			pointer.X += 1
		} else {
			pointer.X = 0
			pointer.Y += 1
		}
	}

	pointer.MultiplyScalar(float32(tilesize))

	return NewRegion(texture, int(pointer.X), int(pointer.Y), tilesize, tilesize)
}
