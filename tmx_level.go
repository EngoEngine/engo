package engi

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

// Just used to create levelTileset->Image
type tilesetImgSrc struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type levelTileset struct {
	Firstgid   int           `xml:"firstgid,attr"`
	Name       string        `xml:"name,attr"`
	TileWidth  int           `xml:"tilewidth,attr"`
	TileHeight int           `xml:"tileheight,attr"`
	ImageSrc   tilesetImgSrc `xml:"image"`
	Image      *Texture
}

type levelLayer struct {
	Name        string `xml:"name,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
	TileMapping []uint32
	// This variable doesn't need to persist, used to fill TileMapping
	CompData []byte `xml:"data"`
}

type Level struct {
	Width      int            `xml:"width,attr"`
	Height     int            `xml:"height,attr"`
	TileWidth  int            `xml:"tilewidth,attr"`
	TileHeight int            `xml:"tileheight,attr"`
	Tilesets   []levelTileset `xml:"tileset"`
	Layers     []levelLayer   `xml:"layer"`
}

// MUST BE base64 ENCODED and COMPRESSED WITH zlib!
func createLevelFromTmx(r Resource) (*Level, error) {
	lvl := &Level{}
	tmx, err := readTmx(r.url)
	if err != nil {
		return lvl, err
	}

	if err := xml.Unmarshal([]byte(tmx), &lvl); err != nil {
		fmt.Printf("Error unmarshalling XML: %v", err)
	}

	// Extract the tile mappings from the compressed data at each layer
	for idx := range lvl.Layers {
		layer := &lvl.Layers[idx]

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
	//for i := 0; i < len(lvl.Tilesets); i++ {
	for k, ts := range lvl.Tilesets {
		//TODO
		url := "data/maps/" + ts.ImageSrc.Source
		name := "doesn't matter"
		r = NewResource(name, url)
		data, err := loadImage(r)
		if err != nil {
			return lvl, err
		}
		ts.Image = NewTexture(data)
		lvl.Tilesets[k] = ts
	}

	return lvl, nil
}

func readTmx(url string) (string, error) {
	file, err := ioutil.ReadFile(url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
