package engo

import (
	"log"
	"math"

	"engo.io/ecs"
)

type AABB struct {
	Min, Max Point
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

// Center positions the space component according to its center instead of its
// top-left point (this avoids doing the same math each time in your systems)
func (sc *SpaceComponent) Center(p Point) {
	xDelta := sc.Width / 2
	yDelta := sc.Height / 2
	// update position according to point being used as our center
	sc.Position.X = p.X - xDelta
	sc.Position.Y = p.Y - yDelta
}

func (sc SpaceComponent) AABB() AABB {
	return AABB{Min: sc.Position, Max: Point{sc.Position.X + sc.Width, sc.Position.Y + sc.Height}}
}

type CollisionComponent struct {
	Solid, Main bool
	Extra       Point
}

type CollisionMessage struct {
	Entity collisionEntity
	To     collisionEntity
}

func (collision CollisionMessage) Type() string {
	return "CollisionMessage"
}

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
		offset := Point{e1.CollisionComponent.Extra.X / 2, e1.CollisionComponent.Extra.Y / 2}
		entityAABB.Min.X -= offset.X
		entityAABB.Min.Y -= offset.Y
		entityAABB.Max.X += offset.X
		entityAABB.Max.Y += offset.Y

		for i2, e2 := range cs.entities {
			if i1 == i2 {
				continue // with other entities, because we won't collide with ourselves
			}

			otherAABB := e2.SpaceComponent.AABB()
			offset = Point{e2.CollisionComponent.Extra.X / 2, e2.CollisionComponent.Extra.Y / 2}
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

				Mailbox.Dispatch(CollisionMessage{Entity: e1, To: e2})
			}
		}
	}
}

func IsIntersecting(rect1 AABB, rect2 AABB) bool {
	if rect1.Max.X > rect2.Min.X && rect1.Min.X < rect2.Max.X && rect1.Max.Y > rect2.Min.Y && rect1.Min.Y < rect2.Max.Y {
		return true
	}

	return false
}

func MinimumTranslation(rect1 AABB, rect2 AABB) Point {
	mtd := Point{}

	left := float64(rect2.Min.X - rect1.Max.X)
	right := float64(rect2.Max.X - rect1.Min.X)
	top := float64(rect2.Min.Y - rect1.Max.Y)
	bottom := float64(rect2.Max.Y - rect1.Min.Y)

	if left > 0 || right < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesnt intercept
	}

	if top > 0 || bottom < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesnt intercept
	}
	if math.Abs(left) < right {
		mtd.X = float32(left)
	} else {
		mtd.X = float32(right)
	}

	if math.Abs(top) < bottom {
		mtd.Y = float32(top)
	} else {
		mtd.Y = float32(bottom)
	}

	if math.Abs(float64(mtd.X)) < math.Abs(float64(mtd.Y)) {
		mtd.Y = 0
	} else {
		mtd.X = 0
	}

	return mtd
}
