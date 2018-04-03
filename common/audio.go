package common

import (
	"log"
	"runtime"

	"engo.io/ecs"

	"github.com/hajimehoshi/oto"
)

const (
	channelNum     = 2
	bytesPerSample = 2

	mask = ^(channelNum*bytesPerSample - 1)
)

// AudioComponent is a Component which is used by the AudioSystem
type AudioComponent struct {
	Player *Player
}

type audioEntity struct {
	*ecs.BasicEntity
	*AudioComponent
}

// AudioSystem is a System that allows for sound effects and / or music
type AudioSystem struct {
	entities []audioEntity

	otoPlayer *oto.Player
}

// New is called when the AudioSystem is added to the world.
func (a *AudioSystem) New(w *ecs.World) {
	var err error
	a.otoPlayer, err = oto.NewPlayer(SampleRate, channelNum, bytesPerSample, 8192)
	if err != nil {
		log.Printf("audio error. Unable to create new OtoPlayer: %v \n\r", err)
	}
	runtime.SetFinalizer(a.otoPlayer, func(p *oto.Player) {
		if err := p.Close(); err != nil {
			log.Printf("audio error. Unable to close OtoPlayer: %v \n\r", err)
		}
	})
	masterVolume = 1
}

// Add adds an entity to the AudioSystem
func (a *AudioSystem) Add(basic *ecs.BasicEntity, audio *AudioComponent) {
	a.entities = append(a.entities, audioEntity{basic, audio})
}

// AddByInterface Allows an Entity to be added directly using the Audioable interface,
// which every entity containing the BasicEntity and AnimationComponent anonymously,
// automatically satisfies.
func (a *AudioSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Audioable)
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

// Update is called once per frame, and updates/plays the players in the AudioSystem
func (a *AudioSystem) Update(dt float32) {

	buf := make([]byte, 4096)
	a.read(buf)

	if _, err := a.otoPlayer.Write(buf); err != nil {
		log.Printf("error copying to OtoPlayer: %v \r\n", err)
	}
}

// Read reads from all the currently playing entities and combines them into a
// single stream that is passed to the oto player.
func (a *AudioSystem) read(b []byte) (int, error) {
	players := make([]*Player, 0)
	for _, e := range a.entities {
		if e.Player.isPlaying {
			players = append(players, e.Player)
		}
	}

	if len(players) == 0 {
		l := len(b)
		l &= mask
		copy(b, make([]byte, l))
		return l, nil
	}

	l := len(b)
	l &= mask

	b16s := [][]int16{}
	for _, player := range players {
		buf, err := player.bufferToInt16(l)
		if err != nil {
			return 0, err
		}
		b16s = append(b16s, buf)
	}
	for i := 0; i < l/2; i++ {
		x := 0
		for _, b16 := range b16s {
			x += int(b16[i])
		}
		if x > (1<<15)-1 {
			x = (1 << 15) - 1
		}
		if x < -(1 << 15) {
			x = -(1 << 15)
		}
		b[2*i] = byte(x)
		b[2*i+1] = byte(x >> 8)
	}

	for _, player := range players {
		if player.eof() {
			if player.Repeat {
				player.Rewind()
			} else {
				player.Pause()
			}
		}
	}

	return l, nil
}
