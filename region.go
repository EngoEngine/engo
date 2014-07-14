// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"math"
)

// A region represents a portion of a texture that can be rendered
// using a Batch.
type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
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

	return &Region{texture, u, v, u2, v2, width, height}
}

// NewRegionFull returns a region that covers the entire texture.
func NewRegionFull(texture *Texture) *Region {
	return NewRegion(texture, 0, 0, int(texture.Width()), int(texture.Height()))
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
