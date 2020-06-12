//+build !vulkan

package common

import (
	"image/color"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
)

type TextureID *gl.Texture

func createTextureID(img Image) TextureID {
	id := engo.Gl.CreateTexture()

	engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, id)

	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, engo.Gl.CLAMP_TO_EDGE)
	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, engo.Gl.CLAMP_TO_EDGE)
	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MIN_FILTER, engo.Gl.LINEAR)
	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_MAG_FILTER, engo.Gl.NEAREST)

	if img.Data() == nil {
		panic("Texture image data is nil.")
	}

	engo.Gl.TexImage2D(engo.Gl.TEXTURE_2D, 0, engo.Gl.RGBA, engo.Gl.RGBA, engo.Gl.UNSIGNED_BYTE, img.Data())
	return id
}

type BufferData struct {
	// Buffer represents the buffer object itself
	// Avoid using it unless your are writing a custom shader
	Buffer *gl.Buffer
	// BufferContent contains the buffer data
	// Avoid using it unless your are writing a custom shader
	BufferContent []float32
}

func clearScreen() {
	engo.Gl.Clear(engo.Gl.COLOR_BUFFER_BIT)
}

func setBackground(c color.Color) {
	r, g, b, a := c.RGBA()

	engo.Gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
}

func enableMultisample() {
	engo.Gl.Enable(engo.Gl.MULTISAMPLE)
}
