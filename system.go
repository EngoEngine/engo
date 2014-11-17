package engi

import (
// "log"
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
		// log.Println("Youre in the club", space)
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

type RenderSystem struct {
	*System
	batch *Batch
}

func (rs *RenderSystem) New() {
	rs.System = &System{}
	rs.batch = NewBatch(Width(), Height())
}

func (rs RenderSystem) Pre() {
	Gl.Clear(Gl.COLOR_BUFFER_BIT)
	rs.batch.Begin()
}

func (rs RenderSystem) Post() {
	rs.batch.End()
}

func (rs *RenderSystem) Update(entity *Entity, dt float32) {
	render, hasRender := entity.GetComponent("RenderComponent").(*RenderComponent)
	space, hasSpace := entity.GetComponent("SpaceComponent").(*SpaceComponent)
	if hasRender && hasSpace {
		switch render.Display.(type) {
		case Drawable:
			drawable := render.Display.(Drawable)
			rs.batch.Draw(drawable, space.Position.X, space.Position.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, 0xffffff, 1)
		case *Font:
			font := render.Display.(*Font)
			font.Print(rs.batch, render.Label, space.Position.X, space.Position.Y, 0xffffff)
		}
	}
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}
