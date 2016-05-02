package engo

// Button is an input which can be either JustPressed, JustReleased or Down. Common uses would be for, a jump key or an action key.
type Button struct {
	Triggers []Key
	Name     string
}

func (b Button) JustPressed() bool {
	for _, trigger := range b.Triggers {
		v := Input.keys.Get(trigger).JustPressed()
		if v {
			return v
		}
	}

	return false
}

func (b Button) JustReleased() bool {
	for _, trigger := range b.Triggers {
		v := Input.keys.Get(trigger).JustReleased()
		if v {
			return v
		}
	}

	return false
}

func (b Button) Down() bool {
	for _, trigger := range b.Triggers {
		v := Input.keys.Get(trigger).Down()
		if v {
			return v
		}
	}

	return false
}
