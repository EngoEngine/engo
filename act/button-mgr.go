package act

const cfgButtonMapSize = 16

type ButtonMgr struct {
	mgr *ActMgr

	nameMap map[string]uintptr
	infoMap map[uintptr]Button
}

type button struct {
	codes []Code
}

////////////////

func NewButtonMgr(mgr *ActMgr) *ButtonMgr {
	obj := new(ButtonMgr)

	obj.mgr = mgr

	obj.nameMap = make(map[string]uintptr, cfgButtonMapSize)
	obj.infoMap = make(map[uintptr]Button, cfgButtonMapSize)

	return obj
}

////////////////

func (ref *ButtonMgr) Id(name string) uintptr {
	return ref.nameMap[name]
}

func (ref *ButtonMgr) SetId(id uintptr, act ...Code) bool {
	if btn, ok := ref.infoMap[id]; ok {
		btn.ref.codes = act
		return true
	}
	return false
}

func (ref *ButtonMgr) SetNamed(name string, act ...Code) uintptr {
	if id, ok := ref.nameMap[name]; ok {
		btn := ref.infoMap[id].ref
		btn.codes = act
		return id
	} else {
		btn := newButton(act)

		id := btn.Id()
		ref.nameMap[name] = id
		ref.infoMap[id] = btn
		return id
	}
}

////////////////

func (ref *ButtonMgr) Idle(id uintptr) bool {
	mgr := ref.mgr
	btn := ref.infoMap[id].ref
	for _, act := range btn.codes {
		if !mgr.Idle(act) {
			return false
		}
	}
	return true
}

func (ref *ButtonMgr) Active(id uintptr) bool {
	mgr := ref.mgr
	btn := ref.infoMap[id].ref
	for _, act := range btn.codes {
		if mgr.Active(act) {
			return true
		}
	}
	return false
}

func (ref *ButtonMgr) JustIdle(id uintptr) bool {
	res := false
	mgr := ref.mgr
	btn := ref.infoMap[id].ref
	for _, act := range btn.codes {
		state := mgr.State(act)
		if StateJustIdle == state {
			res = true
		} else if StateIdle != state {
			return false
		}
	}
	return res
}

func (ref *ButtonMgr) JustActive(id uintptr) bool {
	res := false
	mgr := ref.mgr
	btn := ref.infoMap[id].ref
	for _, act := range btn.codes {
		state := mgr.State(act)
		if StateJustActive == state {
			res = true
		} else if StateIdle != state {
			return false
		}
	}
	return res
}
