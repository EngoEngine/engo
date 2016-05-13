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
