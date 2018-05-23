//+build demo

package main

import (
	"image/color"
	"log"
	"strconv"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var fnt *common.Font

type ClickComponent struct {
	label string
}

type MyLabel struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	ClickComponent
}

type DefaultScene struct{}

func (*DefaultScene) Type() string { return "Default Scene" }

func (*DefaultScene) Preload() {
	err := engo.Files.Load("Roboto-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
}

func (*DefaultScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&ClickSystem{})

	//text
	fnt = &common.Font{
		URL:  "Roboto-Regular.ttf",
		FG:   color.Black,
		Size: 64,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	label1 := MyLabel{BasicEntity: ecs.NewBasic()}
	label1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "LMB: (0, 0)",
	}
	label1.SetShader(common.HUDShader)
	label1.ClickComponent.label = "left click"

	label2 := MyLabel{BasicEntity: ecs.NewBasic()}
	label2.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "RMB: (0, 0)",
	}
	label2.SpaceComponent.Position = engo.Point{0, 150}
	label2.SetShader(common.HUDShader)
	label2.ClickComponent.label = "right click"

	// Add our text to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
			sys.Add(&label2.BasicEntity, &label2.RenderComponent, &label2.SpaceComponent)
		case *ClickSystem:
			sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent, &label1.ClickComponent)
			sys.Add(&label2.BasicEntity, &label2.RenderComponent, &label2.SpaceComponent, &label2.ClickComponent)
		}
	}
}

type mouseState uint

const (
	up mouseState = iota
	down
	justPressed
)

type clickEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
	*ClickComponent
}

type ClickSystem struct {
	entities []clickEntity

	left, right mouseState
}

func (c *ClickSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent, click *ClickComponent) {
	c.entities = append(c.entities, clickEntity{basic, render, space, click})
}

func (c *ClickSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}
func (c *ClickSystem) Update(dt float32) {
	//setup mouse state
	if c.left == justPressed {
		c.left = down
	}
	if c.right == justPressed {
		c.right = down
	}
	if engo.Input.Mouse.Action == engo.Press {
		if engo.Input.Mouse.Button == engo.MouseButtonLeft {
			if c.left == up {
				c.left = justPressed
			}
		} else if engo.Input.Mouse.Button == engo.MouseButtonRight {
			if c.right == up {
				c.right = justPressed
			}
		}
	} else if engo.Input.Mouse.Action == engo.Release {
		if engo.Input.Mouse.Button == engo.MouseButtonLeft {
			c.left = up
		} else if engo.Input.Mouse.Button == engo.MouseButtonRight {
			c.right = up
		}
	}

	//loop through entities
	for _, e := range c.entities {
		switch e.ClickComponent.label {
		case "left click":
			if c.left == justPressed {
				txt := "LMB: (" + strconv.FormatFloat(float64(engo.Input.Mouse.X), 'f', 1, 32) + ", " + strconv.FormatFloat(float64(engo.Input.Mouse.Y), 'f', 1, 32) + ")"
				e.RenderComponent.Drawable.Close()
				e.RenderComponent.Drawable = common.Text{
					Font: fnt,
					Text: txt,
				}
			}
		case "right click":
			if c.right == justPressed {
				txt := "RMB: (" + strconv.FormatFloat(float64(engo.Input.Mouse.X), 'f', 1, 32) + ", " + strconv.FormatFloat(float64(engo.Input.Mouse.Y), 'f', 1, 32) + ")"
				e.RenderComponent.Drawable.Close()
				e.RenderComponent.Drawable = common.Text{
					Font: fnt,
					Text: txt,
				}
			}
		}
	}
}

func main() {
	engo.Run(engo.RunOptions{
		Title:  "Mouse Click Demo",
		Width:  1024,
		Height: 640,
	}, &DefaultScene{})
}
