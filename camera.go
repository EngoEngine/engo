// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	MinZoom float32 = 0.25
	MaxZoom float32 = 3
)

// A rather basic camera
type Camera struct {
	x, y, z  float32
	tracking *Entity // The entity that is currently being followed
}

func (cam *Camera) Setup() {
	cam.x = 0
	cam.y = 0
	cam.z = 1
}

func (cam *Camera) FollowEntity(entity *Entity) {
	cam.tracking = entity
	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		cam.tracking = nil
		return
	}
}

func (cam *Camera) MoveX(value float32) {
	cam.MoveToX(cam.x + value)
}

func (cam *Camera) MoveY(value float32) {
	cam.MoveToY(cam.y + value)
}

func (cam *Camera) Zoom(value float32) {
	cam.ZoomTo(cam.z + value)
}

func (cam *Camera) MoveToX(location float32) {
	cam.x = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.X)
}

func (cam *Camera) MoveToY(location float32) {
	cam.y = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.Y)
}

func (cam *Camera) ZoomTo(zoomLevel float32) {
	cam.z = mgl32.Clamp(zoomLevel, MinZoom, MaxZoom)
}

func (cam *Camera) X() float32 {
	return cam.x
}

func (cam *Camera) Y() float32 {
	return cam.y
}

func (cam *Camera) Z() float32 {
	return cam.z
}

func (cam *Camera) Update(dt float32) {
	if cam.tracking == nil {
		return
	}

	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		return
	}

	cam.centerCam(space.Position.X+space.Width/2, space.Position.Y+space.Height/2, cam.z)
}

func (cam *Camera) centerCam(x, y, z float32) {
	cam.MoveToX(x)
	cam.MoveToY(y)
	cam.ZoomTo(z)
}
