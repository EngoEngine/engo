package engi

import (
	"log"
)

type Systemer interface {
	Update(entity *Entity, dt float32)
	Name() string
	Priority() int
	Pre()
	Post()
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
}

type System struct {
	entities []*Entity
}

func (s System) New()  {}
func (s System) Pre()  {}
func (s System) Post() {}

func (s System) Priority() int {
	return 0
}

func (s System) Entities() []*Entity {
	return s.entities
}

func (s *System) AddEntity(entity *Entity) {
	s.entities = append(s.entities, entity)
}

type CollisionSystem struct {
	*System
}

func (cs *CollisionSystem) New() {
	cs.System = &System{}
}

func (cs *CollisionSystem) Update(entity *Entity, dt float32) {
	space, hasSpace := entity.GetComponent("SpaceComponent").(*SpaceComponent)
	_, hasCollisionMaster := entity.GetComponent("CollisionMasterComponent").(*CollisionMasterComponent)
	if hasSpace && hasCollisionMaster {
		log.Println("Youre in the club", space, collisionMaster)
		for _, other := range cs.Entities() {
			if other.ID() != entity.ID() {
				otherSpace, otherHasSpace := other.GetComponent("SpaceComponent").(*SpaceComponent)
				if otherHasSpace {
					entityAABB := space.AABB()
					otherAABB := otherSpace.AABB()
					if IsIntersecting(entityAABB, otherAABB) {
						mtd := MinimumTranslation(entityAABB, otherAABB)
						space.Position.X += mtd.X
						space.Position.Y += mtd.Y
					}
				}
			}
		}
	}
}

func (cs CollisionSystem) Name() string {
	return "CollisionSystem"
}
