package engo

const (
	AxisMax     float32 = 1
	AxisNeutral float32 = 0
	AxisMin     float32 = -1
)

// NewInputManager holds onto anything input related for engo
func NewInputManager() *InputManager {
	return &InputManager{
		axes:    make(map[string]Axis),
		buttons: make(map[string]Button),
		keys:    NewKeyManager(),
	}
}

// InputManager contains information about all forms of input.
type InputManager struct {
	axes    map[string]Axis
	buttons map[string]Button
	keys    *KeyManager
}

func (im *InputManager) update() {
	im.keys.update()
}

// RegisterAxis registers a new axis which can be used to retrieve inputs which are spectrums.
func (im *InputManager) RegisterAxis(name string, pairs ...AxisPair) {
	im.axes[name] = Axis{
		Name:  name,
		Pairs: pairs,
	}
}

// RegisterButton registers a new button input.
func (im *InputManager) RegisterButton(name string, keys ...Key) {
	im.buttons[name] = Button{
		Triggers: keys,
		Name:     name,
	}
}

// Axis retrieves an Axis with a specified name.
func (im *InputManager) Axis(name string) Axis {
	return im.axes[name]
}

// Button retrieves a Button with a specified name.
func (im *InputManager) Button(name string) Button {
	return im.buttons[name]
}
