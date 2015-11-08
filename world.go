package engi

type World struct {
	entities map[string]*Entity
	systems  []Systemer

	defaultBatch *Batch
	hudBatch     *Batch

	isSetup bool
	paused  bool
}

func (w *World) new() {
	if !w.isSetup {
		w.entities = make(map[string]*Entity)
		if !headless {
			w.defaultBatch = NewBatch(Width(), Height(), batchVert, batchFrag)
			w.hudBatch = NewBatch(Width(), Height(), hudVert, hudFrag)
		}

		// Default WorldBounds values
		WorldBounds.Max = Point{Width(), Height()}

		// Initialize cameraSystem
		cam = &cameraSystem{}
		w.AddSystem(cam)

		w.isSetup = true
	}
}

func (w *World) AddEntity(entity *Entity) {
	w.entities[entity.ID()] = entity

	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.AddEntity(entity)
		}
	}
}

func (w *World) RemoveEntity(entity *Entity) {
	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.RemoveEntity(entity)
		}
	}
	delete(w.entities, entity.ID())
}

func (w *World) AddSystem(system Systemer) {
	system.New()
	system.SetWorld(w)
	w.systems = append(w.systems, system)
}

func (w *World) Entities() []*Entity {
	entities := make([]*Entity, len(w.entities))
	for _, v := range w.entities {
		entities = append(entities, v)
	}

	return entities
}

func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) pre() {
	if !headless {
		Gl.Clear(Gl.COLOR_BUFFER_BIT)
	}
}

func (w *World) post() {}

func (w *World) update(dt float32) {
	w.pre()

	var unp *UnpauseComponent

	for _, system := range w.Systems() {
		if headless && system.SkipOnHeadless() {
			continue // so skip it
		}

		system.Pre()
		for _, entity := range system.Entities() {
			if w.paused {
				ok := entity.GetComponent(&unp)
				if !ok {
					continue // so skip it
				}
			}
			system.Update(entity, dt)
		}
		system.Post()
	}

	w.post()
}

func (w *World) batch(prio PriorityLevel) *Batch {
	if prio >= HUDGround {
		return w.hudBatch
	} else {
		return w.defaultBatch
	}
}
