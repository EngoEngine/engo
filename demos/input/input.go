package main

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

func (*DefaultScene) Preload() {}
func (*DefaultScene) Setup(w *ecs.World) {
	w.AddSystem(&common.RenderSystem{})

	// Register the input system responding
	// to the input actions registered below.
	w.AddSystem(&InputSystem{})

	// Register a button that is triggered when
	// the state of 'Space' or 'Enter' is changed.
	engo.Input.RegisterButton("action", engo.Space, engo.Enter)

	// Register an axis where 'A' will return a
	// negative value and 'D' returns a positive value.
	engo.Input.RegisterAxis("sideways", engo.AxisKeyPair{engo.A, engo.D})
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
	// Look up the button once
	btn := engo.Input.Button("action")

	// And act on the current state
	if btn.Down() {
		fmt.Println("Still down!")
	} else if btn.JustPressed() {
		fmt.Println("Key just pressed!")
	} else if btn.JustReleased() {
		fmt.Println("Key just released!")
	}

	// Check the axis value and act as required
	if v := engo.Input.Axis("sideways").Value(); v != 0 {
		fmt.Printf("Axis value: %.2f \n", v)
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
