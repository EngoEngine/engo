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

////////////////

type ActMgr struct {
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

////////////////

func NewActMgr() *ActMgr {
	obj := new(ActMgr)

	obj.dirtMap = make(map[Code]Code, cfgCodeMapSize)
	obj.codeMap = make(map[Code]actState, cfgCodeMapSize)

	return obj
}

////////////////

func (ref *ActMgr) Clear() {
	for _, act := range ref.dirtMap {
		delete(ref.dirtMap, act)
		st := ref.codeMap[act]

		st.lastState = st.currState
		st.currState = st.currState

		ref.codeMap[act] = st
	}
}

func (ref *ActMgr) Update() {
	for _, cs := range ref.codeLst {
		st := ref.codeMap[cs.act]

		st.lastState = st.currState
		st.currState = cs.state

		ref.dirtMap[cs.act] = cs.act
		ref.codeMap[cs.act] = st
	}
	ref.codeLst = ref.codeLst[:0]
}

func (ref *ActMgr) Idle(act Code) bool {
	st := ref.codeMap[act]
	return (!st.lastState && !st.currState)
}

func (ref *ActMgr) Active(act Code) bool {
	st := ref.codeMap[act]
	return (st.lastState && st.currState)
}

func (ref *ActMgr) JustIdle(act Code) bool {
	st := ref.codeMap[act]
	return (st.lastState && !st.currState)
}

func (ref *ActMgr) JustActive(act Code) bool {
	st := ref.codeMap[act]
	return (!st.lastState && st.currState)
}

////////////////

func (ref *ActMgr) State(act Code) State {
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

func (ref *ActMgr) SetState(act Code, state bool) {
	ref.codeLst = append(ref.codeLst, codeState{act: act, state: state})
}
