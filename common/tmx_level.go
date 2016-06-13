package common

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
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
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// TMXTileset contains the tileset resource parsed from the TileMap XML
type TMXTileset struct {
	Firstgid   int           `xml:"firstgid,attr"`
	Name       string        `xml:"name,attr"`
	TileWidth  int           `xml:"tilewidth,attr"`
	TileHeight int           `xml:"tileheight,attr"`
	TileCount  int           `xml:"tilecount,attr"`
	ImageSrc   TMXTilesetSrc `xml:"image"`
	Image      *TextureResource
}

// TMXTileLayer represents a tile layer parsed from the TileMap XML
type TMXTileLayer struct {
	Name        string `xml:"name,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
	TileMapping []uint32
	// This variable doesn't need to persist, used to fill TileMapping
	CompData []byte `xml:"data"`
}

// TMXImageLayer represents an image layer parsed from the TileMap XML
type TMXImageLayer struct {
	Name     string      `xml:"name,attr"`
	X        float64     `xml:"x,attr"`
	Y        float64     `xml:"y,attr"`
	ImageSrc TMXImageSrc `xml:"image"`
}

// TMXObjectLayer following the Object Layer naming convention in Tiled
type TMXObjectLayer struct {
	//TODO add all object attr available in Tiled
	// 'visible' attr only appears in XML if false --> default 1
	// 'opacity' attr only appears in XML if < 1 --> default 1
	Name    string      `xml:"name,attr"`
	Objects []TMXObject `xml:"object"`
	OffSetX float32     `xml:"offsetx"`
	OffSetY float32     `xml:"offsety"`
	// Visible int         `xml:"visible,attr"`
	// Opacity float32     `xml:"visible,attr"`
}

// TMXObject represents a TMX object with all default Tiled attributes
type TMXObject struct {
	Id       int         `xml:"id,attr"`
	Name     string      `xml:"name,attr"`
	Type     string      `xml:"type,attr"`
	X        float64     `xml:"x,attr"`
	Y        float64     `xml:"y,attr"`
	Width    int         `xml:"width,attr"`
	Height   int         `xml:"height,attr"`
	Polyline TMXPolyline `xml:"polyline"`
}

// TMXPolyline represents a TMX Polyline object with its Points values
type TMXPolyline struct {
	Points string `xml:"points,attr"`
}

// TMXImageSrc represents the actual image source of an image layer
type TMXImageSrc struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// TMXLevel containing all layers and default Tiled attributes
type TMXLevel struct {
	Orientation  string           `xml:"orientation,attr"`
	RenderOrder  string           `xml:"renderorder,attr"`
	Width        int              `xml:"width,attr"`
	Height       int              `xml:"height,attr"`
	TileWidth    int              `xml:"tilewidth,attr"`
	TileHeight   int              `xml:"tileheight,attr"`
	NextObjectId int              `xml:"nextobjectid,attr"`
	Tilesets     []TMXTileset     `xml:"tileset"`
	TileLayers   []TMXTileLayer   `xml:"layer"`
	ImageLayers  []TMXImageLayer  `xml:"imagelayer"`
	ObjectLayers []TMXObjectLayer `xml:"objectgroup"`
}

type ByFirstgid []TMXTileset

// Len returns the length of t
func (t ByFirstgid) Len() int { return len(t) }

// Swap exchanges t's elements i and j
func (t ByFirstgid) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less returns if t's i Firstgid is less than t's j
func (t ByFirstgid) Less(i, j int) bool { return t[i].Firstgid < t[j].Firstgid }

// MUST BE base64 ENCODED and COMPRESSED WITH zlib!
func createLevelFromTmx(tmxBytes []byte, tmxUrl string) (*Level, error) {
	tmxLevel := &TMXLevel{}
	level := &Level{}

	if err := xml.Unmarshal(tmxBytes, &tmxLevel); err != nil {
		return nil, err
	}

	// Extract the tile mappings from the compressed data at each layer
	for i := range tmxLevel.TileLayers {
		layer := &tmxLevel.TileLayers[i]

		// Trim leading/trailing whitespace ( inneficient )
		layer.CompData = []byte(strings.TrimSpace(string(layer.CompData)))

		// Decode it out of base64
		if _, err := base64.StdEncoding.Decode(layer.CompData, layer.CompData); err != nil {
			return nil, err
		}

		// Decompress
		b := bytes.NewReader(layer.CompData)
		zlr, err := zlib.NewReader(b)
		if err != nil {
			return nil, err
		}
		defer zlr.Close()

		tm := make([]uint32, 0)
		var nextInt uint32
		for {
			err = binary.Read(zlr, binary.LittleEndian, &nextInt)
			if err != nil {
				// EOF or unexpected EOF error
				if err == io.EOF {
					break
				}

				return nil, err
			}
			tm = append(tm, nextInt)
		}
		layer.TileMapping = tm
	}

	// Load in the images needed for the tilesets
	for i, ts := range tmxLevel.Tilesets {
		url := path.Join(path.Dir(tmxUrl), ts.ImageSrc.Source)
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
		ts.Image = &texResource
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
		ts[i] = &tilesheet{tts.Image, tts.Firstgid}
	}

	levelTileset := createTileset(level, ts)

	levelTileLayers := make([]*layer, len(tmxLevel.TileLayers))
	for i, tileLayer := range tmxLevel.TileLayers {
		levelTileLayers[i] = &layer{
			tileLayer.Name,
			tileLayer.Width,
			tileLayer.Height,
			tileLayer.TileMapping,
		}
	}

	// create tile layers with tiles
	level.TileLayers = createLevelTiles(level, levelTileLayers, levelTileset)

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

		curImg, err := PreloadedSpriteSingle(url)
		if err != nil {
			return nil, err
		}

		// create image tile
		imageTile := &tile{
			engo.Point{float32(tmxImageLayer.X), float32(tmxImageLayer.Y)},
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

		p1 := engo.Point{x1, y1}
		p2 := engo.Point{x2, y2}
		newLine := &engo.Line{p1, p2}

		lines[i] = newLine
	}

	return lines
}
