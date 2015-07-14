package engi

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Just used to create levelTileset->Image
type TMXTilesetSrc struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type TMXTileset struct {
	Firstgid   int           `xml:"firstgid,attr"`
	Name       string        `xml:"name,attr"`
	TileWidth  int           `xml:"tilewidth,attr"`
	TileHeight int           `xml:"tileheight,attr"`
	ImageSrc   TMXTilesetSrc `xml:"image"`
	Image      *Texture
}

type TMXLayer struct {
	Name        string `xml:"name,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
	TileMapping []uint32
	// This variable doesn't need to persist, used to fill TileMapping
	CompData []byte `xml:"data"`
}

type TMXPolyline struct {
	Points string `xml:"points,attr"`
}

type TMXObj struct {
	X         float64       `xml:"x,attr"`
	Y         float64       `xml:"y,attr"`
	Polylines []TMXPolyline `xml:"polyline"`
}

type TMXObjGroup struct {
	Name    string   `xml:"name,attr"`
	Objects []TMXObj `xml:"object"`
}

type TMXImgSrc struct {
	Source string `xml:"source,attr"`
}

type TMXImgLayer struct {
	Name   string    `xml:"name,attr"`
	X      float64   `xml:"x,attr"`
	Y      float64   `xml:"y,attr"`
	ImgSrc TMXImgSrc `xml:"image"`
}

type TMXLevel struct {
	Width      int           `xml:"width,attr"`
	Height     int           `xml:"height,attr"`
	TileWidth  int           `xml:"tilewidth,attr"`
	TileHeight int           `xml:"tileheight,attr"`
	Tilesets   []TMXTileset  `xml:"tileset"`
	Layers     []TMXLayer    `xml:"layer"`
	ObjGroups  []TMXObjGroup `xml:"objectgroup"`
	ImgLayers  []TMXImgLayer `xml:"imagelayer"`
}

type ByFirstgid []TMXTileset

func (t ByFirstgid) Len() int           { return len(t) }
func (t ByFirstgid) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByFirstgid) Less(i, j int) bool { return t[i].Firstgid < t[j].Firstgid }

// MUST BE base64 ENCODED and COMPRESSED WITH zlib!
func createLevelFromTmx(r Resource) (*Level, error) {
	tlvl := &TMXLevel{}
	lvl := &Level{}

	tmx, err := readTmx(r.url)
	if err != nil {
		return lvl, err
	}

	if err := xml.Unmarshal([]byte(tmx), &tlvl); err != nil {
		fmt.Printf("Error unmarshalling XML: %v", err)
	}

	// Extract the tile mappings from the compressed data at each layer
	for idx := range tlvl.Layers {
		layer := &tlvl.Layers[idx]

		// Trim leading/trailing whitespace ( inneficient )
		layer.CompData = []byte(strings.TrimSpace(string(layer.CompData)))

		// Decode it out of base64
		if n, err := base64.StdEncoding.Decode(layer.CompData, layer.CompData); err != nil {
			fmt.Printf("error after %d bytes: %v", n, err)
			return lvl, err
		}

		// Decompress
		b := bytes.NewReader(layer.CompData)
		zlr, err := zlib.NewReader(b)
		if err != nil {
			fmt.Printf("error: %v", err)
			return lvl, err
		}

		tm := make([]uint32, 0)
		var nextInt uint32
		for {
			err = binary.Read(zlr, binary.LittleEndian, &nextInt)
			if err != nil {
				// this is -generally- EOF
				//fmt.Println("binary.Read failed:", err)
				break
			}
			tm = append(tm, nextInt)
		}
		layer.TileMapping = tm

		zlr.Close()
	}

	// Load in the images needed for the tilesets
	for k, ts := range tlvl.Tilesets {
		ts.Image = Files.Image(path.Base(ts.ImageSrc.Source))
		tlvl.Tilesets[k] = ts
	}

	lvl.Width = tlvl.Width
	lvl.Height = tlvl.Height
	lvl.TileWidth = tlvl.TileWidth
	lvl.TileHeight = tlvl.TileHeight

	// get the tilesheets in order and in generic format
	sort.Sort(ByFirstgid(tlvl.Tilesets))
	ts := make([]*tilesheet, len(tlvl.Tilesets))
	for i, tts := range tlvl.Tilesets {
		ts[i] = &tilesheet{tts.Image, tts.Firstgid}
	}

	lvlTileset := createTileset(lvl, ts)

	lvlLayers := make([]*layer, len(tlvl.Layers))
	for i, tls := range tlvl.Layers {
		lvlLayers[i] = &layer{tls.Name, tls.TileMapping}
	}

	lvl.Tiles = createLevelTiles(lvl, lvlLayers, lvlTileset)

	for _, o := range tlvl.ObjGroups[0].Objects {
		p := o.Polylines[0].Points
		lvl.LineBounds = append(lvl.LineBounds, pointStringToLines(p, o.X, o.Y)...)
	}

	for i := 0; i < len(tlvl.ImgLayers); i++ {
		curImg := Files.Image(path.Base(tlvl.ImgLayers[i].ImgSrc.Source))
		curX := float32(tlvl.ImgLayers[i].X)
		curY := float32(tlvl.ImgLayers[i].Y)
		reg := NewRegion(curImg, 0, 0, curImg.width, curImg.height)
		lvl.Images = append(lvl.Images, &tile{Point{curX, curY}, reg})
	}

	return lvl, nil
}

func pointStringToLines(str string, xOff, yOff float64) []Line {
	pts := strings.Split(str, " ")
	floatPts := make([][]float64, len(pts))
	for i, x := range pts {
		pt := strings.Split(x, ",")
		floatPts[i] = make([]float64, 2)
		floatPts[i][0], _ = strconv.ParseFloat(pt[0], 64)
		floatPts[i][1], _ = strconv.ParseFloat(pt[1], 64)
	}

	lines := make([]Line, len(floatPts)-1)

	// Now to globalize line coordinates
	for i := 0; i < len(floatPts)-1; i++ {
		x1 := float32(floatPts[i][0] + xOff)
		y1 := float32(floatPts[i][1] + yOff)
		x2 := float32(floatPts[i+1][0] + xOff)
		y2 := float32(floatPts[i+1][1] + yOff)
		p1 := Point{x1, y1}
		p2 := Point{x2, y2}
		newLine := Line{p1, p2}
		lines = append(lines, newLine)
	}

	return lines
}

func readTmx(url string) (string, error) {
	file, err := ioutil.ReadFile(url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
