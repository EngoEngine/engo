// Copyright 2014 Harrison Shoebridge. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

const (
	// HighestGround is the highest PriorityLevel that will be rendered
	HighestGround PriorityLevel = 50
	// HUDGround is a PriorityLevel from which everything isn't being affected by the Camera
	HUDGround    PriorityLevel = 40
	Foreground   PriorityLevel = 30
	MiddleGround PriorityLevel = 20
	ScenicGround PriorityLevel = 10
	// Background is the lowest PriorityLevel that will be rendered
	Background PriorityLevel = 0
	// Hidden indicates that it should not be rendered by the RenderSystem
	Hidden PriorityLevel = -1
)

type PriorityLevel int

type Renderable interface {
	Render(b *Batch, render *RenderComponent, space *SpaceComponent)
}

type RenderComponent struct {
	Display      Renderable
	Scale        Point
	Label        string
	priority     PriorityLevel
	Transparency float32
	Color        uint32
}

type renderChangeMessage struct {
	entity      *Entity
	oldPriority PriorityLevel
	newPriority PriorityLevel
}

func (renderChangeMessage) Type() string {
	return "renderChangeMessage"
}

func NewRenderComponent(display Renderable, scale Point, label string) *RenderComponent {
	return &RenderComponent{
		Display:      display,
		Scale:        scale,
		Label:        label,
		priority:     MiddleGround,
		Transparency: 1,
		Color:        0xffffff,
	}
}

func (r *RenderComponent) SetPriority(p PriorityLevel) {
	Mailbox.Dispatch(renderChangeMessage{})
	r.priority = p
}

func (*RenderComponent) Type() string {
	return "RenderComponent"
}

type RenderSystem struct {
	renders map[PriorityLevel][]*Entity
	changed bool
	*System
}

func (rs *RenderSystem) New() {
	rs.renders = make(map[PriorityLevel][]*Entity)
	rs.System = NewSystem()
	rs.ShouldSkipOnHeadless = true

	Mailbox.Listen("renderChangeMessage", func(m Message) {
		rs.changed = true
	})
}

func (rs *RenderSystem) AddEntity(e *Entity) {
	rs.changed = true
	rs.System.AddEntity(e)
}

func (rs *RenderSystem) RemoveEntity(e *Entity) {
	rs.changed = true
	rs.System.RemoveEntity(e)
}

func (rs *RenderSystem) Pre() {
	if !rs.changed {
		return
	}

	rs.renders = make(map[PriorityLevel][]*Entity)
}

func (rs *RenderSystem) Post() {
	var currentBatch *Batch

	for i := Background; i <= HighestGround; i++ {
		if len(rs.renders[i]) == 0 {
			continue
		}

		// Retrieve a batch, may be the default one -- then call .Begin() if we arent already using it
		batch := world.batch(i)
		if batch != currentBatch {
			if currentBatch != nil {
				currentBatch.End()
			}
			batch.Begin()
			currentBatch = batch
		}
		// Then render everything for this level
		for _, entity := range rs.renders[i] {
			var (
				render *RenderComponent
				space  *SpaceComponent
				ok     bool
			)

			if render, ok = entity.ComponentFast(render).(*RenderComponent); !ok {
				continue // with other entities
			}

			if space, ok = entity.ComponentFast(space).(*SpaceComponent); !ok {
				continue // with other entities
			}

			render.Display.Render(batch, render, space)
		}
	}

	if currentBatch != nil {
		currentBatch.End()
	}

	rs.changed = false
}

func (rs *RenderSystem) Update(entity *Entity, dt float32) {
	if !rs.changed {
		return
	}

	var render *RenderComponent
	var ok bool

	if render, ok = entity.ComponentFast(render).(*RenderComponent); !ok {
		return
	}

	rs.renders[render.priority] = append(rs.renders[render.priority], entity)
}

func (*RenderSystem) Type() string {
	return "RenderSystem"
}

func (rs *RenderSystem) Priority() int {
	return 1
}
