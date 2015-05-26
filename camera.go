// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import "math"

// A rather basic camera
type Camera struct {
	DeadzoneSize Point
	pos, to      Point
	tracking     *Entity // The entity that is currently being followed
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

	dWidth := cam.DeadzoneSize.X
	dHeight := cam.DeadzoneSize.Y

	if dWidth == 0 {
		dWidth = 200
	}

	if dHeight == 0 {
		dHeight = 200
	}

	min, max := Point{}, Point{}

	min.X, min.Y = cam.to.X-(dWidth/2), cam.to.Y-(dHeight/2)
	max.X, max.Y = min.X+dWidth, min.Y+dHeight

	if cam.pos.X < min.X {
		cam.pos.X = Clamp(floorFloat32(min.X), WorldBounds.Min.X, WorldBounds.Max.X-width)
	} else if cam.pos.X > max.X {
		cam.pos.X = Clamp(floorFloat32(max.X), WorldBounds.Min.X, WorldBounds.Max.X-width)
	}

	if cam.pos.Y < min.X {
		cam.pos.Y = Clamp(floorFloat32(min.Y), WorldBounds.Min.Y, WorldBounds.Max.Y-height)
	} else if cam.pos.Y > max.Y {
		cam.pos.Y = Clamp(floorFloat32(max.Y), WorldBounds.Min.Y, WorldBounds.Max.Y-height)
	}
}

func floorFloat32(i float32) float32 {
	return float32(math.Floor(float64(i)))
}
