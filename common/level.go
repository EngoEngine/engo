//+build ignore

package common

import "engo.io/engo"

import (
	"engo.io/gl"
)

type Level struct {
	width      int
	height     int
	TileWidth  int
	TileHeight int
	Tiles      []*tile
	LineBounds []*Line
	Images     []*tile
}

func (t *tile) Height() float32 {
	return float32(t.height)
}

func (t *tile) Width() float32 {
	return float32(t.width)
}

func (t *tile) Texture() *gl.Texture {
	return t.Image.Texture()
}

func (t *tile) Close() {
	// noop
}

func (t *tile) View() (float32, float32, float32, float32) {
	return t.Image.u, t.Image.v, t.Image.u2, t.Image.v2
}

type tile struct {
	Point
	height int
	width  int
	Image  *Region
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
		for y := 0; y < lvl.height; y++ {
			for x := 0; x < lvl.width; x++ {
				idx := x + y*lvl.width
				t := &tile{}
				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ts[tileIdx].Image

					t.height = lvl.TileHeight
					t.width = lvl.TileWidth
					t.Point = engo.Point{float32(x * lvl.TileWidth), float32(y * lvl.TileHeight)}
				}
				tilemap = append(tilemap, t)
			}
		}
	}

	return tilemap
}
