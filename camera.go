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

	centerDiference(&cam.Point, Width(), Height(), space)
	// cam.Point = space.Position
}

func (cam *Camera) Update(dt float32) {
	// maxPoints := 8
	if cam.tracking != nil {
		var space *SpaceComponent
		if !cam.tracking.GetComponent(&space) {
			return
		}
		// lerp := float32(.09)
		log.Println(cam.Point)
		centerDiference(&cam.Point, Width(), Height(), space)

	}
}

func centerDiference(to *Point, width, height float32, space *SpaceComponent) {
	to.X += ((space.Position.X + space.Width/2) - (to.X + width/2))
	to.Y += ((space.Position.Y + space.Height/2) - (to.Y + height/2))
}
