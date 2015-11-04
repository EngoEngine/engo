package engi

type Systemer interface {
	Update(entity *Entity, dt float32)
	Type() string
	Priority() int
	Pre()
	Post()
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
	SetWorld(*World)
	SkipOnHeadless() bool
	// Push(message Message)
	// Receive(message Message)
	// Messages() []Message
	// Dismiss(i int)
}

type System struct {
	entities             []*Entity
	messageQueue         []Message
	World                *World
	ShouldSkipOnHeadless bool
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

func (s *System) SetWorld(w *World) {
	s.World = w
}

func (s System) SkipOnHeadless() bool {
	return s.ShouldSkipOnHeadless
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

				Mailbox.Dispatch(CollisionMessage{Entity: entity, To: other})
			}
		}
	}
}

func (*CollisionSystem) Type() string {
	return "CollisionSystem"
}

type PriorityLevel int

const (
	// HighestGround is the highest PriorityLevel that will be rendered
	HighestGround PriorityLevel = 50
	// HUDGround is a PriorityLevel from which everything isn't being affected by the Camera
	HUDGround    PriorityLevel = 40
	Foreground   PriorityLevel = 30
	MiddleGround PriorityLevel = 20
	ScenicGround PriorityLevel = 10
	// Background is the lowest PriorityLevel that will be rendered
	Background PriorityLevel = 0
)

type RenderSystem struct {
	renders map[PriorityLevel][]*Entity
	changed bool
	*System
}

func (rs *RenderSystem) New() {
	rs.renders = make(map[PriorityLevel][]*Entity)
	rs.System = &System{}
	rs.ShouldSkipOnHeadless = true
}

func (rs *RenderSystem) AddEntity(e *Entity) {
	rs.changed = true
	rs.System.AddEntity(e)
}

func (rs RenderSystem) Pre() {
	if !rs.changed {
		return
	}

	rs.renders = make(map[PriorityLevel][]*Entity)
}

type Renderable interface {
	Render(b *Batch, render *RenderComponent, space *SpaceComponent)
}

func (rs *RenderSystem) Post() {
	var currentBatch *Batch

	for i := Background; i <= HighestGround; i++ {
		if len(rs.renders[i]) == 0 {
			continue
		}

		// Retrieve a batch, may be the default one -- then call .Begin() if we arent already using it
		batch := world.batch(i)
		if batch != currentBatch {
			if currentBatch != nil {
				currentBatch.End()
			}
			batch.Begin()
			currentBatch = batch
		}
		// Then render everything for this level
		for _, entity := range rs.renders[i] {
			var render *RenderComponent
			var space *SpaceComponent

			if !entity.GetComponent(&render) || !entity.GetComponent(&space) {
				return
			}

			render.Display.Render(batch, render, space)
		}
	}

	if currentBatch != nil {
		currentBatch.End()
	}

	rs.changed = false
}

func (rs *RenderSystem) Update(entity *Entity, dt float32) {
	if !rs.changed {
		return
	}

	var render *RenderComponent
	if !entity.GetComponent(&render) {
		return
	}

	rs.renders[render.Priority] = append(rs.renders[render.Priority], entity)
}

func (*RenderSystem) Type() string {
	return "RenderSystem"
}

func (rs *RenderSystem) Priority() int {
	return 1
}
