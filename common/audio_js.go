//+build netgo

package common

import (
	"fmt"
	"io"
	"log"

	"engo.io/audio"
	"engo.io/ecs"
	"engo.io/engo"
)

// Load processes the data stream and parses it as an audio file
func (i *audioLoader) Load(url string, data io.Reader) error {

	player, err := audio.NewPlayer(engo.Files.GetRoot()+"/"+url, 0, 0)
	if err != nil {
		return fmt.Errorf("%s (are you running `core.AudioSystemPreload()` before preloading .wav files?)", err.Error())
	}

	i.audios[url] = AudioResource{Player: player, url: url}
	return nil
}

func (a *AudioSystem) New(*ecs.World) {
	a.cachedVolume = MasterVolume

	if a.HeightModifier == 0 {
		a.HeightModifier = defaultHeightModifier
	}

	if !audioSystemPreloaded {
		AudioSystemPreload()
	}
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
					e.AudioComponent.player.Seek(0)
					e.AudioComponent.player.Stop()
					// Remove it from this system, defer because we want to be sure it doesn't interfere with
					// looping over a.entities
					defer a.Remove(*e.BasicEntity)
					continue
				}
			}

			e.AudioComponent.player.Play()
		}
	}
}
