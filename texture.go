// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"encoding/json"
	"log"
	"strings"
)

// A Texture wraps an opengl texture and is mostly used for loading
// images and constructing Regions.
type Texture struct {
	id        *TextureObject
	width     int
	height    int
	minFilter int
	maxFilter int
	uWrap     int
	vWrap     int
}

// NewTexture takes either a string path to an image file, an
// io.Reader containing image date or an image.Image and returns a Texture.
func NewTexture(img Image) *Texture {
	id := GL.CreateTexture()

	GL.BindTexture(GL.TEXTURE_2D, id)

	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_S, GL.CLAMP_TO_EDGE)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_T, GL.CLAMP_TO_EDGE)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MIN_FILTER, GL.LINEAR)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MAG_FILTER, GL.NEAREST)

	if img.Data() == nil {
		print("SDF")
	}

	GL.TexImage2D(GL.TEXTURE_2D, 0, GL.RGBA, img.Width(), img.Height(), 0, GL.RGBA, GL.UNSIGNED_BYTE, img.Data())

	return &Texture{id, img.Width(), img.Height(), GL.LINEAR, GL.LINEAR, GL.CLAMP_TO_EDGE, GL.CLAMP_TO_EDGE}
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

	/*
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
	*/

	var data interface{}
	err := json.Unmarshal([]byte(path), &data)
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
	GL.DeleteTexture(t.id)
}

// Bind will bind the texture.
func (t *Texture) Bind() {
	GL.BindTexture(GL.TEXTURE_2D, t.id)
}

// Unbind will unbind all textures.
func (t *Texture) Unbind() {
	GL.BindTexture(GL.TEXTURE_2D, nil)
}

// Width returns the width of the texture.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() int {
	return t.height
}

// SetFilter sets the filter type used when scaling a texture up or
// down. The default is nearest which will not doing any interpolation
// between pixels.
func (t *Texture) SetFilter(min, max int) {
	t.minFilter = min
	t.maxFilter = max
	t.Bind()
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MIN_FILTER, min)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_MAG_FILTER, max)
}

// Returns the current min and max filters used.
func (t *Texture) Filter() (int, int) {
	return t.minFilter, t.maxFilter
}

func (t *Texture) SetWrap(u, v int) {
	t.uWrap = u
	t.vWrap = v
	t.Bind()
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_S, u)
	GL.TexParameteri(GL.TEXTURE_2D, GL.TEXTURE_WRAP_T, v)
}

func (t *Texture) Wrap() (int, int) {
	return t.uWrap, t.vWrap
}
