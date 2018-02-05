package common

import (
	"sync"

	"engo.io/engo/common/decode/convert"
)

const (
	channelNum     = 2
	bytesPerSample = 2

	mask = ^(channelNum*bytesPerSample - 1)
)

// players is for all the currently playing audio players
type players struct {
	players map[*Player]struct{}
	sync.RWMutex
}

// thePlayers is the specific instance of all the currently playing audio players
var thePlayers = &players{players: make(map[*Player]struct{})}

func (p *players) Read(b []byte) (int, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.players) == 0 {
		l := len(b)
		l &= mask
		copy(b, make([]byte, l))
		return l, nil
	}

	l := len(b)
	l &= mask

	b16s := [][]int16{}
	for player := range p.players {
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

	closed := []*Player{}
	for player := range p.players {
		if player.eof() {
			closed = append(closed, player)
		}
	}
	for _, player := range closed {
		delete(p.players, player)
	}

	return l, nil
}

func (p *players) addPlayer(player *Player) {
	p.Lock()
	p.players[player] = struct{}{}
	p.Unlock()
}

func (p *players) removePlayer(player *Player) {
	p.Lock()
	delete(p.players, player)
	p.Unlock()
}

func (p *players) hasPlayer(player *Player) bool {
	p.RLock()
	_, ok := p.players[player]
	p.RUnlock()
	return ok
}

func (p *players) hasSource(src convert.ReadSeekCloser) bool {
	p.RLock()
	defer p.RUnlock()
	for player := range p.players {
		if player.src == src {
			return true
		}
	}
	return false
}
