// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	//"github.com/errcw/glow/gl-core/3.2/gl"
	"github.com/errcw/glow/gl/2.1/gl"
	"image"
	"log"
)

// A Canvas technically wraps an opengl framebuffer. It is used to
// render to a texture that can then be rendered multiple times with a batch.
type Canvas struct {
	id     uint32
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
	texture.SetFilter(gl.LINEAR, gl.LINEAR)
	texture.SetWrap(gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

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
	gl.Viewport(0, 0, int32(c.Width()), int32(c.Height()))
	gl.BindFramebuffer(gl.FRAMEBUFFER, c.id)
}

// End should be called when done rendering to the canvas.
func (c *Canvas) End() {
	gl.Viewport(0, 0, int32(Width()), int32(Height()))
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
