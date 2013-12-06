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
	id     gl.Uint
	region *Region
	width  int
	height int
}

// NewCanvas constructs a canvas and backing texture with the given
// width and height.
func NewCanvas(width, height int) *Canvas {
	canvas := new(Canvas)
	canvas.width = width
	canvas.height = height

	texture := NewTexture(image.NewRGBA(image.Rect(0, 0, width, height)))
	texture.SetFilter(FilterLinear, FilterLinear)
	texture.SetWrap(WrapClampToEdge, WrapClampToEdge)

	gl.GenFramebuffers(1, &canvas.id)

	texture.Bind()
	gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture.id, 0)

	result := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	texture.Unbind()
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	if result != gl.FRAMEBUFFER_COMPLETE {
		gl.DeleteFramebuffers(1, &canvas.id)
		log.Fatal("canvas couldn't be constructed")
	}

	canvas.region = NewRegion(texture, 0, 0, canvas.width, canvas.height)
	canvas.region.Flip(false, true)

	return canvas
}

// Begin should be called before doing any rendering to the canvas.
func (c *Canvas) Begin() {
	gl.Viewport(0, 0, gl.Sizei(c.Width()), gl.Sizei(c.Height()))
	gl.BindFramebuffer(gl.FRAMEBUFFER, c.id)
}

// End should be called when done rendering to the canvas.
func (c *Canvas) End() {
	gl.Viewport(0, 0, gl.Sizei(Width()), gl.Sizei(Height()))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Width is the width of the canvas.
func (c *Canvas) Width() int {
	return int(c.region.Width())
}

// Height is the height of the canvas.
func (c *Canvas) Height() int {
	return int(c.region.Height())
}

// Region returns the backing texture wrapped in a Region for
// rendering with a batch.
func (c *Canvas) Region() *Region {
	return c.region
}
