package common

import (
	"engo.io/engo"
	"engo.io/engo/math"
	"engo.io/gl"
)

const (
	orth = "orthogonal"
	iso  = "isometric"
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
	// Tileset maps tile IDs to their texture
	Tileset map[int]*Tile
	// TileLayers contains all TileLayer of the level
	TileLayers []*TileLayer
	// ImageLayers contains all ImageLayer of the level
	ImageLayers []*ImageLayer
	// ObjectLayers contains all ObjectLayer of the level
	ObjectLayers []*ObjectLayer
	// tilemap maps tile map position to the tile at that location
	tilemap map[mapPoint]*Tile
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
	Tiles []*Tile
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
	Images []*Tile
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

type mapPoint struct {
	X, Y int
}

// Bounds returns the level boundaries as an engo.AABB object
func (l *Level) Bounds() engo.AABB {
	switch l.Orientation {
	case orth:
		return engo.AABB{
			Min: l.screenPoint(engo.Point{X: 0, Y: 0}),
			Max: l.screenPoint(engo.Point{X: float32(l.width), Y: float32(l.height)}),
		}
	case iso:
		xMin := l.screenPoint(engo.Point{X: 0, Y: float32(l.height)}).X + float32(l.TileWidth)/2
		xMax := l.screenPoint(engo.Point{X: float32(l.width), Y: 0}).X + float32(l.TileWidth)/2
		yMin := l.screenPoint(engo.Point{X: 0, Y: 0}).Y
		yMax := l.screenPoint(engo.Point{X: float32(l.width), Y: float32(l.height)}).Y + float32(l.TileHeight)/2
		return engo.AABB{
			Min: engo.Point{X: xMin, Y: yMin},
			Max: engo.Point{X: xMax, Y: yMax},
		}
	}
	return engo.AABB{}
}

// mapPoint returns the map point of the passed in screen point
func (l *Level) mapPoint(screenPt engo.Point) engo.Point {
	switch l.Orientation {
	case orth:
		screenPt.Multiply(engo.Point{X: 1 / float32(l.TileWidth), Y: 1 / float32(l.TileHeight)})
		return screenPt
	case iso:
		return engo.Point{
			X: (screenPt.X / float32(l.TileWidth)) + (screenPt.Y / float32(l.TileHeight)),
			Y: (screenPt.Y / float32(l.TileHeight)) - (screenPt.X / float32(l.TileWidth)),
		}
	}
	return engo.Point{X: 0, Y: 0}
}

// screenPoint returns the screen point of the passed in map point
func (l *Level) screenPoint(mapPt engo.Point) engo.Point {
	switch l.Orientation {
	case orth:
		mapPt.Multiply(engo.Point{X: float32(l.TileWidth), Y: float32(l.TileHeight)})
		return mapPt
	case iso:
		return engo.Point{
			X: (mapPt.X - mapPt.Y) * float32(l.TileWidth) / 2,
			Y: (mapPt.X + mapPt.Y) * float32(l.TileHeight) / 2,
		}
	}
	return engo.Point{X: 0, Y: 0}
}

// GetTile returns a *Tile at the given point (in space / render coordinates).
func (l *Level) GetTile(pt engo.Point) *Tile {
	mp := l.mapPoint(pt)
	x := int(math.Floor(mp.X))
	y := int(math.Floor(mp.Y))
	t, ok := l.tilemap[mapPoint{X: x, Y: y}]
	if !ok {
		return nil
	}
	return t
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
func (t *Tile) Height() float32 {
	return t.Image.Height()
}

// Width returns the integer width of the tile
func (t *Tile) Width() float32 {
	return t.Image.Width()
}

// Texture returns the tile's Image texture
func (t *Tile) Texture() *gl.Texture {
	return t.Image.id
}

// Close deletes the stored texture of a tile
func (t *Tile) Close() {
	t.Image.Close()
}

// View returns the tile's viewport's min and max X & Y
func (t *Tile) View() (float32, float32, float32, float32) {
	return t.Image.View()
}

// Tile represents a tile in the TMX map.
type Tile struct {
	engo.Point
	Image *Texture
}

type tilesheet struct {
	Image    *TextureResource
	Firstgid int
	Width    int
	Height   int
	Tiles    []Tile
}

type layer struct {
	Name        string
	Width       int
	Height      int
	TileMapping []uint32
}

func createTileset(lvl *Level, sheets []*tilesheet) map[int]*Tile {
	tileset := make(map[int]*Tile)
	deftw := float32(lvl.TileWidth)
	defth := float32(lvl.TileHeight)

	for _, sheet := range sheets {
		var tw, th = deftw, defth
		curGid := sheet.Firstgid
		if sheet.Height != 0 && sheet.Width != 0 {
			tw, th = float32(sheet.Width), float32(sheet.Height)
		}
		for i := range sheet.Tiles {
			tileset[curGid] = &sheet.Tiles[i]
			curGid++
		}
		if sheet.Image == nil {
			continue
		}
		setWidth := sheet.Image.Width / tw
		setHeight := sheet.Image.Height / th
		totalTiles := int(setWidth * setHeight)

		for i := 0; i < totalTiles; i++ {
			t := &Tile{}
			x := float32(i%int(setWidth)) * tw
			y := float32(i/int(setWidth)) * th

			invTexWidth := 1.0 / sheet.Image.Width
			invTexHeight := 1.0 / sheet.Image.Height

			u := x * invTexWidth
			v := y * invTexHeight
			u2 := (x + tw) * invTexWidth
			v2 := (y + th) * invTexHeight
			t.Image = &Texture{
				id:     sheet.Image.Texture,
				width:  tw,
				height: th,
				viewport: engo.AABB{
					Min: engo.Point{X: u, Y: v},
					Max: engo.Point{X: u2, Y: v2},
				},
			}
			tileset[curGid] = t
			curGid++
		}
	}

	return tileset
}

func createLevelTiles(lvl *Level, layers []*layer, ts map[int]*Tile) ([]*TileLayer, map[mapPoint]*Tile) {
	var levelTileLayers []*TileLayer
	tileMapper := make(map[mapPoint]*Tile)

	// Create a TileLayer for each provided layer
	for _, layer := range layers {
		tilemap := make([]*Tile, 0)
		tileLayer := &TileLayer{}
		mapping := layer.TileMapping

		for i := 0; i < lvl.height; i++ {
			for x := 0; x < lvl.width; x++ {
				idx := x + i*lvl.width
				idxPt := mapPoint{X: x, Y: i}
				t := &Tile{}

				if tileIdx := int(mapping[idx]); tileIdx >= 1 {
					t.Image = ts[tileIdx].Image
					t.Point = lvl.screenPoint(engo.Point{X: float32(x), Y: float32(i)})
				}
				tilemap = append(tilemap, t)
				tileMapper[idxPt] = t
			}
		}

		tileLayer.Name = layer.Name
		tileLayer.Width = layer.Width
		tileLayer.Height = layer.Height
		tileLayer.Tiles = tilemap

		levelTileLayers = append(levelTileLayers, tileLayer)
	}

	return levelTileLayers, tileMapper
}
