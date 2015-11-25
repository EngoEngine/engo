package engi

import (
	"testing"
)

const (
	benchmarkComponentCount = 1000
)

type getComponentSystem struct {
	*System
}

func (getComponentSystem) Type() string {
	return "getComponentSystem"
}

func (g *getComponentSystem) New() {
	g.System = NewSystem()
}

func (g *getComponentSystem) Update(entity *Entity, dt float32) {
	var sp *SpaceComponent
	if !entity.Component(&sp) {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if sp == nil {
		return
	}

	if len(entity.components) != 2 {
		return
	}

	var ren *RenderComponent
	if !entity.Component(&ren) {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if ren == nil {
		return
	}
}

func BenchmarkComponent(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystem{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity([]string{"getComponentSystem"})
			e.AddComponent(&SpaceComponent{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

func BenchmarkComponentDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystem{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity([]string{"getComponentSystem"})
			e.AddComponent(&SpaceComponent{})
			e.AddComponent(&RenderComponent{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

type getComponentSystemFast struct {
	*System
}

func (getComponentSystemFast) Type() string {
	return "getComponentSystemFast"
}

func (g *getComponentSystemFast) New() {
	g.System = NewSystem()
}

func (g *getComponentSystemFast) Update(entity *Entity, dt float32) {
	var sp *SpaceComponent
	var ok bool
	if sp, ok = entity.ComponentFast(sp).(*SpaceComponent); !ok {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if sp == nil {
		return
	}

	if len(entity.components) != 2 {
		return
	}

	var ren *RenderComponent
	if ren, ok = entity.ComponentFast(ren).(*RenderComponent); !ok {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if ren == nil {
		return
	}
}

func BenchmarkComponentFast(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystemFast{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity([]string{"getComponentSystemFast"})
			e.AddComponent(&SpaceComponent{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

func BenchmarkComponentFastDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystemFast{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity([]string{"getComponentSystemFast"})
			e.AddComponent(&SpaceComponent{})
			e.AddComponent(&RenderComponent{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}
