// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"encoding/json"
	"github.com/errcw/glow/gl/2.1/gl"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// A Texture wraps an opengl texture and is mostly used for loading
// images and constructing Regions.
type Texture struct {
	id        uint32
	width     int32
	height    int32
	minFilter int32
	maxFilter int32
	uWrap     int
	vWrap     int
}

// NewTexture takes either a string path to an image file, an
// io.Reader containing image date or an image.Image and returns a Texture.
func NewTexture(data interface{}) *Texture {
	var m image.Image

	switch data := data.(type) {
	default:
		log.Fatal("NewTexture needs a string or io.Reader")
	case string:
		file, err := os.Open(data)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case io.Reader:
		img, _, err := image.Decode(data)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case image.Image:
		m = data
	}

	b := m.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	width := int32(m.Bounds().Max.X)
	height := int32(m.Bounds().Max.Y)

	var id uint32
	gl.GenTextures(1, &id)

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, WrapClampToEdge)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, WrapClampToEdge)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, FilterLinearMipMapNearest)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, FilterNearest)
	gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(newm.Pix))

	gl.Disable(gl.TEXTURE_2D)

	return &Texture{id, width, height, gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE}
}

// Split creates Regions from every width, height rect going from left
// to right, then down. This is useful for simple images with uniform cells.
func (t *Texture) Split(w, h int) []*Region {
	x := 0
	y := 0
	width := int(t.Width())
	height := int(t.Height())

	rows := height / h
	cols := width / w

	startX := x
	tiles := make([]*Region, 0)
	for row := 0; row < rows; row++ {
		x = startX
		for col := 0; col < cols; col++ {
			tiles = append(tiles, NewRegion(t, x, y, w, h))
			x += w
		}
		y += h
	}

	return tiles
}

func (t *Texture) Unpack(path string) map[string]*Region {
	regions := make(map[string]*Region)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var data interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	root := data.(map[string]interface{})
	frames := root["frames"].([]interface{})
	for _, frameData := range frames {
		frame := frameData.(map[string]interface{})
		name := strings.Split(frame["filename"].(string), ".")[0]
		rect := frame["frame"].(map[string]interface{})
		x := int(rect["x"].(float64))
		y := int(rect["y"].(float64))
		w := int(rect["w"].(float64))
		h := int(rect["h"].(float64))
		regions[name] = NewRegion(t, x, y, w, h)
	}

	return regions
}

// Delete will dispose of the texture.
func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.id)
}

// Bind will bind the texture.
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind will unbind all textures.
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Width returns the width of the texture.
func (t *Texture) Width() int32 {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() int32 {
	return t.height
}

// SetFilter sets the filter type used when scaling a texture up or
// down. The default is nearest which will not doing any interpolation
// between pixels.
func (t *Texture) SetFilter(min, max int32) {
	t.minFilter = min
	t.maxFilter = max
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, max)
}

// Returns the current min and max filters used.
func (t *Texture) Filter() (int32, int32) {
	return t.minFilter, t.maxFilter
}

func (t *Texture) SetWrap(u, v int) {
	t.uWrap = u
	t.vWrap = v
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(u))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(v))
}

func (t *Texture) Wrap() (int, int) {
	return t.uWrap, t.vWrap
}
