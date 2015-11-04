package engi_test

import (
	"github.com/paked/engi"
	"testing"
)

// BenchmarkCollisionSystem10 creates 10 entities, of which half are solid, and all are Main
func BenchmarkCollisionSystem10(b *testing.B) {
	const count = 10

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		w.AddSystem(&engi.CollisionSystem{})
		for i := 0; i < count; i++ {
			ent := engi.NewEntity([]string{"CollisionSystem"})
			ent.AddComponent(&engi.SpaceComponent{engi.Point{0, 0}, 10, 10})
			ent.AddComponent(&engi.CollisionComponent{Solid: i%2 == 0, Main: true})
			w.AddEntity(ent)
		}
	}
	engi.Bench(b, preload, setup)
}

// BenchmarkCollisionSystem100 creates 100 entities, of which half are solid, and all are Main
func BenchmarkCollisionSystem100(b *testing.B) {
	const count = 100

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		w.AddSystem(&engi.CollisionSystem{})
		for i := 0; i < count; i++ {
			ent := engi.NewEntity([]string{"CollisionSystem"})
			ent.AddComponent(&engi.SpaceComponent{engi.Point{0, 0}, 10, 10})
			ent.AddComponent(&engi.CollisionComponent{Solid: i%2 == 0, Main: true})
			w.AddEntity(ent)
		}
	}
	engi.Bench(b, preload, setup)
}

// BenchmarkCollisionSystem1000 creates 1000 entities, of which half are solid, and all are Main
func BenchmarkCollisionSystem1000(b *testing.B) {
	const count = 1000

	preload := func(w *engi.World) {}
	setup := func(w *engi.World) {
		w.AddSystem(&engi.CollisionSystem{})
		for i := 0; i < count; i++ {
			ent := engi.NewEntity([]string{"CollisionSystem"})
			ent.AddComponent(&engi.SpaceComponent{engi.Point{0, 0}, 10, 10})
			ent.AddComponent(&engi.CollisionComponent{Solid: i%2 == 0, Main: true})
			w.AddEntity(ent)
		}
	}
	engi.Bench(b, preload, setup)
}
