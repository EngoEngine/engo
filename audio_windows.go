package engo

import (
	"log"

	"engo.io/ecs"
	"io"
)

const (
	defaultHeightModifier float32 = 1
)

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
}

func (*AudioComponent) Type() string {
	return "AudioComponent"
}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	ecs.LinearSystem
	HeightModifier float32
}

func (AudioSystem) Type() string {
	return "AudioSystem"
}

func (as *AudioSystem) New(*ecs.World) {
	log.Println("Warning: audio is not yet implemented on Windows")
}

func (as *AudioSystem) UpdateEntity(entity *ecs.Entity, dt float32) {}
