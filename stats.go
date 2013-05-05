package eng

import (
	"log"
	"math"
	"time"
)

type stats struct {
	Elapsed, Dt, Fps, Frames, Period float64
	Then                             time.Time
	show                             bool
}

func NewStats(show bool) *stats {
	st := new(stats)
	st.Period = 1
	st.Update()
	st.show = show
	return st
}

func (t *stats) Update() {
	now := time.Now()
	t.Frames += 1
	t.Dt = now.Sub(t.Then).Seconds()
	t.Elapsed += t.Dt
	t.Then = now

	if t.Elapsed >= t.Period {
		t.Fps = t.Frames / t.Period
		t.Elapsed = math.Mod(t.Elapsed, t.Period)
		t.Frames = 0
		if t.show {
			log.Println(t.Fps)
		}
	}
}
