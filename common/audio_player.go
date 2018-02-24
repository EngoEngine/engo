package common

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"engo.io/engo"
	"engo.io/engo/common/internal/decode/convert"
)

// SampleRate is the sample rate at which the player plays audio. Any audios
// resource that is added to the system is resampled to this sample rate. To
// change the sample rate, you must do so BEFORE adding the audio system to the world.
var SampleRate = 44100

// Player holds the underlying audio data and plays/pauses/stops/rewinds/seeks it.
type Player struct {
	isPlaying bool
	Repeat    bool

	src        convert.ReadSeekCloser
	url        string
	srcEOF     bool
	sampleRate int

	buf    []byte
	pos    int64
	volume float64

	closeCh         chan struct{}
	closedCh        chan struct{}
	readLoopEndedCh chan struct{}
	seekCh          chan seekArgs
	seekedCh        chan error
	proceedCh       chan []int16
	proceededCh     chan proceededValues
	syncCh          chan func()
}

type seekArgs struct {
	offset int64
	whence int
}

type proceededValues struct {
	buf []int16
	err error
}

// URL implements the engo.Resource interface. It retrieves the player's source url.
func (p *Player) URL() string {
	return p.url
}

func newPlayer(src convert.ReadSeekCloser, url string) (*Player, error) {
	p := &Player{
		src:             src,
		url:             url,
		sampleRate:      SampleRate,
		buf:             []byte{},
		volume:          1,
		closeCh:         make(chan struct{}),
		closedCh:        make(chan struct{}),
		readLoopEndedCh: make(chan struct{}),
		seekCh:          make(chan seekArgs),
		seekedCh:        make(chan error),
		proceedCh:       make(chan []int16),
		proceededCh:     make(chan proceededValues),
		syncCh:          make(chan func()),
	}
	// Get the current position of the source.
	pos, err := p.src.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	p.pos = pos
	runtime.SetFinalizer(p, (*Player).Close)

	go func() {
		p.readLoop()
	}()
	return p, nil
}

// Close removes the player from the audio system's players, which are currently playing players.
// it then finalizes and frees the data from the player.
func (p *Player) Close() error {
	runtime.SetFinalizer(p, nil)
	p.isPlaying = false

	select {
	case p.closeCh <- struct{}{}:
		<-p.closedCh
		return nil
	case <-p.readLoopEndedCh:
		return fmt.Errorf("audio: the player is already closed")
	}
}

func (p *Player) bufferToInt16(lengthInBytes int) ([]int16, error) {
	select {
	case p.proceedCh <- make([]int16, lengthInBytes/2):
		r := <-p.proceededCh
		return r.buf, r.err
	case <-p.readLoopEndedCh:
		return nil, fmt.Errorf("audio: the player is already closed")
	}
}

// Play plays the player's audio.
func (p *Player) Play() {
	p.isPlaying = true
}

