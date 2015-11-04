package engi

type Component interface {
	Type() string
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

func (SpaceComponent) Type() string {
	return "SpaceComponent"
}

func (sc SpaceComponent) AABB() AABB {
	return AABB{Min: sc.Position, Max: Point{sc.Position.X + sc.Width, sc.Position.Y + sc.Height}}
}

type CollisionMasterComponent struct {
}

func (CollisionMasterComponent) Type() string {
	return "CollisionMasterComponent"
}

type CollisionComponent struct {
	Solid, Main bool
	Extra       Point
}

func (CollisionComponent) Type() string {
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

func NewRenderComponent(display Renderable, scale Point, label string) *RenderComponent {
	return &RenderComponent{
		Display:      display,
		Scale:        scale,
		Label:        label,
		Priority:     MiddleGround,
		Transparency: 1,
		Color:        0xffffff,
	}
}

func (RenderComponent) Type() string {
	return "RenderComponent"
}

type LinkComponent struct {
	Entity *Entity
}

func (LinkComponent) Type() string {
	return "LinkComponent"
}
