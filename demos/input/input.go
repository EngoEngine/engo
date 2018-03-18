//+build demo

package main

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

func (*DefaultScene) Preload() {}
func (*DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&InputSystem{})

	engo.Input.RegisterAxis("sideways", engo.AxisKeyPair{engo.KeyA, engo.KeyD})
	engo.Input.RegisterButton("action", engo.KeySpace, engo.KeyEnter)
}

func (*DefaultScene) Type() string { return "Game" }

type inputEntity struct {
	*ecs.BasicEntity
}

type InputSystem struct {
	entities []inputEntity
}

func (c *InputSystem) Add(basic *ecs.BasicEntity) {
	c.entities = append(c.entities, inputEntity{basic})
}

func (c *InputSystem) Remove(basic ecs.BasicEntity) {
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

func (c *InputSystem) Update(dt float32) {
	if v := engo.Input.Axis("sideways").Value(); v != 0 {
		fmt.Println(v)
	}

	if btn := engo.Input.Button("action"); btn.JustPressed() {
		fmt.Println("DOWN!")
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Input Demo",
		Width:  1024,
		Height: 640,
	}

	engo.Run(opts, &DefaultScene{})
}
