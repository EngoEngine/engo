package engi

var (
	KEY_STATE_UP        string = "up"
	KEY_STATE_DOWN      string = "down"
	KEY_STATE_JUST_DOWN string = "justdown"
	KEY_STATE_JUST_UP   string = "justup"
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
