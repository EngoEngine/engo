package engi

import (
	"log"
)

type Camera struct {
	Point
	tracking *Entity
}

func (cam *Camera) FollowEntity(entity *Entity) {
	cam.tracking = entity
	var space *SpaceComponent
	if !cam.tracking.GetComponent(&space) {
		return
	}

	centerCam(&cam.Point, Width(), Height(), 1, space)
	// cam.Point = space.Position
}

func (cam *Camera) Update(dt float32) {
	// maxPoints := 8
	if cam.tracking != nil {
		var space *SpaceComponent
		if !cam.tracking.GetComponent(&space) {
			return
		}
		log.Println(cam.Point)
		centerCam(&cam.Point, Width(), Height(), 0.09, space)
	}
}

func centerCam(to *Point, width, height, lerp float32, space *SpaceComponent) {
	to.X += ((space.Position.X + space.Width/2) - (to.X + width/2)) * lerp
	to.Y += ((space.Position.Y + space.Height/2) - (to.Y + height/2)) * lerp
}
