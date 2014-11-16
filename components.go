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

type CollisionMasterComponent struct {
}

func (cm CollisionMasterComponent) Name() string {
	return "CollisionMasterComponent"
}

func (cm CollisionMasterComponent) Is() bool {
	return true
}
