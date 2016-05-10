//+build !windows,!netgo,!android

package core

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"golang.org/x/mobile/exp/audio/al"
	"io"
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
	player     *Player
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

func (a *AudioSystem) New(w *ecs.World) {
	a.cachedVolume = MasterVolume

	if a.HeightModifier == 0 {
		a.HeightModifier = defaultHeightModifier
	}

	if err := al.OpenDevice(); err != nil {
		log.Println("Error initializing AudioSystem:", err)
		return
	}

	var cam *cameraSystem
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *cameraSystem:
			cam = sys
		}
	}

	if cam == nil {
		log.Println("ERROR: CameraSystem not found - have you added the `RenderSystem` before the `AudioSystem`?")
		return
	}

	// TODO: does this break by any chance, if we use multiple scenes? (w/o recreating world)
	engo.Mailbox.Listen("CameraMessage", func(msg engo.Message) {
		_, ok := msg.(CameraMessage)
		if !ok {
			return
		}

		// Hopefully not that much of an issue, when we receive it before the CameraSystem does
		// TODO: but it is when the CameraMessage is not Incremental (i.e. the changes are big)
		al.SetListenerPosition(al.Vector{cam.X() / engo.GameWidth(), cam.Y() / engo.GameHeight(), cam.Z() * a.HeightModifier})
	})
}

func (a *AudioSystem) Update(dt float32) {
	for _, e := range a.entities {
		if e.AudioComponent.player == nil {
			f := engo.Files.Sound(e.AudioComponent.File)
			if f == nil {
				log.Println("Audio file not loaded:", e.AudioComponent.File)
				continue
			}

			var err error
			e.AudioComponent.player, err = NewPlayer(f, 0, 0)
			if err != nil {
				log.Println("Error initializing AudioComponent:", err)
				continue
			}
		}

		if MasterVolume != a.cachedVolume {
			e.AudioComponent.SetVolume(e.AudioComponent.RawVolume)
		}

		if e.AudioComponent.player.State() != Playing {
			if e.AudioComponent.player.State() == Stopped {
				if !e.AudioComponent.Repeat {
					al.RewindSources(e.AudioComponent.player.source)
					al.StopSources(e.AudioComponent.player.source)
					// Remove it from this system, defer because we want to be sure it doesn't interfere with
					// looping over a.entities
					defer a.Remove(*e.BasicEntity)
					continue
				}
			}

			// Prepares if the track hasn't been buffered before.
			if err := e.AudioComponent.player.prepare(e.AudioComponent.Background, 0, false); err != nil {
				log.Println("Error initializing AudioComponent:", err)
				continue
			}

			al.PlaySources(e.AudioComponent.player.source)

			if !e.AudioComponent.Background {
				e.AudioComponent.player.source.SetPosition(al.Vector{
					(e.SpaceComponent.Position.X + e.SpaceComponent.Width/2) / engo.GameWidth(),
					(e.SpaceComponent.Position.Y + e.SpaceComponent.Height/2) / engo.GameHeight(),
					0,
				})
			}
		}
	}
}
