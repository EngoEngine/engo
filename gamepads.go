package engo

import "sync"

// Gamepadbutton is a button on a Gamepad.
type GamepadButton struct {
	lastState    bool
	currentState bool
}

func (b *GamepadButton) set(state bool) {
	b.lastState = b.currentState
	b.currentState = state
}

// State returns the raw state of a key.
func (b *GamepadButton) State() int {
	if b.lastState {
		if b.currentState {
			return KeyStateDown
		}
		return KeyStateJustUp
	}
	if b.currentState {
		return KeyStateJustDown
	}
	return KeyStateUp
}

// JustPressed returns whether a key was just pressed
func (b GamepadButton) JustPressed() bool {
	return (!b.lastState && b.currentState)
}

// JustReleased returns whether a key was just released
func (b GamepadButton) JustReleased() bool {
	return (b.lastState && !b.currentState)
}

// Up returns wheter a key is not being pressed
func (b GamepadButton) Up() bool {
	return (!b.lastState && !b.currentState)
}

// Down returns wether a key is being pressed
func (b GamepadButton) Down() bool {
	return (b.lastState && b.currentState)
}

// GamepadManager manages the gamepads
type GamepadManager struct {
	mutex    sync.RWMutex
	gamepads map[string]*Gamepad
}

// NewGamepadManager creates a new GamepadManager
func NewGamepadManager() *GamepadManager {
	return &GamepadManager{
		gamepads: make(map[string]*Gamepad),
	}
}

// RegisterGamepad registers the gamepad with the given name. It can return an
// error if no suitable gamepads are located.
func (gm *GamepadManager) Register(name string) error {
	return gm.registerGamepadImpl(name)
}

// GetGamepad returns the gamepad previously registered with name.
func (gm *GamepadManager) GetGamepad(name string) *Gamepad {
	return gm.gamepads[name]
}

func (gm *GamepadManager) update() {
	gm.updateImpl()
}
