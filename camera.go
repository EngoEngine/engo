// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

// A Camera represents a viewport with a 2d projection. It allows for
// simple transformations and zooming.
type Camera struct {
	Zoom              float32
	Position          *Vector
	Direction         *Vector
	Up                *Vector
	Projection        *Matrix
	View              *Matrix
	Combined          *Matrix
	InvProjectionView *Matrix
	ViewportWidth     float32
	ViewportHeight    float32
}

// NewCamera returns a camera with a viewport of the given width and height.
func NewCamera(width, height float32) *Camera {
	camera := new(Camera)
	camera.Zoom = 1
	camera.Position = new(Vector)
	camera.Direction = &Vector{0, 0, -1}
	camera.Up = &Vector{0, 1, 0}
	camera.Projection = NewMatrix()
	camera.View = NewMatrix()
	camera.Combined = NewMatrix()
	camera.InvProjectionView = NewMatrix()
	camera.ViewportWidth = width
	camera.ViewportHeight = height
	return camera
}

var tmp = new(Vector)

// Update should be called every time the camer is changed in any way.
func (c *Camera) Update() {
	c.Projection.SetToOrtho(c.Zoom*-c.ViewportWidth/2, c.Zoom*c.ViewportWidth/2, c.Zoom*c.ViewportHeight/2, c.Zoom*-c.ViewportHeight/2, 0, 1)
	c.View.SetToLookAt(c.Position, tmp.Set(c.Position).Add(c.Direction), c.Up)
	c.Combined.Set(c.Projection).Mul(c.View)
	c.InvProjectionView.Set(c.Combined).Inv()
}

// Unproject takes a point in screen space and transforms it to be in
// the view space of the camera.
func (c *Camera) Unproject(vec *Vector) {
	viewportWidth := float32(Width())
	viewportHeight := float32(Height())

	x := vec.X
	y := vec.Y
	y = viewportHeight - y - 1
	vec.X = (2*x)/viewportWidth - 1
	vec.Y = (2*y)/viewportHeight - 1
	vec.Z = 2*vec.Z - 1

	vec.Prj(c.InvProjectionView)
}
