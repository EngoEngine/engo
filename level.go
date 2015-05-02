package engi

type Level struct {
	Tiles []*tile
}

type LevelData struct {
	Width      int
	Height     int
	TileWidth  int
	TileHeight int
	Tileset    []*tile
	Layers     []layer
}

type tile struct {
	Point
	Image *Region
}

type tilesheet struct {
	Image    *Texture
	Firstgid int
}

type layer struct {
	Name        string
	TileMapping []uint32
}

func createTileset(sheets []tilesheet, tw, th int) []*tile {
	tileset := make([]*tile, 0)

	for _, sheet := range sheets {
		setWidth := int(sheet.Image.Width()) / tw
		setHeight := int(sheet.Image.Height()) / th
		totalTiles := setWidth * setHeight

		for i := 0; i < totalTiles; i++ {
			t := &tile{}
			t.Image = regionFromSheet(sheet.Image, tw, th, i)
			tileset = append(tileset, t)
		}
	}

	return tileset
}

func createLevel(ld *LevelData) *Level {
	tilemap := make([]*tile, 0)

	for _, layer := range ld.Layers {
		mapping := layer.TileMapping
		for y := 0; y < ld.Height; y++ {
			for x := 0; x < ld.Width; x++ {
				idx := x + y*ld.Width
				t := &tile{}
				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ld.Tileset[tileIdx].Image
					t.Point = Point{float32(x * ld.TileWidth), float32(y * ld.TileHeight)}
				}
				tilemap = append(tilemap, t)
			}
		}
	}

	return &Level{tilemap}
}

// Works for tiles rendered right-down
func regionFromSheet(sheet *Texture, tw, th int, index int) *Region {
	setWidth := int(sheet.Width()) / tw
	x := (index % setWidth) * tw
	y := (index / setWidth) * th
	return NewRegion(sheet, x, y, tw, th)
}
