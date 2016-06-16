package engo

import "engo.io/engo/act"

type Mouse struct {
	X, Y             float32
	ScrollX, ScrollY float32
	Action           Action
	Button           MouseButton
	Modifer          Modifier
}

// Nitya Note: The mouse doc below seems off, in the old system
// there was no way to detect mouse clicks with axis or buttons ?

// InputManager contains information about all forms of input.
type InputManager struct {
	// Mouse is InputManager's reference to the mouse. It is recommended to use the
	// Axis and Button system if at all possible.
	Mouse Mouse

	ActMgr    *act.ActMgr
	AxisMgr   *act.AxisMgr
	ButtonMgr *act.ButtonMgr
}

// NewInputManager holds onto anything input related for engo
func NewInputManager() *InputManager {
	mgr := act.NewActMgr()
	return &InputManager{
		ActMgr:    mgr,
		AxisMgr:   act.NewAxisMgr(mgr),
		ButtonMgr: act.NewButtonMgr(mgr),
	}
}

func (ref *InputManager) clear() {
	ref.ActMgr.Clear()
}

func (ref *InputManager) update() {
	ref.ActMgr.Update()
}
