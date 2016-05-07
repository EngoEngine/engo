package engo

// An Axis is an input which is a spectrum of values. An example of this is the horizontal movement in a game, or how far a joystick is pressed.
type Axis struct {
	Name  string
	Pairs []AxisPair
}

// Value returns the value of an Axis.
func (a Axis) Value() float32 {
	for _, pair := range a.Pairs {
		v := pair.Value()
		if v != AxisNeutral {
			return v
		}
	}

	return AxisNeutral
}

// An AxisPair is a set of Min/Max values which could possible be used by an Axis.
type AxisPair interface {
	Value() float32
}

// An AxisKeyPair is a set of Min/Max values used for detecting whether or not a key has been pressed.
type AxisKeyPair struct {
	Min Key
	Max Key
}

// Value returns the value of a keypress.
func (keys AxisKeyPair) Value() float32 {
	if Input.keys.Get(keys.Max).Down() {
		return AxisMax
	} else if Input.keys.Get(keys.Min).Down() {
		return AxisMin
	}

	return AxisNeutral
}

const (
	// The vertical mouse axis
	AxisMouseVert = 0
	// The horizontal mouse axis
	AxisMouseHori = 1
)

// AxisMouse is an axis of the mouse direction.
type AxisMouse struct {
	direction int

	old float32
}

// NewAxisMouse creates a new Mouse Axis in either direction AxisMouseVert or AxisMouseHori.
func NewAxisMouse(d int) *AxisMouse {
	old := Input.Mouse.Y
	if old == AxisMouseHori {
		old = Input.Mouse.X
	}

	return &AxisMouse{
		direction: d,
		old:       old,
	}
}

// Value returns the delta of a mouse movement.
func (am *AxisMouse) Value() float32 {
	var diff float32

	if am.direction == AxisMouseHori {
		diff = Input.Mouse.X - am.old

		am.old = Input.Mouse.X
	} else {
		diff = Input.Mouse.Y - am.old

		am.old = Input.Mouse.Y
	}

	return diff
}
