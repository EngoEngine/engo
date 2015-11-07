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

// CameraSystem is a System that manages the state of the Camera
type cameraSystem struct {
	*System
	x, y, z  float32
	tracking *Entity // The entity that is currently being followed
}

func (cameraSystem) Type() string {
	return "cameraSystem"
}

func (cam *cameraSystem) New() {
	cam.System = NewSystem()

	cam.x = WorldBounds.Max.X / 2
	cam.y = WorldBounds.Max.Y / 2
	cam.z = 1

	cam.AddEntity(NewEntity([]string{cam.Type()}))

	Mailbox.Listen("CameraMessage", func(msg Message) {
		cammsg, ok := msg.(CameraMessage)
		if !ok {
			return
		}

		if cammsg.Incremental {
			switch cammsg.Axis {
			case XAxis:
				cam.moveX(cammsg.Value)
			case YAxis:
				cam.moveY(cammsg.Value)
			case ZAxis:
				cam.zoom(cammsg.Value)
			}
		} else {
			switch cammsg.Axis {
			case XAxis:
				cam.moveToX(cammsg.Value)
			case YAxis:
				cam.moveToY(cammsg.Value)
			case ZAxis:
				cam.zoomTo(cammsg.Value)
			}
		}
	})
}

func (cam *cameraSystem) FollowEntity(entity *Entity) {
	cam.tracking = entity
	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		cam.tracking = nil
		return
	}
}

func (cam *cameraSystem) moveX(value float32) {
	cam.moveToX(cam.x + value)
}

func (cam *cameraSystem) moveY(value float32) {
	cam.moveToY(cam.y + value)
}

func (cam *cameraSystem) zoom(value float32) {
	cam.zoomTo(cam.z + value)
}

func (cam *cameraSystem) moveToX(location float32) {
	cam.x = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.X)
}

func (cam *cameraSystem) moveToY(location float32) {
	cam.y = mgl32.Clamp(location, WorldBounds.Min.X, WorldBounds.Max.Y)
}

func (cam *cameraSystem) zoomTo(zoomLevel float32) {
	cam.z = mgl32.Clamp(zoomLevel, MinZoom, MaxZoom)
}

func (cam *cameraSystem) X() float32 {
	return cam.x
}

func (cam *cameraSystem) Y() float32 {
	return cam.y
}

func (cam *cameraSystem) Z() float32 {
	return cam.z
}

func (cam *cameraSystem) Update(entity *Entity, dt float32) {
	if cam.tracking == nil {
		return
	}

	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		return
	}

	cam.centerCam(space.Position.X+space.Width/2, space.Position.Y+space.Height/2, cam.z)
}

func (cam *cameraSystem) centerCam(x, y, z float32) {
	cam.moveToX(x)
	cam.moveToY(y)
	cam.zoomTo(z)
}

// CameraAxis is the axis at which the Camera can/has to move
type CameraAxis uint8

const (
	XAxis CameraAxis = iota
	YAxis
	ZAxis
)

// CameraMessage is a message that can be sent to the Camera (and other Systemers), to indicate movement
type CameraMessage struct {
	Axis        CameraAxis
	Value       float32
	Incremental bool
}

func (CameraMessage) Type() string {
	return "CameraMessage"
}
