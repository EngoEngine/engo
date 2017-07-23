package engo

import "engo.io/engo/act"

type Mouse struct {
	X, Y             float32
	ScrollX, ScrollY float32
	Action           Action
	Button           MouseButton
	Modifer          Modifier
	Vertical         AxisMouse
	Horizontal       AxisMouse
}

// InputMgr contains information about all forms of input.
type InputMgr struct {
	// Mouse is InputMgr's reference to the mouse. It is recommended to use the
	// Axis and Button system if at all possible.
	Mouse Mouse

	acts    *act.ActManager
	axes    *act.AxisManager
	buttons *act.ButtonManager
}

// NewInputMgr holds onto anything input related for engo
func NewInputMgr() *InputMgr {
	mgr := act.NewActManager()
	obj := &InputMgr{
		acts:    mgr,
		axes:    act.NewAxisManager(mgr),
		buttons: act.NewButtonManager(mgr),
	}

	obj.Mouse.Vertical.direction = AxisMouseVert
	obj.Mouse.Horizontal.direction = AxisMouseHori

	return obj
}

func (ref *InputMgr) Axes() *act.AxisManager {
	return ref.axes
}

func (ref *InputMgr) Buttons() *act.ButtonManager {
	return ref.buttons
}

func (ref *InputMgr) clear() {
	ref.acts.Clear()
}

func (ref *InputMgr) update() {
	ref.acts.Update()
}

func (ref *InputMgr) Idle(act act.Code) bool {
	return ref.acts.Idle(act)
}

func (ref *InputMgr) Active(act act.Code) bool {
	return ref.acts.Active(act)
}

func (ref *InputMgr) JustIdle(act act.Code) bool {
	return ref.acts.JustIdle(act)
}

func (ref *InputMgr) JustActive(act act.Code) bool {
	return ref.acts.JustActive(act)
}

func (ref *InputMgr) State(act act.Code) act.State {
	return ref.acts.State(act)
}

func (ref *InputMgr) SetState(act act.Code, state bool) {
	ref.acts.SetState(act, state)
}

const (
	// AxisMouseVert is vertical mouse axis
	AxisMouseVert uint32 = 0
	// AxisMouseHori is vertical mouse axis
	AxisMouseHori uint32 = 1
)

// AxisMouse is an axis for a single x or y component of the Mouse. The value returned from it is
// the delta movement, since the previous call and it is not constrained by the AxisMin and AxisMax values.
type AxisMouse struct {
	// direction is the value storing either AxisMouseVert and AxisMouseHori. It determines which directional
	// component to operate on.
	direction uint32
	// old is the delta from the previous calling of Value.
	old float32
}

// Value returns the delta of a mouse movement.
func (am *AxisMouse) Value() float32 {
	var diff float32

	if am.direction == AxisMouseHori {
		diff = Input.Mouse.X - am.old
		am.old = Input.Mouse.X
	} else {
		diff = Input.Mouse.Y - am.old
		am.old = Input.Mouse.Y
	}

	return diff
}
