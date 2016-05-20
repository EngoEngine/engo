package common

import (
	"engo.io/engo"
	"engo.io/gl"
)

type Level struct {
	width      int
	height     int
	TileWidth  int
	TileHeight int
	Tiles      []*tile
	LineBounds []*engo.Line
	Images     []*tile
}

func (t *tile) Height() float32 {
	return t.Image.Height()
}

func (t *tile) Width() float32 {
	return t.Image.Width()
}

func (t *tile) Texture() *gl.Texture {
	return t.Image.id
}

func (t *tile) Close() {
	t.Image.Close()
}

func (t *tile) View() (float32, float32, float32, float32) {
	return t.Image.View()
}

type tile struct {
	engo.Point
	Image *Texture
}

type tilesheet struct {
	Image    *TextureResource
	Firstgid int
}

type layer struct {
	Name        string
	TileMapping []uint32
}

func createTileset(lvl *Level, sheets []*tilesheet) []*tile {
	tileset := make([]*tile, 0)
	tw := float32(lvl.TileWidth)
	th := float32(lvl.TileHeight)

	for _, sheet := range sheets {
		setWidth := sheet.Image.Width / tw
		setHeight := sheet.Image.Height / th
		totalTiles := int(setWidth * setHeight)

		for i := 0; i < totalTiles; i++ {
			t := &tile{}
			x := float32(i%int(setWidth)) * tw
			y := float32(i/int(setWidth)) * th

			invTexWidth := 1.0 / float32(sheet.Image.Width)
			invTexHeight := 1.0 / float32(sheet.Image.Height)

			u := float32(x) * invTexWidth
			v := float32(y) * invTexHeight
			u2 := float32(x+tw) * invTexWidth
			v2 := float32(y+th) * invTexHeight
			t.Image = &Texture{id: sheet.Image.Texture, width: tw, height: th, viewport: engo.AABB{engo.Point{u, v}, engo.Point{u2, v2}}}
			tileset = append(tileset, t)
		}
	}

	return tileset
}

func createLevelTiles(lvl *Level, layers []*layer, ts []*tile) []*tile {
	tilemap := make([]*tile, 0)

	for _, lay := range layers {
		mapping := lay.TileMapping
		for y := 0; y < lvl.height; y++ {
			for x := 0; x < lvl.width; x++ {
				idx := x + y*lvl.width
				t := &tile{}
				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ts[tileIdx].Image
					t.Point = engo.Point{float32(x * lvl.TileWidth), float32(y * lvl.TileHeight)}

				}

				tilemap = append(tilemap, t)
			}
		}
	}

	return tilemap
}
