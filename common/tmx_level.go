package common

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"

	"engo.io/engo"
)

// TMXTilesetSrc is just used to create levelTileset->Image
type TMXTilesetSrc struct {
	// Source holds the URI of the tileset image
	Source string `xml:"source,attr"`
	// Width of each tile in the tileset image
	Width int `xml:"width,attr"`
	// Height of each tile in the tileset image
	Height int `xml:"height,attr"`
}

// TMXTileset contains the tileset resource parsed from the TileMap XML
type TMXTileset struct {
	// Firstgid is the first assigned gid of the tileset
	Firstgid int `xml:"firstgid,attr"`
	// Name of the tileset in Tiled
	Name string `xml:"name,attr"`
	// TileWidth defines the width of each tile
	TileWidth int `xml:"tilewidth,attr"`
	// TileHeight defines the height of each tile
	TileHeight int `xml:"tileheight,attr"`
	// TileCount holds the total tile count in this tileset
	TileCount int `xml:"tilecount,attr"`
	// ImageSrc contains the TMXTilesetSrc which defines the tileset image
	ImageSrc TMXTilesetSrc `xml:"image"`
	// Image holds the reference of the tileset's TextureResource
	Image *TextureResource
	// Tiles are any tiles inside the tileset
	Tiles []TMXTile `xml:"tile"`
}

// TMXTileLayer represents a tile layer parsed from the TileMap XML
type TMXTileLayer struct {
	// Name of the tile layer in Tiled
	Name string `xml:"name,attr"`
	// Width is the integer width of each tile in this layer
	Width int `xml:"width,attr"`
	// Height is the integer height of each tile in this layer
	Height int `xml:"height,attr"`
	// TileMapping contains the generated tilemapping list
	TileMapping []uint32
	// TMXData is the encoded tile layer grid
	Data TMXData `xml:"data"`
}

// TMXData represents custom properties which "Can be used as a child
// of the map, tileset, tile (when part of a tileset), layer, objectgroup,
// object and imagelayer elements"
type TMXData struct {
	// Encoding defines the fomat of the Data field; Valid values are
	// one of "", "base64", or "csv"
	Encoding string `xml:"encoding,attr"`
	// Compression defines the compression applied to base64 data
	// valid values are one of "", "zlib", "gzip"
	Compression string `xml:"compression,attr"`
	// Data contains an encoded list of uint32 guids defining
	// the tile layout of a layer
	Data string `xml:",innerxml,"`
	// Tiles is an array of tiles containing guids. Not set if
	// other encodings are used
	Tiles []TMXTile `xml:"tile"`
}

// TMXTile represents a single tile on a tile layer.
type TMXTile struct {
	// gid represents a single tile encoded with its flip
	// orientation
	Gid uint32 `xml:"gid,attr"`
	// ImageSrc is the image on the tile
	ImageSrc TMXImageSrc `xml:"image"`
	// Image holds the reference of the tileset's TextureResource
	Image *TextureResource
}

// ErrUnknownEncoding is when the encoding of the TMX is unsupported.
var ErrUnknownEncoding = errors.New("Unknown Encoding")

// ErrUnknownCompression is when the compression of the TMX is unsupported.
var ErrUnknownCompression = errors.New("Unknown Compression")

// tileDecode() creates a decoded array of gids from xml tile tags
func (d *TMXData) decodeTile() ([]uint32, error) {
	tm := make([]uint32, 0)
	for _, t := range d.Tiles {
		tm = append(tm, t.Gid)
	}
	return tm, nil
}

// decodeCSV() creates a decoded array of gids from csv
func (d *TMXData) decodeCSV() ([]uint32, error) {
	b := strings.NewReader(strings.TrimSpace(d.Data))
	cr := csv.NewReader(b)
	// We allow variable number of fields per record to allow line ending commas and then
	// empty strings appearing as a field. Later, we filter empty strings. This trick is
	// needed to allow TilEd-style CSVs with line-ending commas but no comma at the end
	// of last line.
	cr.FieldsPerRecord = -1
	tm := make([]uint32, 0)
	if recs, err := cr.ReadAll(); err == nil {
		if len(recs) < 1 {
			return nil, errors.New("No csv records found")
		}
		for _, rec := range recs {
			for i, id := range rec {
				// An empty string appearing after last comma. We filter it.
				if id == "" && i == len(rec)-1 {
					continue
				}
				if nextInt, err2 := strconv.ParseUint(id, 10, 32); err == nil {
					tm = append(tm, uint32(nextInt))
				} else {
					return nil, err2
				}
			}
		}
		if len(tm) < 1 {
			return nil, errors.New("No Data Returned")
		}
	} else {
		return nil, err
	}
	return tm, nil
}

