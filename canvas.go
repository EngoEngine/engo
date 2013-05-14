// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	gl "github.com/chsc/gogl/gl33"
	"image"
	"log"
)

// A Canvas technically wraps an opengl framebuffer. It is used to
// render to a texture that can then be rendered multiple times with a batch.
type Canvas struct {
	id      gl.Uint
	texture *Texture
	width   int
	height  int
}

// NewCanvas constructs a canvas and backing texture with the given
// width and height.
func NewCanvas(width, height int) *Canvas {
	canvas := new(Canvas)
	canvas.width = width
	canvas.height = height

	canvas.texture = NewTexture(image.NewRGBA(image.Rect(0, 0, width, height)))
	canvas.texture.SetFilter(FilterLinear, FilterLinear)
	canvas.texture.SetWrap(WrapClampToEdge, WrapClampToEdge)

	gl.GenFramebuffers(1, &canvas.id)

	canvas.texture.Bind()
	gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, canvas.texture.id, 0)

	result := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	canvas.texture.Unbind()
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	if result != gl.FRAMEBUFFER_COMPLETE {
		gl.DeleteFramebuffers(1, &canvas.id)
		log.Fatal("canvas couldn't be constructed")
	}

	return canvas
}

// Begin should be called before doing any rendering to the canvas.
func (c *Canvas) Begin() {
	gl.Viewport(0, 0, gl.Sizei(c.texture.Width()), gl.Sizei(c.texture.Height()))
	gl.BindFramebuffer(gl.FRAMEBUFFER, c.id)
}

// End should be called when done rendering to the canvas.
func (c *Canvas) End() {
	gl.Viewport(0, 0, gl.Sizei(Width()), gl.Sizei(Height()))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Texture returns the backing texture that will be rendered to. This
// can be wrapped in a Region for rendering with a batch. The texture
// will most likely be flipped upside down. Region.Flip(false, true)
// can be used when the region is made to deal with that.
func (c *Canvas) Texture() *Texture {
	return c.texture
}

// Width is the width of the canvas.
func (c *Canvas) Width() int {
	return c.texture.Width()
}

// Height is the height of the canvas.
func (c *Canvas) Height() int {
	return c.texture.Height()
}
