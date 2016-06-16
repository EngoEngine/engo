package main

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/act"
	"engo.io/engo/common"
)

type DefaultScene struct{}

func (*DefaultScene) Preload() {}

func (*DefaultScene) Setup(w *ecs.World) {
	// Register a button that is triggered when
	// the state of 'Space' or 'Enter' is changed.
	engo.Input.ButtonMgr.SetButton("action", act.KeySpace, act.KeyEnter)

	// Register an axis where 'A' will return a
	// negative value and 'D' returns a positive value.
	engo.Input.AxisMgr.SetAxis("sideways", act.AxisPair{act.KeyA, act.KeyD})

	w.AddSystem(&common.RenderSystem{})

	// Register the input system responding
	// to the input actions registered above.
	w.AddSystem(&InputSystem{})
}

func (*DefaultScene) Type() string { return "Game" }

type inputEntity struct {
	*ecs.BasicEntity
}

type InputSystem struct {
	action   uintptr
	sideways uintptr
	entities []inputEntity
}

func (c *InputSystem) New(w *ecs.World) {
	c.action = engo.Input.ButtonMgr.GetId("action")
	c.sideways = engo.Input.AxisMgr.GetId("sideways")

	if 0 == c.action {
		fmt.Println("Action button not found, using default fallback!")
		c.action = engo.Input.ButtonMgr.SetButton("action", act.KeySpace, act.KeyEnter)
	}
	if 0 == c.sideways {
		fmt.Println("Sideway axis not found, using default fallback!")
		c.sideways = engo.Input.AxisMgr.SetAxis("sideways", act.AxisPair{act.KeyA, act.KeyD})
	}
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

	// Nitya Note: Make sure to try to press both
	// keys and just idle/active behave as expected ?

	// Look up the button state and act on it
	if engo.Input.ButtonMgr.Active(c.action) {
		fmt.Println("The action key is still down.")
	}
	if engo.Input.ButtonMgr.JustIdle(c.action) {
		fmt.Println("An action key was just released!")
	}
	if engo.Input.ButtonMgr.JustActive(c.action) {
		fmt.Println("An action key was just pressed!")
	}

	// Nitya Note: Make sure to try to press both up and
	// down at the same time, result is zero as expected ?

	// Check some states on the axis
	if engo.Input.AxisMgr.JustActive(c.sideways) {
		fmt.Println("The axis just woke up.")
	}
	if engo.Input.AxisMgr.MinJustActive(c.sideways) {
		fmt.Println(" - The axis is going down.")
	}
	if engo.Input.AxisMgr.MaxJustActive(c.sideways) {
		fmt.Println(" - The axis is going up.")
	}

	// Only log the axis value when active
	if engo.Input.AxisMgr.Active(c.sideways) {
		fmt.Printf(" --> Axis value: %.2f \n", engo.Input.AxisMgr.Value(c.sideways))
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