// Decode takes the encoded data from a tmx map file and
// unpacks it an array of uint32 guids
func (d *TMXData) Decode() ([]uint32, error) {
	// Tile tag and CSV encodings
	if len(d.Tiles) > 0 {
		return d.decodeTile()
	} else if d.Encoding == "csv" {
		return d.decodeCSV()
	}

	// Only encoding in the standard is base64
	var breader io.Reader
	if d.Encoding == "base64" {
		buff, err := base64.StdEncoding.DecodeString(strings.TrimSpace(d.Data))
		if err != nil {
			return nil, err
		}
		breader = bytes.NewReader(buff)
	} else {
		return nil, ErrUnknownEncoding
	}

	// Setup decompression if needed
	var zreader io.Reader
	if d.Compression == "" {
		zreader = breader
	} else if d.Compression == "zlib" {
		z, err := zlib.NewReader(breader)
		if err != nil {
			return nil, err
		}
		defer z.Close()
		zreader = z
	} else if d.Compression == "gzip" {
		z, err := gzip.NewReader(breader)
		if err != nil {
			return nil, err
		}
		defer z.Close()
		zreader = z
	} else {
		return nil, ErrUnknownCompression
	}

	// decode into tm and return it.
	tm := make([]uint32, 0)
	var nextInt uint32
	for {
		err := binary.Read(zreader, binary.LittleEndian, &nextInt)
		if err != nil {
			// EOF or unexpected EOF error
			if err == io.EOF {
				break
			}

			return nil, err
		}
		tm = append(tm, nextInt)
	}
	return tm, nil
}

// TMXImageLayer represents an image layer parsed from the TileMap XML
type TMXImageLayer struct {
	// Name of the image layer in Tiled
	Name string `xml:"name,attr"`
	// X holds the defined X coordinate in Tiled
	X float64 `xml:"x,attr"`
	// Y holds the defined Y coordinate in Tiled
	Y float64 `xml:"y,attr"`
	// ImageSrc contains the TMXImageSrc which defines the image filename
	ImageSrc TMXImageSrc `xml:"image"`
}

// TMXObjectLayer following the Object Layer naming convention in Tiled
type TMXObjectLayer struct {
	// Name of the object layer in Tiled
	Name string `xml:"name,attr"`
	// Objects contains the list all objects in this layer
	Objects []TMXObject `xml:"object"`
	// OffSetX is the parsed X offset for the object layer
	OffSetX float32 `xml:"offsetx"`
	// OffSetY is the parsed Y offset for the object layer
	OffSetY float32 `xml:"offsety"`
	//TODO add all object attr available in Tiled
	// 'visible' attr only appears in XML if false --> default 1
	// 'opacity' attr only appears in XML if < 1 --> default 1
	// Visible int         `xml:"visible,attr"`
	// Opacity float32     `xml:"visible,attr"`
}

// TMXObject represents a TMX object with all default Tiled attributes
type TMXObject struct {
	// Id is the unique ID of each object defined by Tiled
	Id int `xml:"id,attr"`
	// Name defines the name of the object given in Tiled
	Name string `xml:"name,attr"`
	// Type contains the string type which was given in Tiled
	Type string `xml:"type,attr"`
	// X holds the X float64 coordinate of the object in the map
	X float64 `xml:"x,attr"`
	// Y holds the Y float64 coordinate of the object in the map
	Y float64 `xml:"y,attr"`
	// Width is the integer width of the object
	Width int `xml:"width,attr"`
	// Height is the integer height of the object
	Height int `xml:"height,attr"`
	// Gid is the integer global tile ID
	Gid int `xml:"gid,attr"`
	// Polyline contains the TMXPolyline object if the parsed object has a polyline points string
	Polyline TMXPolyline `xml:"polyline"`
}

// TMXPolyline represents a TMX Polyline object with its Points values
type TMXPolyline struct {
	// Points contains the original, unaltered points string from the TMZ XML
	Points string `xml:"points,attr"`
}

// TMXImageSrc represents the actual image source of an image layer
type TMXImageSrc struct {
	// Source holds the URI of the image
	Source string `xml:"source,attr"`
	// Width is the integer width of the image
	Width int `xml:"width,attr"`
	// Height is the integer height of the image
	Height int `xml:"height,attr"`
}

