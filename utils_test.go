package engo

import (
	"testing"

	"engo.io/ecs"
)

type inlineGame struct {
	preloadFunc func()
	setupFunc   func(*ecs.World)
}

func (m *inlineGame) Preload() {
	m.preloadFunc()
}

func (m *inlineGame) Setup(w *ecs.World) {
	m.setupFunc(w)
}

// Bench is a helper-function to easily benchmark one frame, given a preload / setup function
func Bench(b *testing.B, preload func(), setup func(w *ecs.World)) {
	g := &inlineGame{preloadFunc: preload, setupFunc: setup}

	w := &ecs.World{}
	w.New()

	g.Preload()
	g.Setup(w)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w.Update(1 / 120) // 120 fps
	}
}
