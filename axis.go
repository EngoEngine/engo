package engo

// An Axis is an input which is a spectrum of values. An example of this is the horizontal movement in a game, or how far a joystick is pressed.
type Axis struct {
	// Name represents the name of the axis (Horizontal, Vertical)
	Name string
	// Pairs represents the axis pairs of this acis
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

// AxisMouseDirection is the direction (X or Y) which the mouse is being tracked for.
type AxisMouseDirection uint

const (
	// AxisMouseVert is vertical mouse axis
	AxisMouseVert AxisMouseDirection = 0
	// AxisMouseHori is vertical mouse axis
	AxisMouseHori AxisMouseDirection = 1
)

// AxisMouse is an axis for a single x or y component of the Mouse. The value returned from it is
// the delta movement, since the previous call and it is not constrained by the AxisMin and AxisMax values.
type AxisMouse struct {
	// direction is the value storing either AxisMouseVert and AxisMouseHori. It determines which directional
	// component to operate on.
	direction AxisMouseDirection
	// old is the delta from the previous calling of Value.
	old float32
}

// NewAxisMouse creates a new Mouse Axis in either direction AxisMouseVert or AxisMouseHori.
func NewAxisMouse(d AxisMouseDirection) *AxisMouse {
	old := Input.Mouse.Y
	if d == AxisMouseHori {
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
		diff = (Input.Mouse.X - am.old + (ResizeXOffset / (2 * GetGlobalScale().X * CanvasScale())))
		am.old = (Input.Mouse.X + (ResizeXOffset / (2 * GetGlobalScale().X * CanvasScale())))
	} else {
		diff = (Input.Mouse.Y - am.old + (ResizeYOffset / (2 * GetGlobalScale().Y * CanvasScale())))
		am.old = (Input.Mouse.Y + (ResizeYOffset / (2 * GetGlobalScale().Y * CanvasScale())))
	}

	return diff
}