// TMXLevel containing all layers and default Tiled attributes
type TMXLevel struct {
	// Orientation is the parsed level orientation from the TMX XML, like orthogonal
	Orientation string `xml:"orientation,attr"`
	// RenderOrder is the in Tiled specified TileMap render order, like right-down
	RenderOrder string `xml:"renderorder,attr"`
	// Width is the integer width of the parsed level
	Width int `xml:"width,attr"`
	// Height is the integer height of the parsed level
	Height int `xml:"height,attr"`
	// TileWidth defines the width of each tile in the level
	TileWidth int `xml:"tilewidth,attr"`
	// TileHeight defines the height of each tile in the level
	TileHeight int `xml:"tileheight,attr"`
	// NextObjectId is the next free Object ID defined by Tiled
	NextObjectId int `xml:"nextobjectid,attr"`
	// Tilesets contains a list of all parsed TMXTileset objects
	Tilesets []TMXTileset `xml:"tileset"`
	// TileLayers contains a list of all parsed TMXTileLayer objects
	TileLayers []TMXTileLayer `xml:"layer"`
	// ImageLayers contains a list of all parsed TMXImageLayer objects
	ImageLayers []TMXImageLayer `xml:"imagelayer"`
	// ObjectLayers contains a list of all parsed TMXObjectLayer objects
	ObjectLayers []TMXObjectLayer `xml:"objectgroup"`
}

// ByFirstgid is a []TMXTileset sorted by First gid
type ByFirstgid []TMXTileset

// Len returns the length of t
func (t ByFirstgid) Len() int { return len(t) }

// Swap exchanges t's elements i and j
func (t ByFirstgid) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less returns if t's i Firstgid is less than t's j
func (t ByFirstgid) Less(i, j int) bool { return t[i].Firstgid < t[j].Firstgid }

func getTextureResourcesFromTmx(url string) (*TextureResource, error) {
	if err := engo.Files.Load(url); err != nil {
		return nil, err
	}
	image, err := engo.Files.Resource(url)
	if err != nil {
		return nil, err
	}
	texResource, ok := image.(TextureResource)
	if !ok {
		return nil, fmt.Errorf("resource is not of type 'TextureResource': %q", url)
	}
	return &texResource, nil
}

