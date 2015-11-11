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
	mapper map[Key]KeyState
	mutex  sync.RWMutex
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
}

func init() {
	Keys.mapper = make(map[Key]KeyState)
}
