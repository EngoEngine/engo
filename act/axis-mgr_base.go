//+build !netgo

package act

import "unsafe"

type Axis struct {
	ref *axis
}

func newAxis(act []AxisPair) Axis {
	obj := new(axis)
	obj.pairs = act

	return Axis{ref: obj}
}

func (this Axis) Id() uintptr {
	return uintptr(unsafe.Pointer(this.ref))
}
