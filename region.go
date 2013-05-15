// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	gl "github.com/chsc/gogl/gl33"
	"math"
)

// A region represents a portion of a texture that can be rendered
// using a Batch.
type Region struct {
	texture       *Texture
	u, v          gl.Float
	u2, v2        gl.Float
	width, height int
}

// NewRegion constructs an image from the rectangle x, y, w, h on the
// given texture.
func NewRegion(texture *Texture, x, y, w, h int) *Region {
	invTexWidth := 1.0 / float32(texture.Width())
	invTexHeight := 1.0 / float32(texture.Height())

	u := float32(x) * invTexWidth
	v := float32(y) * invTexHeight
	u2 := float32(x+w) * invTexWidth
	v2 := float32(y+h) * invTexHeight
	width := int(math.Abs(float64(w)))
	height := int(math.Abs(float64(h)))

	return &Region{texture, gl.Float(u), gl.Float(v), gl.Float(u2), gl.Float(v2), width, height}
}

// Flip will swap the region's image on the x and/or y axes.
func (r *Region) Flip(x, y bool) {
	if x {
		tmp := r.u
		r.u = r.u2
		r.u2 = tmp
	}
	if y {
		tmp := r.v
		r.v = r.v2
		r.v2 = tmp
	}
}

func (r *Region) Width() float32 {
	return float32(r.width)
}

func (r *Region) Height() float32 {
	return float32(r.height)
}
