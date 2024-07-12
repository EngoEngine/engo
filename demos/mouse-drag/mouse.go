//go:build demo
// +build demo

package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
	engo.Files.Load("guy.png")
	engo.Files.Load("banana.png")
	engo.Files.Load("red-cherry.png")
	engo.Files.Load("watermelon.png")
}

func (*DefaultScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&DragSystem{})

	guyTexture, _ := common.LoadedSprite("guy.png")
	guy := Guy{BasicEntity: ecs.NewBasic()}
	guy.RenderComponent = common.RenderComponent{
		Drawable: guyTexture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    guyTexture.Width() * guy.RenderComponent.Scale.X,
		Height:   guyTexture.Height() * guy.RenderComponent.Scale.Y,
		Rotation: 15,
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

	bananaTexture, _ := common.LoadedSprite("banana.png")
	banana := Guy{BasicEntity: ecs.NewBasic()}
	banana.RenderComponent = common.RenderComponent{
		Drawable: bananaTexture,
	}
	banana.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{200, 0},
		Width:    bananaTexture.Width() * banana.RenderComponent.Scale.X,
		Height:   bananaTexture.Height() * banana.RenderComponent.Scale.Y,
		Rotation: 30,
	}
	bananaOutline := []engo.Line{}
	points := []float32{0, 95, 11, 93, 47, 89, 60.5, 82.3, 76.7, 66.2, 81.9, 57.0, 85.3, 31.5,
		78.9, 14.6, 78.6, 6.4, 72.2, 0, 88.1, 0, 93.0, 4.1, 98.8, 20.1, 110.2,
		42.6, 110.2, 60.4, 97.9, 91.7, 84.8, 105.4, 71.4, 113.6, 50.1, 119.8,
		28.7, 119.8, 28.7, 119.8, 28.7, 119.8, 12, 112.5, 0, 103.7, 0, 95.0}
	for i := 0; i < len(points)-2; i += 2 {
		line := engo.Line{
			P1: engo.Point{
				X: points[i],
				Y: points[i+1],
			},
			P2: engo.Point{
				X: points[i+2],
				Y: points[i+3],
			},
		}
		bananaOutline = append(bananaOutline, line)
	}
	banana.SpaceComponent.AddShape(common.Shape{Lines: bananaOutline})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&banana.BasicEntity, &banana.RenderComponent, &banana.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&banana.BasicEntity, &banana.MouseComponent, &banana.SpaceComponent, &banana.RenderComponent)
		case *DragSystem:
			sys.Add(&banana.BasicEntity, &banana.SpaceComponent, &banana.MouseComponent, &banana.dragComponent)
		}
	}

	watermelonTexture, _ := common.LoadedSprite("watermelon.png")
	watermelon := Guy{BasicEntity: ecs.NewBasic()}
	watermelon.RenderComponent = common.RenderComponent{
		Drawable: watermelonTexture,
	}
	watermelon.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{400, 0},
		Width:    watermelonTexture.Width() * watermelon.RenderComponent.Scale.X,
		Height:   watermelonTexture.Height() * watermelon.RenderComponent.Scale.Y,
		Rotation: 45,
	}
	watermelon.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 61, Cy: 50, Rx: 61, Ry: 50}})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&watermelon.BasicEntity, &watermelon.RenderComponent, &watermelon.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&watermelon.BasicEntity, &watermelon.MouseComponent, &watermelon.SpaceComponent, &watermelon.RenderComponent)
		case *DragSystem:
			sys.Add(&watermelon.BasicEntity, &watermelon.SpaceComponent, &watermelon.MouseComponent, &watermelon.dragComponent)
		}
	}

	cherryTexture, _ := common.LoadedSprite("red-cherry.png")
	cherry := Guy{BasicEntity: ecs.NewBasic()}
	cherry.RenderComponent = common.RenderComponent{
		Drawable: cherryTexture,
	}
	cherry.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{600, 0},
		Width:    cherryTexture.Width() * cherry.RenderComponent.Scale.X,
		Height:   cherryTexture.Height() * cherry.RenderComponent.Scale.Y,
		Rotation: 60,
	}
	cherry.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 25, Cy: 57.5, Rx: 25, Ry: 25.5}})
	cherry.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 59, Cy: 75, Rx: 26, Ry: 25}})
	cherryOutline := []engo.Line{}
	points = []float32{36.2, 37.5, 46.4, 22.7, 29.9, 18.2, 28.5, 16.4, 44.1, 3.5,
		50, 0, 61.2, 0, 69.7, 5.0, 83.1, 5.1, 90.2, 9.0, 98.8, 22.6, 106.0, 45.9,
		106, 49.4, 100.7, 49.3, 79.8, 39.1, 70, 28.9, 69, 24.7, 65.7, 37.8, 65.7,
		54.2, 59.5, 54.1, 59.2, 35.2, 62.9, 22.7, 54.2, 24.0, 40.5, 42.0, 36.2, 37.5}
	for i := 0; i < len(points)-2; i += 2 {
		line := engo.Line{
			P1: engo.Point{
				X: points[i],
				Y: points[i+1],
			},
			P2: engo.Point{
				X: points[i+2],
				Y: points[i+3],
			},
		}
		cherryOutline = append(cherryOutline, line)
	}
	cherry.SpaceComponent.AddShape(common.Shape{Lines: cherryOutline})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&cherry.BasicEntity, &cherry.RenderComponent, &cherry.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&cherry.BasicEntity, &cherry.MouseComponent, &cherry.SpaceComponent, &cherry.RenderComponent)
		case *DragSystem:
			sys.Add(&cherry.BasicEntity, &cherry.SpaceComponent, &cherry.MouseComponent, &cherry.dragComponent)
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
