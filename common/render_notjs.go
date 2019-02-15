//+build !js

package common

import "unsafe"

// Unsafe interface for instance comparison
type iface struct {
	Type, Data unsafe.Pointer
}

func compareShaders(a, b Shader) bool {
	// comparing the instances using unsafe pointers seems to be a little faster (about 0.007 ns on a "slow" machine)
	// so using this unsafe method to compare shader != prevShader gives a nice performance boost when using many entities.
	aIface := *(*iface)(unsafe.Pointer(&a))
	bIface := *(*iface)(unsafe.Pointer(&b))
	return aIface.Data == bIface.Data
}
