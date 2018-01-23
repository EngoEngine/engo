//+build netgo windows android

package common

import (
	"io"

	"engo.io/ecs"
)

const (
	defaultHeightModifier float32 = 1
)

// MasterVolume provides a number that all underlying player volumes are scaled by
var MasterVolume float64 = 1

// ReadSeekCloser is an io.ReadSeeker and io.Closer.
type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

// AudioComponent is a Component which is used by the AudioSystem
type AudioComponent struct {
	File       string
	Repeat     bool
	Background bool
	RawVolume  float64
}

// AudioSystemPreload enables the preloading of `.wav` files - has no effect here (because AudioSystem is unimplemented)
func AudioSystemPreload() {}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	HeightModifier float32
}

// New is not implemented
func (as *AudioSystem) New(*ecs.World) {
	notImplemented("audio")
}

// Add is not implemented
func (as *AudioSystem) Add(*ecs.BasicEntity, *AudioComponent, *SpaceComponent) {}

// AddByInterface is not implemented
func (as *AudioSystem) AddByInterface(o Audioable) {}

// Remove is no implemented
func (as *AudioSystem) Remove(basic ecs.BasicEntity) {}

// Update is not implemented
func (as *AudioSystem) Update(dt float32) {}
