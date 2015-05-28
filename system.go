package engi

type Systemer interface {
	Update(entity *Entity, dt float32)
	Name() string
	Priority() int
	Pre()
	Post()
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
	// Push(message Message)
	// Receive(message Message)
	// Messages() []Message
	// Dismiss(i int)
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
		if other.ID() != entity.ID() && other.Exists {

			var r *RenderComponent
			other.GetComponent(&r)

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

				Mailbox.Dispatch("CollisionMessage", CollisionMessage{Entity: entity, To: other})
			}
		}
	}
}

func (cs CollisionSystem) Name() string {
	return "CollisionSystem"
}

type PriorityLevel int

const (
	Foreground   PriorityLevel = 2
	MiddleGround PriorityLevel = 1
	Background   PriorityLevel = 0
)

type RenderSystem struct {
	renders map[PriorityLevel][]*Entity
	*System
}

func (rs *RenderSystem) New() {
	rs.renders = make(map[PriorityLevel][]*Entity)
	rs.System = &System{}
}

func (rs RenderSystem) Pre() {
	delete(rs.renders, Foreground)
	delete(rs.renders, MiddleGround)
	delete(rs.renders, Background)
}

func (rs RenderSystem) Post() {
	for i := 2; i >= 0; i-- {
		for _, entity := range rs.renders[PriorityLevel(i)] {
			var render *RenderComponent
			var space *SpaceComponent

			if !entity.GetComponent(&render) || !entity.GetComponent(&space) {
				return
			}

			switch render.Display.(type) {
			case Drawable:
				drawable := render.Display.(Drawable)
				Wo.Batch().Draw(drawable, space.Position.X-Cam.pos.X, space.Position.Y-Cam.pos.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, render.Color, render.Transparency)
			case *Font:
				font := render.Display.(*Font)
				font.Print(Wo.Batch(), render.Label, space.Position.X-Cam.pos.X, space.Position.Y-Cam.pos.Y, 0xffffff)
			case *Text:
				text := render.Display.(*Text)
				text.Draw(Wo.Batch(), space.Position)
			case *Level:
				level := render.Display.(*Level)
				for _, img := range level.Images {
					if img.Image != nil {
						Wo.Batch().Draw(img.Image, img.X-Cam.pos.X, img.Y-Cam.pos.Y, 0, 0, 1, 1, 0, 0xffffff, 1)
					}
				}
				for _, tile := range level.Tiles {
					if tile.Image != nil {
						Wo.Batch().Draw(tile.Image, (tile.X+space.Position.X)-Cam.pos.X, (tile.Y+space.Position.Y)-Cam.pos.Y, 0, 0, 1, 1, 0, 0xffffff, 1)
					}
				}
			}
		}

	}
}

func (rs *RenderSystem) Update(entity *Entity, dt float32) {
	var render *RenderComponent
	if !entity.GetComponent(&render) {
		return
	}

	rs.renders[render.Priority] = append(rs.renders[render.Priority], entity)
	/*var render *RenderComponent
	var space *SpaceComponent

	if !entity.GetComponent(&render) || !entity.GetComponent(&space) {
		return
	}

	switch render.Display.(type) {
	case Drawable:
		drawable := render.Display.(Drawable)
		Wo.Batch().Draw(drawable, space.Position.X-Cam.pos.X, space.Position.Y-Cam.pos.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, 0xffffff, 1)
	case *Font:
		font := render.Display.(*Font)
		font.Print(Wo.Batch(), render.Label, space.Position.X-Cam.pos.X, space.Position.Y-Cam.pos.Y, 0xffffff)
	case *Text:
		text := render.Display.(*Text)
		text.Draw(Wo.Batch(), space.Position)
	case *Level:
		level := render.Display.(*Level)
		for _, img := range level.Images {
			if img.Image != nil {
				Wo.Batch().Draw(img.Image, img.X-Cam.pos.X, img.Y-Cam.pos.Y, 0, 0, 1, 1, 0, 0xffffff, 1)
			}
		}
		for _, tile := range level.Tiles {
			if tile.Image != nil {
				Wo.Batch().Draw(tile.Image, (tile.X+space.Position.X)-Cam.pos.X, (tile.Y+space.Position.Y)-Cam.pos.Y, 0, 0, 1, 1, 0, 0xffffff, 1)
			}
		}
	}*/
}

func (rs RenderSystem) Name() string {
	return "RenderSystem"
}

func (rs RenderSystem) Priority() int {
	return 1
}
