//+build demo

package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type MyShape struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type ClickableShape struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent
	OnClickComponent
}

type OnClickComponent struct {
	On bool
}

type OnClickEntity struct {
	*ecs.BasicEntity
	*common.MouseComponent
	*common.RenderComponent
	*OnClickComponent
}

func (e *OnClickEntity) Light(on bool) {
	e.On = on
	if on {
		e.RenderComponent.Color = color.RGBA{0, 255, 0, 255}
	} else {
		e.RenderComponent.Color = color.RGBA{255, 0, 0, 255}
	}
}

type OnClickSystem struct {
	entities map[uint64]OnClickEntity
}

func (s *OnClickSystem) New(w *ecs.World) {
	s.entities = make(map[uint64]OnClickEntity)
}

func (s *OnClickSystem) Add(basic *ecs.BasicEntity, mouse *common.MouseComponent, render *common.RenderComponent, onClick *OnClickComponent) {
	s.entities[basic.ID()] = OnClickEntity{basic, mouse, render, onClick}
}

func (s *OnClickSystem) Remove(basic ecs.BasicEntity) {
	delete(s.entities, basic.ID())
}

func (s *OnClickSystem) Update(dt float32) {
	for _, ent := range s.entities {
		if ent.Clicked {
			if s.CanLight(ent) {
				ent.Light(!ent.On)
				if !ent.On {
					for _, child := range ent.Descendents() {
						if e, ok := s.entities[child.ID()]; ok {
							e.Light(false)
						}
					}
				}
			}
		}
	}
}

func (s *OnClickSystem) CanLight(ent OnClickEntity) bool {
	if ent.Parent() == nil {
		return true
	}
	if parent, ok := s.entities[ent.Parent().ID()]; ok {
		return parent.On
	}
	return false
}

type DefaultScene struct{}

func (*DefaultScene) Type() string { return "DefaultScene" }

func (*DefaultScene) Preload() {}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.RGBA{55, 55, 55, 255})
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&OnClickSystem{})

	rects := []engo.Point{
		engo.Point{
			X: 225,
			Y: 15,
		},
		engo.Point{
			X: 25,
			Y: 135,
		},
		engo.Point{
			X: 225,
			Y: 135,
		},
		engo.Point{
			X: 425,
			Y: 135,
		},
		engo.Point{
			X: 25,
			Y: 255,
		},
		engo.Point{
			X: 120,
			Y: 255,
		},
		engo.Point{
			X: 225,
			Y: 255,
		},
		engo.Point{
			X: 425,
			Y: 255,
		},
	}
	rectBasics := make([]*ecs.BasicEntity, len(rects))

	for i, pt := range rects {
		rect := ClickableShape{BasicEntity: ecs.NewBasic()}
		rectBasics[i] = &rect.BasicEntity
		rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{pt.X, pt.Y}, Width: 50, Height: 50}
		rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{0, 255, 0, 255}}
		rect.On = true

		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&rect.BasicEntity, &rect.RenderComponent, &rect.SpaceComponent)
			case *common.MouseSystem:
				sys.Add(&rect.BasicEntity, &rect.MouseComponent, &rect.SpaceComponent, &rect.RenderComponent)
			case *OnClickSystem:
				sys.Add(&rect.BasicEntity, &rect.MouseComponent, &rect.RenderComponent, &rect.OnClickComponent)
			}
		}
	}

	lines := []engo.Line{
		engo.Line{
			P1: engo.Point{X: 250, Y: 65},
			P2: engo.Point{X: 250, Y: 100},
		},
		engo.Line{
			P1: engo.Point{X: 50, Y: 100},
			P2: engo.Point{X: 450, Y: 100},
		},
		engo.Line{
			P1: engo.Point{X: 50, Y: 100},
			P2: engo.Point{X: 50, Y: 135},
		},
		engo.Line{
			P1: engo.Point{X: 250, Y: 100},
			P2: engo.Point{X: 250, Y: 135},
		},
		engo.Line{
			P1: engo.Point{X: 448, Y: 100},
			P2: engo.Point{X: 448, Y: 135},
		},
		engo.Line{
			P1: engo.Point{X: 50, Y: 185},
			P2: engo.Point{X: 50, Y: 220},
		},
		engo.Line{
			P1: engo.Point{X: 50, Y: 220},
			P2: engo.Point{X: 150, Y: 220},
		},
		engo.Line{
			P1: engo.Point{X: 50, Y: 220},
			P2: engo.Point{X: 50, Y: 255},
		},
		engo.Line{
			P1: engo.Point{X: 148, Y: 220},
			P2: engo.Point{X: 148, Y: 255},
		},
		engo.Line{
			P1: engo.Point{X: 250, Y: 185},
			P2: engo.Point{X: 250, Y: 255},
		},
		engo.Line{
			P1: engo.Point{X: 450, Y: 185},
			P2: engo.Point{X: 450, Y: 255},
		},
	}
	for i := 0; i < len(lines); i++ {
		line := MyShape{BasicEntity: ecs.NewBasic()}
		line.SpaceComponent = common.SpaceComponent{Position: engo.Point{lines[i].P1.X, lines[i].P1.Y}, Width: 2, Height: lines[i].Magnitude(), Rotation: lines[i].AngleDeg() + 180}
		line.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{0, 0, 0, 255}}

		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&line.BasicEntity, &line.RenderComponent, &line.SpaceComponent)
			}
		}
	}

	// Setup parent-child relationships
	rectBasics[0].AppendChild(rectBasics[1])
	rectBasics[0].AppendChild(rectBasics[2])
	rectBasics[0].AppendChild(rectBasics[3])

	rectBasics[1].AppendChild(rectBasics[4])
	rectBasics[1].AppendChild(rectBasics[5])

	rectBasics[2].AppendChild(rectBasics[6])

	rectBasics[3].AppendChild(rectBasics[7])
}

func main() {
	opts := engo.RunOptions{
		Title:  "Parenting",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &DefaultScene{})
}
