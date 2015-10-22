package engi

import (
	"sync"
)

var (
	KEY_STATE_UP        string = "up"
	KEY_STATE_DOWN      string = "down"
	KEY_STATE_JUST_DOWN string = "justdown"
	KEY_STATE_JUST_UP   string = "justup"

	keyStates map[Key]bool

	Keys KeyManager
)

type KeyManager struct {
	// TODO: benchmark to figure out if an array (with unused items) is faster
	mapper map[Key]KeyState

	// TODO: backwards compatible; remove at some point
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

	mutex sync.RWMutex
}

func (km *KeyManager) Get(k Key) KeyState {
	km.mutex.RLock()
	defer km.mutex.RUnlock()
	ks, ok := km.mapper[k]
	if !ok {
		return KeyState{false, false}
	}
	return ks
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
	Keys.mutex.Lock()
	for key, down := range keyStates {
		ks := Keys.mapper[key]
		ks.set(down)
		Keys.mapper[key] = ks
	}
	Keys.mutex.Unlock()

	// TODO: backwards compatible; remove at some point
	Keys.KEY_W.set(keyStates[W])
	Keys.KEY_A.set(keyStates[A])
	Keys.KEY_S.set(keyStates[S])
	Keys.KEY_D.set(keyStates[D])

	Keys.KEY_UP.set(keyStates[ArrowUp])
	Keys.KEY_DOWN.set(keyStates[ArrowDown])
	Keys.KEY_LEFT.set(keyStates[ArrowLeft])
	Keys.KEY_RIGHT.set(keyStates[ArrowRight])

	Keys.KEY_SPACE.set(keyStates[Space])
	Keys.KEY_ESCAPE.set(keyStates[Escape])
	Keys.KEY_CONTROL.set(keyStates[LeftControl])
	Keys.SHIFT.set(keyStates[LeftShift])

	Keys.KEY_PLUS.set(keyStates[NumAdd])
	Keys.KEY_MINUS.set(keyStates[NumSubtract])
}

func init() {
	Keys.mapper = make(map[Key]KeyState)
}
