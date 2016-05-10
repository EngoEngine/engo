package core

import (
	"image"
	"image/draw"
	_ "image/png"
	"io"

	"engo.io/engo"
	"engo.io/gl"
)

type TextureResource struct {
	Texture *gl.Texture
	Width   float32
	Height  float32
	url     string
}

func (t TextureResource) URL() string {
	return t.url
}

type imageLoader struct {
	images map[string]TextureResource
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

func (i *imageLoader) Load(url string, data io.Reader) error {
	img, _, err := image.Decode(data)
	if err != nil {
		return err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return NewTexture(&ImageObject{newm}), nil
}

func (i *imageLoader) Unload(url string) {
	delete(i.images, url)
}

func (i *imageLoader) Resource(url string) (engo.Resource, bool) {
	texture, ok := i.images[url]
	return texture, ok
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

func NewTexture(img Image) TextureResource {
	var id *gl.Texture
	if !engo.Headless() {
		id = engo.Gl.CreateTexture()

		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, id)

		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, engo.Gl.CLAMP_TO_EDGE)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, engo.Gl.CLAMP_TO_EDGE)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MIN_FILTER, engo.Gl.LINEAR)
		engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MAG_FILTER, engo.Gl.NEAREST)

		if img.Data() == nil {
			panic("Texture image data is nil.")
		}

		engo.Gl.TexImage2D(engo.Gl.TEXTURE_2D, 0, engo.Gl.RGBA, engo.Gl.RGBA, engo.Gl.UNSIGNED_BYTE, img.Data())
	}

	return TextureResource{Texture: id, Width: float32(img.Width()), Height: float32(img.Height())}
}

// ImageObject is a pure Go implementation of a `Drawable`
type ImageObject struct {
	data *image.NRGBA
}

// NewImageObject creates a new ImageObject given the image.NRGBA reference
func NewImageObject(img *image.NRGBA) *ImageObject {
	return &ImageObject{img}
}

// Data returns the entire image.NRGBA object
func (i *ImageObject) Data() interface{} {
	return i.data
}

// Width returns the maximum X coordinate of the image
func (i *ImageObject) Width() int {
	return i.data.Rect.Max.X
}

// Height returns the maximum Y coordinate of the image
func (i *ImageObject) Height() int {
	return i.data.Rect.Max.Y
}

func init() {
	engo.Files.Register(".jpg", &imageLoader{})
	engo.Files.Register(".png", &imageLoader{})
	engo.Files.Register(".gif", &imageLoader{})
}