// createLevelFromTmx unmarshales and unpacks tmx data into a Level
func createLevelFromTmx(tmxBytes []byte, tmxURL string) (*Level, error) {
	tmxLevel := &TMXLevel{}
	level := &Level{}

	if err := xml.Unmarshal(tmxBytes, &tmxLevel); err != nil {
		return nil, err
	}

	// Extract the tile mappings from the compressed data at each layer
	for i := range tmxLevel.TileLayers {
		layer := &tmxLevel.TileLayers[i]
		tm, err := layer.Data.Decode()
		if err != nil {
			return nil, err
		}
		layer.TileMapping = tm
	}

	// Load in the images needed for the tilesets
	for i, ts := range tmxLevel.Tilesets {
		for j, tile := range ts.Tiles {
			texResource, err := getTextureResourcesFromTmx(path.Join(path.Dir(tmxURL), tile.ImageSrc.Source))
			if err != nil {
				return nil, err
			}
			ts.Tiles[j].Image = texResource
		}
		if ts.ImageSrc.Source == "" {
			tmxLevel.Tilesets[i] = ts
			continue
		}
		texResource, err := getTextureResourcesFromTmx(path.Join(path.Dir(tmxURL), ts.ImageSrc.Source))
		if err != nil {
			return nil, err
		}
		ts.Image = texResource
		tmxLevel.Tilesets[i] = ts
	}

	level.width = tmxLevel.Width
	level.height = tmxLevel.Height
	level.TileWidth = tmxLevel.TileWidth
	level.TileHeight = tmxLevel.TileHeight
	level.Orientation = tmxLevel.Orientation
	level.RenderOrder = tmxLevel.RenderOrder
	level.NextObjectId = tmxLevel.NextObjectId

	// get the tilesheets in order and in generic format
	sort.Sort(ByFirstgid(tmxLevel.Tilesets))
	ts := make([]*tilesheet, len(tmxLevel.Tilesets))
	for i, tts := range tmxLevel.Tilesets {
		tiles := make([]Tile, 0)
		for _, t := range tts.Tiles {
			tile := Tile{}
			tile.Image = &Texture{t.Image.Texture, t.Image.Width, t.Image.Height, engo.AABB{Max: engo.Point{X: 1.0, Y: 1.0}}}
			tile.Point = engo.Point{X: 0, Y: 0}
			tiles = append(tiles, tile)
		}
		ts[i] = &tilesheet{tts.Image, tts.Firstgid, tts.TileWidth, tts.TileHeight, tiles}
	}

	level.Tileset = createTileset(level, ts)

	levelTileLayers := make([]*layer, len(tmxLevel.TileLayers))
	for i, tileLayer := range tmxLevel.TileLayers {
		levelTileLayers[i] = &layer{
			tileLayer.Name,
			tileLayer.Width,
			tileLayer.Height,
			tileLayer.TileMapping,
		}
	}

	// create Tile layers with tiles
	level.TileLayers, level.tilemap = createLevelTiles(level, levelTileLayers, level.Tileset)

	// create object layers
	for _, objectLayer := range tmxLevel.ObjectLayers {

		levelObjectLayer := &ObjectLayer{
			Name:    objectLayer.Name,
			OffSetX: objectLayer.OffSetX,
			OffSetY: objectLayer.OffSetY,
		}

		// check all objects in layer
		for _, tmxObject := range objectLayer.Objects {

			// check if object is a Polyline object
			if tmxObject.Polyline.Points != "" {

				points := tmxObject.Polyline.Points

				polylineObject := &PolylineObject{
					Id:     tmxObject.Id,
					Name:   tmxObject.Name,
					Type:   tmxObject.Type,
					X:      tmxObject.X,
					Y:      tmxObject.Y,
					Points: points,
				}

				polylineObject.LineBounds =
					append(
						polylineObject.LineBounds,
						pointStringToLines(points, tmxObject.X, tmxObject.Y)...,
					)

				levelObjectLayer.PolyObjects =
					append(levelObjectLayer.PolyObjects, polylineObject)

			} else {
				// non-Polyline object
				object := &Object{
					tmxObject.Id,
					tmxObject.Name,
					tmxObject.Type,
					tmxObject.X,
					tmxObject.Y,
					tmxObject.Width,
					tmxObject.Height,
				}

				levelObjectLayer.Objects = append(levelObjectLayer.Objects, object)

			}
		}

		level.ObjectLayers = append(level.ObjectLayers, levelObjectLayer)
	}

	// One image per image layer
	for _, tmxImageLayer := range tmxLevel.ImageLayers {

		imageLayer := &ImageLayer{
			Name:   tmxImageLayer.Name,
			Width:  tmxImageLayer.ImageSrc.Width,
			Height: tmxImageLayer.ImageSrc.Height,
			Source: tmxImageLayer.ImageSrc.Source,
		}

		url := path.Base(tmxImageLayer.ImageSrc.Source)
		if err := engo.Files.Load(url); err != nil {
			return nil, err
		}

		curImg, err := LoadedSprite(url)
		if err != nil {
			return nil, err
		}

		// create image tile
		imageTile := &Tile{
			engo.Point{X: float32(tmxImageLayer.X), Y: float32(tmxImageLayer.Y)},
			curImg,
		}

		imageLayer.Images = append(imageLayer.Images, imageTile)
		level.ImageLayers = append(level.ImageLayers, imageLayer)
	}

	return level, nil
}

func pointStringToLines(str string, xOff, yOff float64) []*engo.Line {
	pts := strings.Split(str, " ")
	floatPts := make([][]float64, len(pts))
	for i, x := range pts {
		pt := strings.Split(x, ",")
		floatPts[i] = make([]float64, 2)
		floatPts[i][0], _ = strconv.ParseFloat(pt[0], 64)
		floatPts[i][1], _ = strconv.ParseFloat(pt[1], 64)
	}

	lines := make([]*engo.Line, len(floatPts)-1)

	// Now to globalize line coordinates
	for i := 0; i < len(floatPts)-1; i++ {
		x1 := float32(floatPts[i][0] + xOff)
		y1 := float32(floatPts[i][1] + yOff)
		x2 := float32(floatPts[i+1][0] + xOff)
		y2 := float32(floatPts[i+1][1] + yOff)

		p1 := engo.Point{X: x1, Y: y1}
		p2 := engo.Point{X: x2, Y: y2}
		newLine := &engo.Line{P1: p1, P2: p2}

		lines[i] = newLine
	}

	return lines
}
