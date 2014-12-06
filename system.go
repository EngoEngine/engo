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
	Push(message Message)
	Receive(message Message)
	Messages() []Message
	Dismiss(i int)
}

type System struct {
	entities     []*Entity
	messageQueue []Message
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

func (system *System) Push(message Message) {
	system.messageQueue = append(system.messageQueue, message)
}

func (system System) Receive(message Message) {}

func (system System) Messages() []Message {
	return system.messageQueue
}

func (system *System) Dismiss(i int) {
	log.Println(i, len(system.messageQueue))
	system.messageQueue = system.messageQueue[:i+copy(system.messageQueue[i:], system.messageQueue[i+1:])]
}

type CollisionSystem struct {
	*System
}

func (cs *CollisionSystem) New() {
	cs.System = &System{}
}

func (cs *CollisionSystem) Update(entity *Entity, dt float32) {
	var space *SpaceComponent
	var collisionComponent *CollisionComponent
	if !entity.GetComponent(&space) || !entity.GetComponent(&collisionComponent) {
		return
	}

	if !collisionComponent.Main {
		return
	}

	for _, other := range cs.Entities() {
		if other.ID() != entity.ID() {

			var r *RenderComponent
			other.GetComponent(&r)
			t, ok := r.Display.(*Tilemap)
			if ok {
				CollideTilemap(entity, other, t)
				return
			}

			var otherSpace *SpaceComponent
			var otherCollision *CollisionComponent
			if !other.GetComponent(&otherSpace) || !other.GetComponent(&otherCollision) {
				return
			}

			entityAABB := space.AABB()
			offset := Point{collisionComponent.Extra.X / 2, collisionComponent.Extra.Y / 2}
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
				if otherCollision.Solid && collisionComponent.Solid {
					mtd := MinimumTranslation(entityAABB, otherAABB)
					space.Position.X += mtd.X
					space.Position.Y += mtd.Y
				}

				Mailbox.Dispatch(CollisionMessage{entity})
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
	var render *RenderComponent
	var space *SpaceComponent

	if !entity.GetComponent(&render) || !entity.GetComponent(&space) {
		return
	}

	switch render.Display.(type) {
	case Drawable:
		drawable := render.Display.(Drawable)
		rs.batch.Draw(drawable, space.Position.X-Cam.X, space.Position.Y-Cam.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, 0xffffff, 1)
	case *Font:
		font := render.Display.(*Font)
		font.Print(rs.batch, render.Label, space.Position.X-Cam.X, space.Position.Y-Cam.Y, 0xffffff)
	case *Text:
		text := render.Display.(*Text)
		text.Draw(rs.batch, space.Position)
	case *Tilemap:
		tilemap := render.Display.(*Tilemap)
		for _, slice := range tilemap.Tiles {
			for _, tile := range slice {
				if tile.Image != nil {
					rs.batch.Draw(tile.Image, (tile.X+space.Position.X)-Cam.X, (tile.Y+space.Position.Y)-Cam.Y, 0, 0, 1, 1, 0, 0xffffff, 1)
				}
			}
		}
	}
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}
