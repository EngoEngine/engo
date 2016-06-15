//+build netgo

package act

var axiIdCounter = uintptr(0)

type Axis struct {
	id  uintptr
	ref *axis
}

////////////////

func newAxis(act []AxisPair) Axis {
	obj := new(axis)
	obj.pairs = act
	axiIdCounter++

	return Axis{
		id:  axiIdCounter,
		ref: obj,
	}
}

////////////////

func (this Axis) Id() uintptr {
	return this.id
}
