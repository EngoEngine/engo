package common

import (
	"io"
	"log"

	"engo.io/ecs"
	"github.com/hajimehoshi/oto"
)

// MasterVolume provides a number that all underlying player volumes are scaled by
// the value can only be set from 0 to 1
var MasterVolume float64 = 1

// AudioComponent is a Component which is used by the AudioSystem
type AudioComponent struct {
	Repeat  bool
	Player  *Player
	playing bool
	volume  float64
}

// SetVolume sets the AudioComponent's volume
// volume can only be set from 0 to 1
func (ac *AudioComponent) SetVolume(volume float64) {
	if volume <= 0 || volume >= 1 {
		log.Println("Volume can only be set between zero and one. Volume was not set.")
		return
	}
	ac.volume = volume
	ac.Player.setVolume(volume * MasterVolume)
}

type audioEntity struct {
	*ecs.BasicEntity
	*AudioComponent
}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	entities []audioEntity

	OtoPlayer *oto.Player
}

// New is called when the AudioSystem is added to the world. If you use multiple scenes
// make sure you add a Hide method to it and close the OtoPlayer. To be completely safe, also
// add Exit methods to your scenes that close it.
func (a *AudioSystem) New(w *ecs.World) {
	var err error
	a.OtoPlayer, err = oto.NewPlayer(SampleRate, channelNum, bytesPerSample, 8192)
	if err != nil {
		log.Printf("audio error. Unable to create new OtoPlayer: %v \n\r", err)
	}
}

// Add adds a new entity to the AudioSystem. AudioComponent is always required, and the SpaceComponent is
// required as soon as AudioComponent.Background is false. (So if it's not a background noise, we want to know
// where it's originated from)
func (a *AudioSystem) Add(basic *ecs.BasicEntity, audio *AudioComponent) {
	audio.volume = 1
	a.entities = append(a.entities, audioEntity{basic, audio})
}

// AddByInterface adds an entity to the system using the Audioable interface. This allows for entities to be added without specifying each component
func (a *AudioSystem) AddByInterface(o Audioable) {
	a.Add(o.GetBasicEntity(), o.GetAudioComponent())
}

// Remove removes an entity from the AudioSystem
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

// Update is called once per frame, and updates/plays the entities in the AudioSystem
func (a *AudioSystem) Update(dt float32) {
	for _, e := range a.entities {
		e.playing = thePlayers.hasPlayer(e.Player)
	}

	if _, err := io.CopyN(a.OtoPlayer, thePlayers, 4096); err != nil {
		log.Printf("error copying to OtoPlayer: %v \r\n", err)
	}

	for _, e := range a.entities {
		if e.AudioComponent.Repeat && e.playing != thePlayers.hasPlayer(e.Player) {
			e.Player.Rewind()
			e.Player.Play()
		}
	}
}
