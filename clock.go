package engo

import (
	"time"

	"github.com/luxengine/math"
)

// A Clock is a measurement built in `engo` to measure the actual frames per seconds (framerate).
type Clock struct {
	elapsed float32
	delta   float32
	fps     float32
	frames  uint64
	start   time.Time
	frame   time.Time
}

// NewClock creates a new timer which allows you to measure ticks per seconds. Be sure to call `Tick()` whenever you
// want a tick to occur - it does not automatically tick each frame.
func NewClock() *Clock {
	clock := new(Clock)
	clock.start = time.Now()
	clock.Tick()
	return clock
}

// Tick indicates a new tick/frame has occurred.
func (c *Clock) Tick() {
	now := time.Now()
	c.frames += 1
	if !c.frame.IsZero() {
		c.delta = float32(now.Sub(c.frame).Seconds())
	}

	c.elapsed += c.delta
	c.frame = now

	if c.elapsed >= 1 {
		c.fps = float32(c.frames)
		c.elapsed = math.Mod(c.elapsed, 1)
		c.frames = 0
	}
}

// Delta is the amount of seconds between the last tick and the one before that
func (c *Clock) Delta() float32 {
	return c.delta
}

// FPS is the amount of frames per second, computed every time a tick occurs at least a second after the previous update
func (c *Clock) FPS() float32 {
	return c.fps
}

// Time is the number of seconds the clock has been running
func (c *Clock) Time() float32 {
	return float32(time.Now().Sub(c.start).Seconds())
}
