package common

import (
	"errors"
	"io"
	"log"

	"engo.io/ecs"
	"engo.io/engo"

	"github.com/hajimehoshi/oto"
)

const (
	channelNum     = 2
	bytesPerSample = 2

	mask = ^(channelNum*bytesPerSample - 1)
)

// stepPlayer is used for headless mode audio, such as for tests
// you can control exactly how many steps it takes which allows for verification
// of the PCM data written to it. This replaces oto.Player when run in headless
type stepPlayer struct {
	ThrowWriteError bool
	stepStart       chan []byte
	stepDone        chan struct{}
}

func (l *stepPlayer) Write(b []byte) (int, error) {
	if l.ThrowWriteError {
		return 0, errors.New("write error")
	}
	l.stepStart <- b
	<-l.stepDone

	return len(b), nil
}

func (l *stepPlayer) Close() error {
	return nil
}

func (l *stepPlayer) Bytes() []byte {
	return <-l.stepStart
}

func (l *stepPlayer) Step() {
	if len(l.stepDone) > 0 {
		return
	}
	l.stepDone <- struct{}{}
}

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

	otoPlayer                   io.WriteCloser
	bufsize                     int
	closeCh, pauseCh, restartCh chan struct{}
	playerCh                    chan []*Player
}

// New is called when the AudioSystem is added to the world.
func (a *AudioSystem) New(w *ecs.World) {
	var err error
	switch engo.CurrentBackEnd {
	case engo.BackEndMobile:
		a.bufsize = 12288
	default:
		a.bufsize = 8192
	}
	if engo.Headless() {
		a.otoPlayer = &stepPlayer{
			stepStart: make(chan []byte),
			stepDone:  make(chan struct{}, 1),
		}
	} else {
		a.otoPlayer, err = oto.NewPlayer(SampleRate, channelNum, bytesPerSample, a.bufsize)
		if err != nil {
			log.Printf("audio error. Unable to create new OtoPlayer: %v \n\r", err)
		}
	}
	// run oto on a separate thread so it doesn't slow down updates
	a.closeCh = make(chan struct{}, 1)
	a.pauseCh = make(chan struct{}, 1)
	a.restartCh = make(chan struct{}, 1)
	a.playerCh = make(chan []*Player, 25)
	go func() {
		players := make([]*Player, 0)
	loop:
		for {
			select {
			case <-a.closeCh:
				break loop
			case <-a.pauseCh:
				<-a.restartCh
			case players = <-a.playerCh:
			default:
				buf := make([]byte, 2048)
				a.read(buf, players)

				if _, err := a.otoPlayer.Write(buf); err != nil {
					log.Printf("error copying to OtoPlayer: %v \r\n", err)
				}
			}
		}
	}()
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

// Update doesn't do anything since audio is run on it's own thread
func (a *AudioSystem) Update(dt float32) {
	if len(a.playerCh) >= 25 { //if the channel is full just return so we don't block the update loop
		return
	}
	players := make([]*Player, 0)
	for _, e := range a.entities {
		if e.Player.isPlaying {
			players = append(players, e.Player)
		}
	}
	a.playerCh <- players
}

// Read reads from all the currently playing entities and combines them into a
// single stream that is passed to the oto player.
func (a *AudioSystem) read(b []byte, players []*Player) (int, error) {
	l := len(b)
	l &= mask

	if len(players) == 0 {
		copy(b, make([]byte, l))
		return l, nil
	}

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

// Close closes the AudioSystem's loop. After this is called the AudioSystem
// can no longer play audio.
func (a *AudioSystem) Close() {
	if len(a.closeCh) > 0 { //so it doesn't block
		return
	}
	a.closeCh <- struct{}{}
}

// Pause pauses the AudioSystem's loop. Call Restart to continue playing audio.
func (a *AudioSystem) Pause() {
	if len(a.pauseCh) > 0 { // so it doesn't block
		return
	}
	a.pauseCh <- struct{}{}
}

// Restart restarts the AudioSystem's loop when it's paused.
func (a *AudioSystem) Restart() {
	if len(a.restartCh) > 0 { // so it doesn't block
		return
	}
	a.restartCh <- struct{}{}
}
