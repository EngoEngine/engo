//+build netgo

package act

var btnIdCounter = uintptr(0)

type Button struct {
	id  uintptr
	ref *button
}

func newButton(act []Code) Button {
	obj := new(button)
	obj.codes = act
	btnIdCounter++

	return Button{
		id:  btnIdCounter,
		ref: obj,
	}
}

func (this Button) Id() uintptr {
	return this.id
}
