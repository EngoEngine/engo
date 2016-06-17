package act

type Code int32
type State uint8

const cfgCodeMapSize = 16

const KeyCode = Code(1 << 16)
const MouseCode = Code(2 << 16)

const StateIdle = State(0)
const StateActive = State(1)
const StateJustIdle = State(2)
const StateJustActive = State(3)

type ActManager struct {
	codeLst []codeState
	dirtMap map[Code]Code
	codeMap map[Code]actState
}

type actState struct {
	lastState bool
	currState bool
}

type codeState struct {
	act   Code
	state bool
}

func NewActManager() *ActManager {
	obj := new(ActManager)

	obj.dirtMap = make(map[Code]Code, cfgCodeMapSize)
	obj.codeMap = make(map[Code]actState, cfgCodeMapSize)

	return obj
}

func (ref *ActManager) Clear() {
	for _, act := range ref.dirtMap {
		delete(ref.dirtMap, act)
		st := ref.codeMap[act]

		st.lastState = st.currState
		st.currState = st.currState

		ref.codeMap[act] = st
	}
}

func (ref *ActManager) Update() {
	for _, cs := range ref.codeLst {
		st := ref.codeMap[cs.act]

		st.lastState = st.currState
		st.currState = cs.state

		ref.dirtMap[cs.act] = cs.act
		ref.codeMap[cs.act] = st
	}
	ref.codeLst = ref.codeLst[:0]
}

func (ref *ActManager) Idle(act Code) bool {
	st := ref.codeMap[act]
	return (!st.lastState && !st.currState)
}

func (ref *ActManager) Active(act Code) bool {
	st := ref.codeMap[act]
	return (st.lastState && st.currState)
}

func (ref *ActManager) JustIdle(act Code) bool {
	st := ref.codeMap[act]
	return (st.lastState && !st.currState)
}

func (ref *ActManager) JustActive(act Code) bool {
	st := ref.codeMap[act]
	return (!st.lastState && st.currState)
}

func (ref *ActManager) State(act Code) State {
	st := ref.codeMap[act]
	if st.lastState {
		if st.currState {
			return StateActive
		} else {
			return StateJustIdle
		}
	} else {
		if st.currState {
			return StateJustActive
		} else {
			return StateIdle
		}
	}
}

func (ref *ActManager) SetState(act Code, state bool) {
	ref.codeLst = append(ref.codeLst, codeState{act: act, state: state})
}
