//+build demo

package main

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Guy struct {
	ecs.BasicEntity
	common.MouseComponent
	common.RenderComponent
	common.SpaceComponent
	dragComponent
}

type DefaultScene struct{}

func (*DefaultScene) Type() string {
	return "Default Scene"
}

func (*DefaultScene) Preload() {
	engo.Files.Load("icon.png")
}

func (*DefaultScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&DragSystem{})

	guyTexture, _ := common.LoadedSprite("icon.png")
	guy := Guy{BasicEntity: ecs.NewBasic()}
	guy.RenderComponent = common.RenderComponent{
		Drawable: guyTexture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    guyTexture.Width() * guy.RenderComponent.Scale.X,
		Height:   guyTexture.Height() * guy.RenderComponent.Scale.Y,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&guy.BasicEntity, &guy.MouseComponent, &guy.SpaceComponent, &guy.RenderComponent)
		case *DragSystem:
			sys.Add(&guy.BasicEntity, &guy.SpaceComponent, &guy.MouseComponent, &guy.dragComponent)
		}
	}
}

type dragComponent struct {
	following  bool
	xoff, yoff float32
}

type dragEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.MouseComponent
	*dragComponent
}

type DragSystem struct {
	entities []dragEntity
}

func (d *DragSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, mouse *common.MouseComponent, drag *dragComponent) {
	d.entities = append(d.entities, dragEntity{basic, space, mouse, drag})
}

func (d *DragSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range d.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		d.entities = append(d.entities[:delete], d.entities[delete+1:]...)
	}
}

func (d *DragSystem) Update(dt float32) {
	for _, e := range d.entities {
		if e.MouseComponent.Clicked {
			e.following = true
			e.xoff = engo.Input.Mouse.X - e.SpaceComponent.Position.X
			e.yoff = engo.Input.Mouse.Y - e.SpaceComponent.Position.Y
		}
		if e.MouseComponent.Released {
			e.following = false
			e.xoff = 0
			e.yoff = 0
		}
		if e.following {
			e.SpaceComponent.Position.Set(engo.Input.Mouse.X-e.xoff, engo.Input.Mouse.Y-e.yoff)
		}
	}
}

func main() {
	engo.Run(engo.RunOptions{
		Title:  "Mouse Drag Demo",
		Width:  1024,
		Height: 640,
	}, &DefaultScene{})
}
