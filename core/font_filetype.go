package core

import (
	"image"
	"image/draw"
	_ "image/png"
	"io"

	"engo.io/engo"
	"engo.io/gl"
	"github.com/golang/freetype/truetype"
)

type TextureResource struct {
	Font *truetype.Font

	Texture *gl.Texture
	Width   float32
	Height  float32
	url     string
}

func (t TextureResource) URL() string {
	return t.url
}

type fontLoader struct {
	fonts map[string]TextureResource
	/*
		// Load loads the given resource into memory.
		Load(url string, data io.Reader) error

		// Unload releases the given resource from memory.
		Unload(url string) error

		// Resource returns the given resource, and a boolean indicating whether the
		// resource was loaded.
		Resource(url string) (Resource, bool)
	*/
}

func (i *fontLoader) Load(url string, data io.Reader) error {
	img, _, err := image.Decode(data)
	if err != nil {
		return err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return NewTexture(&ImageObject{newm}), nil
}

func (i *fontLoader) Unload(url string) {
	delete(i.fonts, url)
}

func (i *fontLoader) Resource(url string) (engo.Resource, bool) {
	texture, ok := i.fonts[url]
	return texture, ok
}

func init() {
	engo.Files.Register(".ttf", &fontLoader{})
}
