package common

import (
	"engo.io/engo"
	"engo.io/gl"
)

// Level is a parsed TMX level containing all layers and default Tiled attributes
type Level struct {
	// Orientation is the parsed level orientation from the TMX XML, like orthogonal, isometric, etc.
	Orientation string
	// RenderOrder is the in Tiled specified TileMap render order, like right-down, right-up, etc.
	RenderOrder string
	width       int
	height      int
	// TileWidth defines the width of each tile in the level
	TileWidth int
	// TileHeight defines the height of each tile in the level
	TileHeight int
	// NextObjectId is the next free Object ID defined by Tiled
	NextObjectId int
	// TileLayers contains all TileLayer of the level
	TileLayers []*TileLayer
	// ImageLayers contains all ImageLayer of the level
	ImageLayers []*ImageLayer
	// ObjectLayers contains all ObjectLayer of the level
	ObjectLayers []*ObjectLayer
}

// TileLayer contains a list of its tiles plus all default Tiled attributes
type TileLayer struct {
	// Name defines the name of the tile layer given in the TMX XML / Tiled
	Name string
	// Width is the integer width of each tile in this layer
	Width int
	// Height is the integer height of each tile in this layer
	Height int
	// Tiles contains the list of tiles
	Tiles []*tile
}

// ImageLayer contains a list of its images plus all default Tiled attributes
type ImageLayer struct {
	// Name defines the name of the image layer given in the TMX XML / Tiled
	Name string
	// Width is the integer width of each image in this layer
	Width int
	// Height is the integer height of each image in this layer
	Height int
	// Source contains the original image filename
	Source string
	// Images contains the list of all image tiles
	Images []*tile
}

// ObjectLayer contains a list of its standard objects as well as a list of all its polyline objects
type ObjectLayer struct {
	// Name defines the name of the object layer given in the TMX XML / Tiled
	Name string
	// OffSetX is the parsed X offset for the object layer
	OffSetX float32
	// OffSetY is the parsed Y offset for the object layer
	OffSetY float32
	// Objects contains the list of (regular) Object objects
	Objects []*Object
	// PolyObjects contains the list of PolylineObject objects
	PolyObjects []*PolylineObject
}

// Object is a standard TMX object with all its default Tiled attributes
type Object struct {
	// Id is the unique ID of each object defined by Tiled
	Id int
	// Name defines the name of the object given in Tiled
	Name string
	// Type contains the string type which was given in Tiled
	Type string
	// X holds the X float64 coordinate of the object in the map
	X float64
	// X holds the X float64 coordinate of the object in the map
	Y float64
	// Width is the integer width of the object
	Width int
	// Height is the integer height of the object
	Height int
}

// PolylineObject is a TMX polyline object with all its default Tiled attributes
type PolylineObject struct {
	// Id is the unique ID of each polyline object defined by Tiled
	Id int
	// Name defines the name of the polyline object given in Tiled
	Name string
	// Type contains the string type which was given in Tiled
	Type string
	// X holds the X float64 coordinate of the polyline in the map
	X float64
	// Y holds the Y float64 coordinate of the polyline in the map
	Y float64
	// Points contains the original, unaltered points string from the TMZ XML
	Points string
	// LineBounds is the list of engo.Line objects generated from the points string
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
