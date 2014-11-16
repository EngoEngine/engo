package engi

type Component interface {
	Name() string
}

type PositionComponent struct {
	X, Y int
}

func (pc PositionComponent) Name() string {
	return "Position"
}

type SpaceComponent struct {
	Position Point
	Width    float32
	Height   float32
}

func (sc SpaceComponent) Name() string {
	return "SpaceComponent"
}
