package engi

import (
	"github.com/paked/engi/ecs"
	"log"
	"math"
)

type AABB struct {
	Min, Max Point
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

func (*SpaceComponent) Type() string {
	return "SpaceComponent"
}

func (sc SpaceComponent) AABB() AABB {
	return AABB{Min: sc.Position, Max: Point{sc.Position.X + sc.Width, sc.Position.Y + sc.Height}}
}

type CollisionMasterComponent struct {
}

func (*CollisionMasterComponent) Type() string {
	return "CollisionMasterComponent"
}

func (cm CollisionMasterComponent) Is() bool {
	return true
}

type CollisionComponent struct {
	Solid, Main bool
	Extra       Point
}

func (*CollisionComponent) Type() string {
	return "CollisionComponent"
}

type CollisionMessage struct {
	Entity *ecs.Entity
	To     *ecs.Entity
}

func (collision CollisionMessage) Type() string {
	return "CollisionMessage"
}

type CollisionSystem struct {
	*ecs.System
}

func (cs *CollisionSystem) New(*ecs.World) {
	cs.System = ecs.NewSystem()
}

func (cs *CollisionSystem) RunInParallel() bool {
	return len(cs.Entities()) > 40 // turning point for CollisionSystem
}

func (cs *CollisionSystem) Update(entity *ecs.Entity, dt float32) {
	var (
		space     *SpaceComponent
		collision *CollisionComponent
		ok        bool
	)

	if space, ok = entity.ComponentFast(space).(*SpaceComponent); !ok {
		return
	}

	if collision, ok = entity.ComponentFast(collision).(*CollisionComponent); !ok {
		return
	}

	if !collision.Main {
		return
	}

	var (
		otherSpace     *SpaceComponent
		otherCollision *CollisionComponent
	)

	for _, other := range cs.Entities() {
		if other.ID() != entity.ID() {
			if otherSpace, ok = other.ComponentFast(otherSpace).(*SpaceComponent); !ok {
				return
			}

			if otherCollision, ok = other.ComponentFast(otherCollision).(*CollisionComponent); !ok {
				return
			}

			entityAABB := space.AABB()
			offset := Point{collision.Extra.X / 2, collision.Extra.Y / 2}
			entityAABB.Min.X -= offset.X
			entityAABB.Min.Y -= offset.Y
			entityAABB.Max.X += offset.X
			entityAABB.Max.Y += offset.Y
			otherAABB := otherSpace.AABB()
			offset = Point{otherCollision.Extra.X / 2, otherCollision.Extra.Y / 2}
			otherAABB.Min.X -= offset.X
			otherAABB.Min.Y -= offset.Y
			otherAABB.Max.X += offset.X
			otherAABB.Max.Y += offset.Y
			if IsIntersecting(entityAABB, otherAABB) {
				if otherCollision.Solid && collision.Solid {
					mtd := MinimumTranslation(entityAABB, otherAABB)
					space.Position.X += mtd.X
					space.Position.Y += mtd.Y
				}

				Mailbox.Dispatch(CollisionMessage{Entity: entity, To: other})
			}
		}
	}
}

func (*CollisionSystem) Type() string {
	return "CollisionSystem"
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
