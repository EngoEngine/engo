//+build !netgo

package act

import "unsafe"

type Button struct {
	ref *button
}

////////////////

func newButton(act []Code) Button {
	obj := new(button)
	obj.codes = act

	return Button{ref: obj}
}

////////////////

func (this Button) Id() uintptr {
	return uintptr(unsafe.Pointer(this.ref))
}
