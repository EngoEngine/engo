package act

const cfgAxisMapSize = 16

type AxisManager struct {
	mgr *ActManager

	nameMap map[string]uintptr
	infoMap map[uintptr]Axis
}

type axis struct {
	pairs []AxisPair
}

type AxisPair struct {
	Min Code
	Max Code
}

func NewAxisManager(mgr *ActManager) *AxisManager {
	obj := new(AxisManager)

	obj.mgr = mgr

	obj.nameMap = make(map[string]uintptr, cfgAxisMapSize)
	obj.infoMap = make(map[uintptr]Axis, cfgAxisMapSize)

	return obj
}

func (ref *AxisManager) Id(name string) uintptr {
	return ref.nameMap[name]
}

func (ref *AxisManager) SetById(id uintptr, act ...AxisPair) bool {
	if axi, ok := ref.infoMap[id]; ok {
		axi.ref.pairs = act
		return true
	}
	return false
}

func (ref *AxisManager) SetByName(name string, act ...AxisPair) uintptr {
	if id, ok := ref.nameMap[name]; ok {
		axi := ref.infoMap[id].ref
		axi.pairs = act
		return id
	} else {
		axi := newAxis(act)

		id := axi.Id()
		ref.nameMap[name] = id
		ref.infoMap[id] = axi
		return id
	}
}

func (ref *AxisManager) Value(id uintptr) float32 {
	mgr := ref.mgr
	min := float32(0.0)
	max := float32(0.0)
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if mgr.Active(act.Min) {
			min = -1.0
		}
		if mgr.Active(act.Max) {
			max = 1.0
		}
	}
	return (min + max)
}

func (ref *AxisManager) Idle(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if !mgr.Idle(act.Min) {
			return false
		}
		if !mgr.Idle(act.Max) {
			return false
		}
	}
	return true
}

func (ref *AxisManager) Active(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if mgr.Active(act.Min) {
			return true
		}
		if mgr.Active(act.Max) {
			return true
		}
	}
	return false
}

func (ref *AxisManager) JustIdle(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		min := mgr.State(act.Min)
		if StateJustIdle == min {
			res = true
		} else if StateIdle != min {
			return false
		}
		max := mgr.State(act.Max)
		if StateJustIdle == max {
			res = true
		} else if StateIdle != max {
			return false
		}
	}
	return res
}

func (ref *AxisManager) JustActive(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		min := mgr.State(act.Min)
		if StateJustActive == min {
			res = true
		} else if StateIdle != min {
			return false
		}
		max := mgr.State(act.Max)
		if StateJustActive == max {
			res = true
		} else if StateIdle != max {
			return false
		}
	}
	return res
}

func (ref *AxisManager) MinIdle(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if !mgr.Idle(act.Min) {
			return false
		}
	}
	return true
}

func (ref *AxisManager) MinActive(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if mgr.Active(act.Min) {
			return true
		}
	}
	return false
}

func (ref *AxisManager) MinJustIdle(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		min := mgr.State(act.Min)
		if StateJustIdle == min {
			res = true
		} else if StateIdle != min {
			return false
		}
	}
	return res
}

func (ref *AxisManager) MinJustActive(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		min := mgr.State(act.Min)
		if StateJustActive == min {
			res = true
		} else if StateIdle != min {
			return false
		}
	}
	return res
}

func (ref *AxisManager) MaxIdle(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if !mgr.Idle(act.Max) {
			return false
		}
	}
	return true
}

func (ref *AxisManager) MaxActive(id uintptr) bool {
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		if mgr.Active(act.Max) {
			return true
		}
	}
	return false
}

func (ref *AxisManager) MaxJustIdle(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		max := mgr.State(act.Max)
		if StateJustIdle == max {
			res = true
		} else if StateIdle != max {
			return false
		}
	}
	return res
}

func (ref *AxisManager) MaxJustActive(id uintptr) bool {
	res := false
	mgr := ref.mgr
	axi := ref.infoMap[id].ref
	for _, act := range axi.pairs {
		max := mgr.State(act.Max)
		if StateJustActive == max {
			res = true
		} else if StateIdle != max {
			return false
		}
	}
	return res
}
