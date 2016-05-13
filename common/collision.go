package common

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/luxengine/math"
)

type SpaceComponent struct {
	Position engo.Point
	Width    float32
	Height   float32
	Rotation float32 // angle in degrees for the rotation to apply clockwise
}

// Center positions the space component according to its center instead of its
// top-left point (this avoids doing the same math each time in your systems)
func (sc *SpaceComponent) Center(p engo.Point) {
	xDelta := sc.Width / 2
	yDelta := sc.Height / 2
	// update position according to point being used as our center
	sc.Position.X = p.X - xDelta
	sc.Position.Y = p.Y - yDelta
}

// AABB returns the minimum and maximum point for the given SpaceComponent. It hereby takes into account the
// rotation of the Component - it may very well be that the Minimum as given by engo.AABB, is smaller than the Position
// of the object (i.e. when rotated). As this method takes into account the rotation, it should be used only when
// required.
func (sc SpaceComponent) AABB() engo.AABB {
	if sc.Rotation == 0 {
		return engo.AABB{
			Min: sc.Position,
			Max: engo.Point{sc.Position.X + sc.Width, sc.Position.Y + sc.Height},
		}
	}

	sin, cos := math.Sincos(sc.Rotation * math.Pi / 180)
	xmin := sc.Position.X
	xmax := sc.Position.X + sc.Width*cos - sc.Height*sin
	ymin := sc.Position.Y
	ymax := sc.Position.Y + sc.Height*cos + sc.Width*sin

	var (
		X_MIN, X_MAX, Y_MIN, Y_MAX float32
	)

	if xmin < xmax {
		X_MIN = xmin
		X_MAX = xmax
	} else {
		X_MIN = xmax
		X_MAX = xmin
	}

	if ymin < ymax {
		Y_MIN = ymin
		Y_MAX = ymax
	} else {
		Y_MIN = ymax
		Y_MAX = ymin
	}

	return engo.AABB{engo.Point{X_MIN, Y_MIN}, engo.Point{X_MAX, Y_MAX}}
}

type CollisionComponent struct {
	Solid, Main bool
	Extra       engo.Point
	Collides    bool // Collides is true if the component is colliding with something during this pass
}

type CollisionMessage struct {
	Entity collisionEntity
	To     collisionEntity
}

func (CollisionMessage) Type() string { return "CollisionMessage" }

type collisionEntity struct {
	*ecs.BasicEntity
	*CollisionComponent
	*SpaceComponent
}

type CollisionSystem struct {
	entities []collisionEntity
}

func (c *CollisionSystem) Add(basic *ecs.BasicEntity, collision *CollisionComponent, space *SpaceComponent) {
	c.entities = append(c.entities, collisionEntity{basic, collision, space})
}

func (c *CollisionSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (cs *CollisionSystem) Update(dt float32) {
	for i1, e1 := range cs.entities {
		if !e1.CollisionComponent.Main {
			continue // with other entities
		}

		entityAABB := e1.SpaceComponent.AABB()
		offset := engo.Point{e1.CollisionComponent.Extra.X / 2, e1.CollisionComponent.Extra.Y / 2}
		entityAABB.Min.X -= offset.X
		entityAABB.Min.Y -= offset.Y
		entityAABB.Max.X += offset.X
		entityAABB.Max.Y += offset.Y

		var collided bool

		for i2, e2 := range cs.entities {
			if i1 == i2 {
				continue // with other entities, because we won't collide with ourselves
			}

			otherAABB := e2.SpaceComponent.AABB()
			offset = engo.Point{e2.CollisionComponent.Extra.X / 2, e2.CollisionComponent.Extra.Y / 2}
			otherAABB.Min.X -= offset.X
			otherAABB.Min.Y -= offset.Y
			otherAABB.Max.X += offset.X
			otherAABB.Max.Y += offset.Y

			if IsIntersecting(entityAABB, otherAABB) {
				if e1.CollisionComponent.Solid && e2.CollisionComponent.Solid {
					mtd := MinimumTranslation(entityAABB, otherAABB)
					e1.SpaceComponent.Position.X += mtd.X
					e1.SpaceComponent.Position.Y += mtd.Y
				}

				collided = true
				engo.Mailbox.Dispatch(CollisionMessage{Entity: e1, To: e2})
			}
		}

		e1.CollisionComponent.Collides = collided
	}
}

func IsIntersecting(rect1 engo.AABB, rect2 engo.AABB) bool {
	if rect1.Max.X > rect2.Min.X && rect1.Min.X < rect2.Max.X && rect1.Max.Y > rect2.Min.Y && rect1.Min.Y < rect2.Max.Y {
		return true
	}

	return false
}

func MinimumTranslation(rect1 engo.AABB, rect2 engo.AABB) engo.Point {
	mtd := engo.Point{}

	left := rect2.Min.X - rect1.Max.X
	right := rect2.Max.X - rect1.Min.X
	top := rect2.Min.Y - rect1.Max.Y
	bottom := rect2.Max.Y - rect1.Min.Y

	if left > 0 || right < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesn't intercept
	}

	if top > 0 || bottom < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesn't intercept
	}
	if math.Abs(left) < right {
		mtd.X = left
	} else {
		mtd.X = right
	}

	if math.Abs(top) < bottom {
		mtd.Y = top
	} else {
		mtd.Y = bottom
	}

	if math.Abs(mtd.X) < math.Abs(mtd.Y) {
		mtd.Y = 0
	} else {
		mtd.X = 0
	}

	return mtd
}
