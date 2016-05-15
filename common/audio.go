//+build !android,!windows
package common

import (
	"io"
	"log"

	"engo.io/audio"
	"engo.io/ecs"
)

const (
	defaultHeightModifier float32 = 1
)

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
	player     *audio.Player
	RawVolume  float64
}

func (ac *AudioComponent) SetVolume(volume float64) {
	ac.RawVolume = volume
	ac.player.SetVolume(volume * MasterVolume)
}

type audioEntity struct {
	*ecs.BasicEntity
	*AudioComponent
	*SpaceComponent
}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	entities       []audioEntity
	HeightModifier float32

	cachedVolume float64
}

var audioSystemPreloaded bool

// AudioSystemPreload has to be called before preloading any `.wav` files
func AudioSystemPreload() {
	if err := audio.Preload(); err != nil {
		log.Println("Error initializing AudioSystem:", err)
		return
	}
	audioSystemPreloaded = true
}

// Add adds a new entity to the AudioSystem. AudioComponent is always required, and the SpaceComponent is
// required as soon as AudioComponent.Background is false. (So if it's not a background noise, we want to know
// where it's originated from)
func (a *AudioSystem) Add(basic *ecs.BasicEntity, audio *AudioComponent, space *SpaceComponent) {
	a.entities = append(a.entities, audioEntity{basic, audio, space})
}

func (a *AudioSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range a.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		a.entities = append(a.entities[:delete], a.entities[delete+1:]...)
	}
}
