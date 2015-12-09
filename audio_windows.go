package engi

import (
	"log"

	"github.com/paked/engi/ecs"
)

const (
	defaultHeightModifier float32 = 1
)

// AudioComponent is a Component which is used by the AudioSystem
type AudioComponent struct {
	File       string
	Repeat     bool
	Background bool
	player     *Player
}

func (*AudioComponent) Type() string {
	return "AudioComponent"
}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	*ecs.System
	HeightModifier float32
}

func (AudioSystem) Type() string {
	return "AudioSystem"
}

func (as *AudioSystem) New(*ecs.World) {
	as.System = ecs.NewSystem()

	log.Println("Warning: audio is not yet implemented on Windows")
}

func (as *AudioSystem) Update(entity *ecs.Entity, dt float32) {}
