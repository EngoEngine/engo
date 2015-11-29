package engi

import (
	"github.com/paked/engi/ecs"
	"testing"
)

type NilSystem struct {
	*ecs.System
}

func (ns *NilSystem) New() {
	ns.System = ecs.NewSystem()
}

func (*NilSystem) Update(*ecs.Entity, float32) {}

func (*NilSystem) Type() string {
	return "NilSystem"
}

// BenchmarkEmpty creates the game, and measures the runtime of a single frame, w/o anything set up
func BenchmarkEmpty(b *testing.B) {
	preload := func() {}
	setup := func(w *ecs.World) {}
	Bench(b, preload, setup)
}

// BenchmarkSystem10 creates 10 `NilSystem`s
func BenchmarkSystem10(b *testing.B) {
	const count = 10

	preload := func() {}
	setup := func(w *ecs.World) {
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
	setup := func(w *ecs.World) {
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
	setup := func(w *ecs.World) {
		w.AddSystem(&NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(ecs.NewEntity([]string{"NilSystem"}))
		}
	}
	Bench(b, preload, setup)
}

// BenchmarkEntity1000 creates 1000 `Entity`s which all depend on the `NilSystem`
func BenchmarkEntity1000(b *testing.B) {
	const count = 1000

	preload := func() {}
	setup := func(w *ecs.World) {
		w.AddSystem(&NilSystem{})
		for i := 0; i < count; i++ {
			w.AddEntity(ecs.NewEntity([]string{"NilSystem"}))
		}
	}
	Bench(b, preload, setup)
}
