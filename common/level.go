package common

import (
	"engo.io/engo"
	"engo.io/gl"
)

// Level is a parsed TMX level containing all layers and default Tiled attributes
type Level struct {
	Orientation  string
	RenderOrder  string
	width        int
	height       int
	TileWidth    int
	TileHeight   int
	NextObjectId int
	TileLayers   []*TileLayer
	ImageLayers  []*ImageLayer
	ObjectLayers []*ObjectLayer
	// old
	// Tiles      []*tile
	// LineBounds []*engo.Line
	// Images     []*tile
}

// TileLayer contains a list of its tiles plus all default Tiled attributes
type TileLayer struct {
	Name   string
	Width  int
	Height int
	Tiles  []*tile
}

// ImageLayer contains a list of its images plus all default Tiled attributes
type ImageLayer struct {
	Name   string
	Width  int
	Height int
	Source string
	Images []*tile
}

// ObjectLayer contains a list of its standard objects as well as a list of all its polyline objects
type ObjectLayer struct {
	Name        string
	OffSetX     float32
	OffSetY     float32
	Objects     []*Object
	PolyObjects []*PolylineObject
}

// Object is a standard TMX object with all its default Tiled attributes
type Object struct {
	Id     int
	Name   string
	Type   string
	X      float64
	Y      float64
	Width  int
	Height int
}

// PolylineObject is a TMX polyline object with all its default Tiled attributes
type PolylineObject struct {
	Id         int
	Name       string
	Type       string
	X          float64
	Y          float64
	Points     string
	LineBounds []*engo.Line
}

// Bounds returns the level boundaries as an engo.AABB object
func (l *Level) Bounds() engo.AABB {
	return engo.AABB{
		Min: engo.Point{0, 0},
		Max: engo.Point{
			float32(l.TileWidth * l.width),
			float32(l.TileHeight * l.height),
		},
	}
}

// Width returns the integer width of the level
func (l *Level) Width() int {
	return l.width
}

// Height returns the integer height of the level
func (l *Level) Height() int {
	return l.height
}

// Height returns the integer height of the tile
func (t *tile) Height() float32 {
	return t.Image.Height()
}

// Width returns the integer width of the tile
func (t *tile) Width() float32 {
	return t.Image.Width()
}

// Texture returns the tile's Image texture
func (t *tile) Texture() *gl.Texture {
	return t.Image.id
}

// Close deletes the stored texture of a tile
func (t *tile) Close() {
	t.Image.Close()
}

// View returns the tile's viewport's min and max X & Y
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
	Width       int
	Height      int
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
			t.Image = &Texture{
				id:     sheet.Image.Texture,
				width:  tw,
				height: th,
				viewport: engo.AABB{
					engo.Point{u, v},
					engo.Point{u2, v2},
				},
			}
			tileset = append(tileset, t)
		}
	}

	return tileset
}

func createLevelTiles(lvl *Level, layers []*layer, ts []*tile) []*TileLayer {
	var levelTileLayers []*TileLayer

	// Create a TileLayer for each provided layer
	for _, layer := range layers {
		tilemap := make([]*tile, 0)
		tileLayer := &TileLayer{}
		mapping := layer.TileMapping

		for i := 0; i < lvl.height; i++ {
			for x := 0; x < lvl.width; x++ {
				idx := x + i*lvl.width
				t := &tile{}

				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ts[tileIdx].Image
					t.Point = engo.Point{
						float32(x * lvl.TileWidth),
						float32(i * lvl.TileHeight),
					}
				}
				tilemap = append(tilemap, t)
			}
		}

		tileLayer.Name = layer.Name
		tileLayer.Width = layer.Width
		tileLayer.Height = layer.Height
		tileLayer.Tiles = tilemap

		levelTileLayers = append(levelTileLayers, tileLayer)
	}

	return levelTileLayers
}
