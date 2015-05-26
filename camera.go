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

	posX := cam.to.X
	posY := cam.to.Y

	minX, minY := posX-100, posY-100
	maxX, maxY := minX+200, minY+200

	if cam.pos.X < minX {
		cam.pos.X = Clamp(floorFloat32(minX), 0, WorldBounds.Max.X-width)
	}

	if cam.pos.X > maxX {
		cam.pos.X = Clamp(floorFloat32(maxX), 0, WorldBounds.Max.X-width)
	}

	if posY < minY {
		cam.pos.Y = Clamp(floorFloat32(minY), 0, WorldBounds.Max.Y-height)
	}

	if posY > maxY {
		cam.pos.Y = Clamp(floorFloat32(maxY), 0, WorldBounds.Max.Y-height)
	}

}

func floorFloat32(i float32) float32 {
	return float32(math.Floor(float64(i)))
}
