package main

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/act"
)

type DefaultScene struct{}

func (*DefaultScene) Preload() {}

func (*DefaultScene) Setup(w *ecs.World) {
	// Register a button that is triggered when
	// the state of 'Space' or 'Enter' is changed.
	engo.Buttons.SetNamed("action", act.KeySpace, act.KeyEnter)

	// Register an axis where 'A' will return a
	// negative value and 'D' returns a positive value.
	engo.Axes.SetNamed("sideways", act.AxisPair{act.KeyA, act.KeyD})

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
	c.action = engo.Buttons.Id("action")
	c.sideways = engo.Axes.Id("sideways")

	if 0 == c.action {
		fmt.Println("Action button not found, using default fallback!")
		c.action = engo.Buttons.SetNamed("action", act.KeySpace, act.KeyEnter)
	}
	if 0 == c.sideways {
		fmt.Println("Sideway axis not found, using default fallback!")
		c.sideways = engo.Axes.SetNamed("sideways", act.AxisPair{act.KeyA, act.KeyD})
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

	// Note: Most of these states are mutualy exlusive, in
	// production code "else if" would be prefered! To demo
	// difrent behaviors the if statements are seperated here.

	// Checking specific input codes
	if engo.Input.JustIdle(act.KeyA) {
		fmt.Println("Key code A was just released!")
	} // else
	if engo.Input.JustActive(act.KeyA) {
		fmt.Println("Key code A was just pressed!")
	}

	// Note: Make sure to try to press both keys
	// and just idle/active behave as expected ?

	// Look up the button state and act on it
	if engo.Buttons.Active(c.action) {
		fmt.Println("The action key is still down.")
	} // else
	if engo.Buttons.JustIdle(c.action) {
		fmt.Println("An action key was just released!")
	} // else
	if engo.Buttons.JustActive(c.action) {
		fmt.Println("An action key was just pressed!")
	}

	// Note: Make sure to try to press both up and down
	// at the same time, the axis is zero as expected ?

	// Check some states on the axis
	if engo.Axes.JustIdle(c.sideways) {
		fmt.Println("The axis fell a sleep.")
	} // else
	if engo.Axes.JustActive(c.sideways) {
		fmt.Println("The axis just woke up.")
	}

	// Check the state on the min side
	if engo.Axes.MinJustIdle(c.sideways) {
		if engo.Axes.MaxActive(c.sideways) {
			fmt.Println(" - The axis is going up.")
		}
	} // else
	if engo.Axes.MinJustActive(c.sideways) {
		if !engo.Axes.MaxActive(c.sideways) {
			fmt.Println(" - The axis is going down.")
		} else {
			fmt.Println(" - The axis is neutralized.")
		}
	}

	// Lets do the same for the max side
	if engo.Axes.MaxJustIdle(c.sideways) {
		if engo.Axes.MinActive(c.sideways) {
			fmt.Println(" - The axis is going down.")
		}
	} // else
	if engo.Axes.MaxJustActive(c.sideways) {
		if !engo.Axes.MinActive(c.sideways) {
			fmt.Println(" - The axis is going up.")
		} else {
			fmt.Println(" - The axis is neutralized.")
		}
	}

	// Log when the axis is active
	if engo.Axes.Active(c.sideways) {
		fmt.Printf(" --> Axis value: %.2f \n", engo.Axes.Value(c.sideways))
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