func (p *Player) readLoop() {
	defer func() {
		// Note: the error is ignored
		p.src.Close()
		// Receiving from a closed channel returns quickly
		// i.e. `case <-p.readLoopEndedCh:` can check if this loops is ended.
		close(p.readLoopEndedCh)
	}()

	t := time.After(0)
	var readErr error
	for {
		select {
		case <-p.closeCh:
			p.closedCh <- struct{}{}
			return

		case s := <-p.seekCh:
			pos, err := p.src.Seek(s.offset, s.whence)
			p.buf = nil
			p.pos = pos
			p.srcEOF = false
			p.seekedCh <- err
			t = time.After(time.Millisecond)
			break

		case <-t:
			// If the buffer has 1 second, that's enough.
			if len(p.buf) >= p.sampleRate*bytesPerSample*channelNum {
				t = time.After(100 * time.Millisecond)
				break
			}

			// Try to read the buffer for 1/60[s].
			s := 60
			if engo.CurrentBackEnd == engo.BackEndWeb {
				s = 20
				if engo.IsAndroidChrome() {
					s = 10
				}
			}
			l := p.sampleRate * bytesPerSample * channelNum / s
			l &= mask
			buf := make([]byte, l)
			n, err := p.src.Read(buf)

			p.buf = append(p.buf, buf[:n]...)
			if err == io.EOF {
				p.srcEOF = true
			}
			if p.srcEOF && len(p.buf) == 0 {
				t = nil
				break
			}
			if err != nil && err != io.EOF {
				readErr = err
				t = nil
				break
			}
			if engo.CurrentBackEnd == engo.BackEndWeb {
				t = time.After(10 * time.Millisecond)
			} else {
				t = time.After(time.Millisecond)
			}

		case buf := <-p.proceedCh:
			if readErr != nil {
				p.proceededCh <- proceededValues{buf, readErr}
				return
			}

			lengthInBytes := len(buf) * 2
			l := lengthInBytes

			if len(p.buf) < lengthInBytes && !p.srcEOF {
				p.proceededCh <- proceededValues{buf, nil}
				break
			}
			if l > len(p.buf) {
				l = len(p.buf)
			}
			for i := 0; i < l/2; i++ {
				buf[i] = int16(p.buf[2*i]) | (int16(p.buf[2*i+1]) << 8)
				buf[i] = int16(float64(buf[i]) * p.volume)
			}
			p.pos += int64(l)
			p.buf = p.buf[l:]

			p.proceededCh <- proceededValues{buf, nil}

		case f := <-p.syncCh:
			f()
		}
	}
}

func (p *Player) sync(f func()) bool {
	ch := make(chan struct{})
	ff := func() {
		f()
		close(ch)
	}
	select {
	case p.syncCh <- ff:
		<-ch
		return true
	case <-p.readLoopEndedCh:
		return false
	}
}

func (p *Player) eof() bool {
	r := false
	p.sync(func() {
		r = p.srcEOF && len(p.buf) == 0
	})
	return r
}

// IsPlaying returns boolean indicating whether the player is playing.
func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

// Rewind rewinds the current position to the start.
//
// Rewind returns error when seeking the source stream returns error.
func (p *Player) Rewind() error {
	return p.Seek(0)
}

// Seek seeks the position with the given offset.
//
// Seek returns error when seeking the source stream returns error.
func (p *Player) Seek(offset time.Duration) error {
	o := int64(offset) * bytesPerSample * channelNum * int64(p.sampleRate) / int64(time.Second)
	o &= mask
	select {
	case p.seekCh <- seekArgs{o, io.SeekStart}:
		return <-p.seekedCh
	case <-p.readLoopEndedCh:
		return fmt.Errorf("audio: the player is already closed")
	}
}

// Pause pauses the playing.
func (p *Player) Pause() {
	p.isPlaying = false
}

// Current returns the current position.
func (p *Player) Current() time.Duration {
	sample := int64(0)
	p.sync(func() {
		sample = p.pos / bytesPerSample / channelNum
	})
	return time.Duration(sample) * time.Second / time.Duration(p.sampleRate)
}

// GetVolume gets the Player's volume
func (p *Player) GetVolume() float64 {
	v := 0.0
	p.sync(func() {
		v = p.volume
	})
	return v
}

// SetVolume sets the Player's volume
// volume can only be set from 0 to 1
func (p *Player) SetVolume(volume float64) {
	// The condition must be true when volume is NaN.
	if !(0 <= volume && volume <= 1) {
		log.Println("Volume can only be set between zero and one. Volume was not set.")
		return
	}

	p.sync(func() {
		p.volume = volume * masterVolume
	})
}

var masterVolume float64

// SetMasterVolume sets the master volume. The masterVolume is multiplied by all
// the other volumes to get the volume of each entity played.
// Value must be between 0 and 1 or else it doesn't set.
func SetMasterVolume(volume float64) {
	if volume <= 0 || volume >= 1 {
		log.Println("Master Volume can only be set between zero and one. Volume was not set.")
		return
	}
	masterVolume = volume
}

// GetMasterVolume gets the master volume of the audio system.
func GetMasterVolume() float64 {
	return masterVolume
}
