package engi

type Level struct {
	Width      int
	Height     int
	TileWidth  int
	TileHeight int
	Tiles      []*tile
	LineBounds []Line
	Images     []*tile
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

func createTileset(lvl *Level, sheets []*tilesheet) []*tile {
	tileset := make([]*tile, 0)
	tw := lvl.TileWidth
	th := lvl.TileHeight

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

func createLevelTiles(lvl *Level, layers []*layer, ts []*tile) []*tile {
	tilemap := make([]*tile, 0)

	for _, lay := range layers {
		mapping := lay.TileMapping
		for y := 0; y < lvl.Height; y++ {
			for x := 0; x < lvl.Width; x++ {
				idx := x + y*lvl.Width
				t := &tile{}
				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ts[tileIdx].Image
					t.Point = Point{float32(x * lvl.TileWidth), float32(y * lvl.TileHeight)}
				}
				tilemap = append(tilemap, t)
			}
		}
	}

	return tilemap
}

// Works for tiles rendered right-down
func regionFromSheet(sheet *Texture, tw, th int, index int) *Region {
	setWidth := int(sheet.Width()) / tw
	x := (index % setWidth) * tw
	y := (index / setWidth) * th
	return NewRegion(sheet, float32(x), float32(y), float32(tw), float32(th))
}
