// +build netgo

package engo

import (
	"engo.io/gl"
	"github.com/gopherjs/gopherjs/js"
)

type Image interface {
	Data() *js.Object
	Width() int
	Height() int
}

func NewTexture(img Image) *Texture {
	var id *gl.Texture
	if !headless {
		id = Gl.CreateTexture()
		Gl.BindTexture(Gl.TEXTURE_2D, id)
		Gl.TexImage2D(Gl.TEXTURE_2D, 0, Gl.RGBA, Gl.RGBA, Gl.UNSIGNED_BYTE, img.Data())
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MAG_FILTER, Gl.LINEAR)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MIN_FILTER, Gl.LINEAR_MIPMAP_NEAREST)
		Gl.GenerateMipmap(Gl.TEXTURE_2D)
		/*
			Gl.BindTexture(Gl.TEXTURE_2D, nil)
			Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_S, Gl.CLAMP_TO_EDGE)
			Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_T, Gl.CLAMP_TO_EDGE)
		*/
		if img.Data() == nil {
			panic("Texture image data is nil.")
		}

	}

	return &Texture{id, float32(img.Width()), float32(img.Height())}
}
