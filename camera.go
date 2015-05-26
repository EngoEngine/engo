// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import "math"

// A rather basic camera
type Camera struct {
	pos, to  Point
	tracking *Entity // The entity that is currently being followed
}

func (cam *Camera) FollowEntity(entity *Entity) {
	cam.tracking = entity
	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		return
	}

	cam.to = space.Position
	cam.centerCam(Width(), Height(), 1, space)
}

func (cam *Camera) Update(dt float32) {
	// maxPoints := 8
	if cam.tracking != nil {
		var space *SpaceComponent
		if !cam.tracking.GetComponent(&space) {
			return
		}
		cam.centerCam(Width(), Height(), 0.09, space)
	}
}

func (cam *Camera) centerCam(width, height, lerp float32, space *SpaceComponent) {
	cam.to.X += ((space.Position.X + space.Width/2) - (cam.to.X + width/2)) * lerp
	cam.to.Y += ((space.Position.Y + space.Height/2) - (cam.to.Y + height/2)) * lerp

	if cam.to.X > 0 && cam.to.X+width < WorldBounds.Max.X {
		cam.pos.X = floorFloat32(cam.to.X)
	}

	if cam.to.Y > 0 && cam.to.Y+height < WorldBounds.Max.Y {
		cam.pos.Y = floorFloat32(cam.pos.Y)
	}

}

func floorFloat32(i float32) float32 {
	return float32(math.Floor(float64(i)))
}
