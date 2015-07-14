package engi

type Component interface {
	Name() string
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

func (sc SpaceComponent) Name() string {
	return "SpaceComponent"
}

func (sc SpaceComponent) AABB() AABB {
	return AABB{Min: sc.Position, Max: Point{sc.Position.X + sc.Width, sc.Position.Y + sc.Height}}
}

type PhysicsComponent struct {
	// Is the entity touching the ground?
	Grounded bool

	Gravity  float32
	Velocity Point

	// Ground and Air Friction
	G_Friction float32
	A_Friction float32
}

func (spc PhysicsComponent) Name() string {
	return "PhysicsComponent"
}

type CollisionMasterComponent struct {
}

func (cm CollisionMasterComponent) Name() string {
	return "CollisionMasterComponent"
}

type CollisionComponent struct {
	Solid, Main bool
	Extra       Point
}

func (cc CollisionComponent) Name() string {
	return "CollisionComponent"
}

func (cm CollisionMasterComponent) Is() bool {
	return true
}

type RenderComponent struct {
	Display      Renderable
	Scale        Point
	Label        string
	Priority     PriorityLevel
	Transparency float32
	Color        uint32
}

func NewRenderComponent(display Renderable, scale Point, label string) RenderComponent {
	return RenderComponent{
		Display:      display,
		Scale:        scale,
		Label:        label,
		Priority:     MiddleGround,
		Transparency: 1,
		Color:        0xffffff,
	}
}

func (rc RenderComponent) Name() string {
	return "RenderComponent"
}

type LinkComponent struct {
	Entity *Entity
}

func (lc LinkComponent) Name() string {
	return "LinkComponent"
}
