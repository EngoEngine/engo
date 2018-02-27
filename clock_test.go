package engo

import (
	"testing"
	"time"
)

func TestClockTick(t *testing.T) {
	clock := NewClock()
	before := clock.counter
	clock.Tick()
	after := clock.counter
	if after != before+1 {
		t.Errorf("Tick did not increase the counter, before: %v, after: %v", before, after)
	}
	<-time.After(1 * time.Second)
	clock.Tick()
	if clock.counter != 0 {
		t.Error("Clock's count did not reset to 0 after waiting over 1 second")
	}
}

func TestClockFPS(t *testing.T) {
	clock := NewClock()
	for i := 0; i < 6; i++ {
		<-time.After(time.Second / 6)
		clock.Tick()
	}
	if !FloatEqual(clock.FPS(), float32(6)) {
		t.Errorf("Clock's FPS did not match 6 fps, was %v", clock.FPS())
	}
}

func TestClockTime(t *testing.T) {
	clock := NewClock()
	<-time.After(1 * time.Second)
	if !FloatEqual(clock.Time(), float32(1)) {
		t.Errorf("Clock's duration from Time did not match 1 second, was %v", clock.Time())
	}
}
