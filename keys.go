package engi

var (
	KEY_STATE_UP        string = "up"
	KEY_STATE_DOWN      string = "down"
	KEY_STATE_JUST_DOWN string = "justdown"
	KEY_STATE_JUST_UP   string = "justup"

	states map[Key]bool

	Keys KeyManager
)

type KeyManager struct {
	KEY_W       KeyState
	KEY_A       KeyState
	KEY_S       KeyState
	KEY_D       KeyState
	KEY_UP      KeyState
	KEY_DOWN    KeyState
	KEY_LEFT    KeyState
	KEY_RIGHT   KeyState
	KEY_SPACE   KeyState
	KEY_CONTROL KeyState
	KEY_ESCAPE  KeyState
	SHIFT       KeyState
	KEY_PLUS    KeyState
	KEY_MINUS   KeyState
}

type KeyState struct {
	lastState    bool
	currentState bool
}

func (key *KeyState) set(state bool) {
	key.lastState = key.currentState
	key.currentState = state
}

func (key *KeyState) State() string {
	if !key.lastState && key.currentState {
		return KEY_STATE_JUST_DOWN
	} else if key.lastState && !key.currentState {
		return KEY_STATE_JUST_UP
	} else if key.lastState && key.currentState {
		return KEY_STATE_DOWN
	} else if !key.lastState && !key.currentState {
		return KEY_STATE_UP
	}

	return KEY_STATE_UP
}

func (key KeyState) JustPressed() bool {
	return key.State() == KEY_STATE_JUST_DOWN
}

func (key KeyState) JustReleased() bool {
	return key.State() == KEY_STATE_JUST_UP
}

func (key KeyState) Up() bool {
	return key.State() == KEY_STATE_UP
}

func (key KeyState) Down() bool {
	return key.State() == KEY_STATE_DOWN
}

func keysUpdate() {
	Keys.KEY_W.set(states[W])
	Keys.KEY_A.set(states[A])
	Keys.KEY_S.set(states[S])
	Keys.KEY_D.set(states[D])

	Keys.KEY_UP.set(states[ArrowUp])
	Keys.KEY_DOWN.set(states[ArrowDown])
	Keys.KEY_LEFT.set(states[ArrowLeft])
	Keys.KEY_RIGHT.set(states[ArrowRight])

	Keys.KEY_SPACE.set(states[Space])
	Keys.KEY_ESCAPE.set(states[Escape])
	Keys.KEY_CONTROL.set(states[LeftControl])
	Keys.SHIFT.set(states[LeftShift])

	Keys.KEY_PLUS.set(states[NumAdd])
	Keys.KEY_MINUS.set(states[NumSubtract])
}
