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

type CollisionMasterComponent struct {
}

func (cm CollisionMasterComponent) Name() string {
	return "CollisionMasterComponent"
}

func (cm CollisionMasterComponent) Is() bool {
	return true
}
