package engo

import (
	"testing"
	"time"
)

// testTime is the time interface where testTime.Now() is controllable using
// testTime.curTime
type testTime struct {
	curTime int64
}

func (t testTime) Now() int64 {
	return t.curTime
}

func TestClockTick(t *testing.T) {
	theTimer = testTime{0}
	clock := NewClock()
	before := clock.counter
	clock.Tick()
	after := clock.counter
	if after != before+1 {
		t.Errorf("Tick did not increase the counter, before: %v, after: %v", before, after)
	}
	theTimer = testTime{1000000000}
	clock.Tick()
	if clock.counter != 0 {
		t.Error("Clock's count did not reset to 0 after waiting over 1 second")
	}
}

func TestClockFPS(t *testing.T) {
	data := []struct {
		fps int
	}{
		{5},
		{10},
		{20},
		{30},
		{60},
	}
	for _, d := range data {
		theTimer = testTime{0}
		clock := NewClock()
		tickTime := 1000000000 / d.fps
		curTime := int64(0)
		for i := 0; i < d.fps; i++ {
			curTime += int64(tickTime) + 1
			theTimer = testTime{curTime}
			clock.Tick()
		}
		if clock.FPS() != float32(d.fps) {
			t.Errorf("Clock's FPS did not match %v fps, was %v", d.fps, clock.FPS())
		}
	}
}

func TestClockDelta(t *testing.T) {
	data := []struct {
		delta int64
	}{
		{6000000000},
		{16666667},
		{33333333},
		{66666666},
		{60},
	}
	for _, d := range data {
		exp := float32(d.delta) / 1000000000.0
		theTimer = testTime{0}
		clock := NewClock()
		theTimer = testTime{d.delta}
		clock.Tick()
		if clock.Delta() != exp {
			t.Errorf("Clock's Delta did not match %v, was %v", exp, clock.Delta())
		}
	}
}

func TestClockTime(t *testing.T) {
	data := []struct {
		time int64
	}{
		{6000000000},
		{16666667},
		{33333333},
		{66666666},
		{60},
	}
	for _, d := range data {
		exp := float32(d.time) / 1000000000.0
		theTimer = testTime{0}
		clock := NewClock()
		theTimer = testTime{d.time}
		if clock.Time() != exp {
			t.Errorf("Clock's duration from Time() did not match %v seconds, was %v", exp, clock.Time())
		}
	}
}

func TestTheTimerNow(t *testing.T) {
	theTimer = realTime{}
	res := time.Now().UnixNano() - theTimer.Now()
	if res < -1000000 { //this is an arbitrary choice, could be smaller or larger  depending on machine testing on
		t.Error("theTimer when it's a realTime did not produce time.Now().UnixNano()")
	}
}
