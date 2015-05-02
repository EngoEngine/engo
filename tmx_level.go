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
	"strings"
)

// Just used to create levelTileset->Image
type TMXtilesetSrc struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type TMXtileset struct {
	Firstgid   int           `xml:"firstgid,attr"`
	Name       string        `xml:"name,attr"`
	TileWidth  int           `xml:"tilewidth,attr"`
	TileHeight int           `xml:"tileheight,attr"`
	ImageSrc   TMXtilesetSrc `xml:"image"`
	Image      *Texture
}

type TMXlayer struct {
	Name        string `xml:"name,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
	TileMapping []uint32
	// This variable doesn't need to persist, used to fill TileMapping
	CompData []byte `xml:"data"`
}

type TMXlevel struct {
	Width      int          `xml:"width,attr"`
	Height     int          `xml:"height,attr"`
	TileWidth  int          `xml:"tilewidth,attr"`
	TileHeight int          `xml:"tileheight,attr"`
	Tilesets   []TMXtileset `xml:"tileset"`
	Layers     []TMXlayer   `xml:"layer"`
}

type ByFirstgid []TMXtileset

func (t ByFirstgid) Len() int           { return len(t) }
func (t ByFirstgid) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByFirstgid) Less(i, j int) bool { return t[i].Firstgid < t[j].Firstgid }

// MUST BE base64 ENCODED and COMPRESSED WITH zlib!
func createLevelFromTmx(r Resource) (*Level, error) {
	tlvl := &TMXlevel{}
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

	ld := &LevelData{}
	ld.Width = tlvl.Width
	ld.Height = tlvl.Height
	ld.TileWidth = tlvl.TileWidth
	ld.TileHeight = tlvl.TileHeight

	// get the tilesheets in order and in generic format
	sort.Sort(ByFirstgid(tlvl.Tilesets))
	ts := make([]tilesheet, len(tlvl.Tilesets))
	for i, tts := range tlvl.Tilesets {
		ts[i] = tilesheet{tts.Image, tts.Firstgid}
	}
	ld.Tileset = createTileset(ts, ld.TileWidth, ld.TileHeight)

	ls := make([]layer, len(tlvl.Layers))
	for i, tls := range tlvl.Layers {
		ls[i] = layer{tls.Name, tls.TileMapping}
	}

	ld.Layers = ls
	lvl = createLevel(ld)
	return lvl, nil
}

func readTmx(url string) (string, error) {
	file, err := ioutil.ReadFile(url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
