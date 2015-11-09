package engi

import (
	"testing"
)

type NilSystem struct {
	*System
}

func (ns *NilSystem) New() {
	ns.System = NewSystem()
}

func (*NilSystem) Update(*Entity, float32) {}

func (*NilSystem) Type() string {
	return "NilSystem"
}

// BenchmarkEmpty creates the game, and measures the runtime of a single frame, w/o anything set up
func BenchmarkEmpty(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {}
	Bench(b, preload, setup)
}

// BenchmarkSystem10 creates 10 `NilSystem`s
func BenchmarkSystem10(b *testing.B) {
	const count = 10

	preload := func() {}
	setup := func(w *World) {
		for i := 0; i < count; i++ {
			w.AddSystem(&NilSystem{})
		}
	}
	Bench(b, preload, setup)
}

// BenchmarkSystem1000 creates 1000 `NilSystem`s
func BenchmarkSystem1000(b *testing.B) {
	const count = 1000

	preload := func() {}
	setup := func(w *World) {
		for i := 0; i < count; i++ {
			w.AddSystem(&NilSystem{})
		}
	}
	Bench(b, preload, setup)
}

// BenchmarkEntity10 creates 10 `Entity`s which all depend on the `NilSystem`
func BenchmarkEntity10(b *testing.B) {
	const count = 10

	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(NewEntity([]string{"NilSystem"}))
		}
	}
	Bench(b, preload, setup)
}

// BenchmarkEntity1000 creates 1000 `Entity`s which all depend on the `NilSystem`
func BenchmarkEntity1000(b *testing.B) {
	const count = 1000

	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(NewEntity([]string{"NilSystem"}))
		}
	}
	Bench(b, preload, setup)
}
