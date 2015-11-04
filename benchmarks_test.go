package engi_test

import (
	"github.com/paked/engi"
	"testing"
)

// BenchmarkEmpty creates the game, and measures the runtime of a single frame, w/o anything set up
func BenchmarkEmpty(b *testing.B) {
	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {}
	engi.Bench(b, preload, setup)
}

// BenchmarkSystem10 creates 10 `NilSystem`s
func BenchmarkSystem10(b *testing.B) {
	const count = 10

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		for i := 0; i < count; i++ {
			w.AddSystem(&engi.NilSystem{})
		}
	}
	engi.Bench(b, preload, setup)
}

// BenchmarkSystem1000 creates 1000 `NilSystem`s
func BenchmarkSystem1000(b *testing.B) {
	const count = 1000

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		for i := 0; i < count; i++ {
			w.AddSystem(&engi.NilSystem{})
		}
	}
	engi.Bench(b, preload, setup)
}

// BenchmarkEntity10 creates 10 `Entity`s which all depend on the `NilSystem`
func BenchmarkEntity10(b *testing.B) {
	const count = 10

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		w.AddSystem(&engi.NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(engi.NewEntity([]string{"NilSystem"}))
		}
	}
	engi.Bench(b, preload, setup)
}

// BenchmarkEntity1000 creates 1000 `Entity`s which all depend on the `NilSystem`
func BenchmarkEntity1000(b *testing.B) {
	const count = 1000

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		w.AddSystem(&engi.NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(engi.NewEntity([]string{"NilSystem"}))
		}
	}
	engi.Bench(b, preload, setup)
}
