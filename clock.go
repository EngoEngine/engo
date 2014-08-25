// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"math"
	"time"
)

type Clock struct {
	elapsed float64
	delta   float64
	fps     float64
	frames  uint64
	start   time.Time
	frame   time.Time
}

func NewClock() *Clock {
	clock := new(Clock)
	clock.start = time.Now()
	clock.Tick()
	return clock
}

func (c *Clock) Tick() {
	now := time.Now()
	c.frames += 1
	c.delta = now.Sub(c.frame).Seconds()
	c.elapsed += c.delta
	c.frame = now

	if c.elapsed >= 1 {
		c.fps = float64(c.frames)
		c.elapsed = math.Mod(c.elapsed, 1)
		c.frames = 0
	}
}

func (c *Clock) Delta() float32 {
	return float32(c.delta)
}

func (c *Clock) Fps() float32 {
	return float32(c.fps)
}

func (c *Clock) Time() float32 {
	return float32(time.Now().Sub(c.start).Seconds())
}
