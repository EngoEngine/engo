//+build !windows,!netgo,!android

package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"engo.io/audio"
	"engo.io/ecs"
	"engo.io/engo"
	"golang.org/x/mobile/exp/audio/al"
)

// Load processes the data stream and parses it as an audio file
func (i *audioLoader) Load(url string, data io.Reader) error {
	audioBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	audioBuffer := bytes.NewReader(audioBytes)
	player, err := audio.NewPlayer(&readSeekCloserBuffer{audioBuffer}, 0, 0)
	if err != nil {
		return fmt.Errorf("%s (are you running `core.AudioSystemPreload()` before preloading .wav files?)", err.Error())
	}

	i.audios[url] = AudioResource{Player: player, url: url}
	return nil
}

func (a *AudioSystem) New(w *ecs.World) {
	a.cachedVolume = MasterVolume

	if a.HeightModifier == 0 {
		a.HeightModifier = defaultHeightModifier
	}

	if !audioSystemPreloaded {
		AudioSystemPreload()
	}

	var cam *cameraSystem
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *cameraSystem:
			cam = sys
		}
	}

	if cam == nil {
		log.Println("[ERROR] CameraSystem not found - have you added the `RenderSystem` before the `AudioSystem`?")
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
			playerRes, err := engo.Files.Resource(e.AudioComponent.File)
			if err != nil {
				log.Println("[ERROR] [AudioSystem]:", err)
				continue // with other entities
			}

			player, ok := playerRes.(AudioResource)
			if !ok {
				log.Println("[ERROR] [AudioSystem]: Loaded audio file is not of type `AudioResource`:", e.AudioComponent.File)
				continue // with other entities
			}

			e.AudioComponent.player = player.Player
		}

		if MasterVolume != a.cachedVolume {
			e.AudioComponent.SetVolume(e.AudioComponent.RawVolume)
		}

		if e.AudioComponent.player.State() != audio.Playing {
			if e.AudioComponent.player.State() == audio.Stopped {
				if !e.AudioComponent.Repeat {
					al.RewindSources(e.AudioComponent.player.Source())
					al.StopSources(e.AudioComponent.player.Source())
					// Remove it from this system, defer because we want to be sure it doesn't interfere with
					// looping over a.entities
					defer a.Remove(*e.BasicEntity)
					continue
				}
			}

			// Prepares if the track hasn't been buffered before.
			if err := e.AudioComponent.player.Prepare(e.AudioComponent.Background, 0, false); err != nil {
				log.Println("Error initializing AudioComponent:", err)
				continue
			}

			al.PlaySources(e.AudioComponent.player.Source())

			if !e.AudioComponent.Background {
				e.AudioComponent.player.Source().SetPosition(al.Vector{
					(e.SpaceComponent.Position.X + e.SpaceComponent.Width/2) / engo.GameWidth(),
					(e.SpaceComponent.Position.Y + e.SpaceComponent.Height/2) / engo.GameHeight(),
					0,
				})
			}
		}
	}
}
