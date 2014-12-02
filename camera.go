package engi

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

	cam.Point = space.Position
}

func (cam *Camera) Update(dt float32) {
	// maxPoints := 8
	if cam.tracking != nil {
		var space *SpaceComponent
		if !cam.tracking.GetComponent(&space) {
			return
		}
		lerp := float32(.07)
		cam.X += ((space.Position.X + space.Width/2) - (cam.X + Width()/2)) * lerp
		cam.Y += ((space.Position.Y + space.Height/2) - (cam.Y + Height()/2)) * lerp
	}
}
