package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

var (
	iconScene *IconScene
	rockScene *RockScene
)

// IconScene is responsible for managing the icon
type IconScene struct{}

func (game *IconScene) Preload() {
	engi.Files.Add("data/icon.png")
}

func (game *IconScene) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "RockScene", WaitTime: time.Second * 3})

	guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engi.Files.Image("icon.png")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "icon")
	collision := &engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	w.AddEntity(guy)
}

func (*IconScene) Hide()        {}
func (*IconScene) Show()        {}
func (*IconScene) Type() string { return "IconScene" }

// RockScene is responsible for managing the rock
type RockScene struct{}

func (*RockScene) Preload() {
	engi.Files.Add("data/rock.png")
}

func (game *RockScene) Setup(w *ecs.World) {
	engi.SetBg(0x2d3739)

	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "IconScene", WaitTime: time.Second * 3})

	guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engi.Files.Image("rock.png")
	render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "rock")
	collision := &engi.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	w.AddEntity(guy)
}

func (*RockScene) Hide()        {}
func (*RockScene) Show()        {}
func (*RockScene) Type() string { return "RockScene" }

// SceneSwitcherSystem is a System that actually calls SetScene
type SceneSwitcherSystem struct {
	*ecs.System

	NextScene     string
	WaitTime      time.Duration
	secondsWaited float32
}

func (*SceneSwitcherSystem) Type() string {
	return "SceneSwitcherSystem"
}

func (s *SceneSwitcherSystem) New(*ecs.World) {
	s.System = ecs.NewSystem()
	s.System.AddEntity(ecs.NewEntity([]string{s.Type()}))
}

func (s *SceneSwitcherSystem) Update(e *ecs.Entity, dt float32) {
	s.secondsWaited += dt
	if float64(s.secondsWaited) > s.WaitTime.Seconds() {
		s.secondsWaited = 0

		// Change the world to s.NextScene, and don't override / force World re-creation
		engi.SetSceneByName(s.NextScene, false)

		fmt.Println("Switched to", s.NextScene)
	}
}

// ScaleSystem is the System which scales the Entities inside
type ScaleSystem struct {
	*ecs.System
}

func (ScaleSystem) Type() string {
	return "ScaleSystem"
}

func (s *ScaleSystem) New(*ecs.World) {
	s.System = ecs.NewSystem()
}

func (c *ScaleSystem) Update(e *ecs.Entity, dt float32) {
	var render *engi.RenderComponent
	if !e.Component(&render) {
		return
	}
	var mod float32

	if rand.Int()%2 == 0 {
		mod = 0.1
	} else {
		mod = -0.1
	}

	if render.Scale().X+mod >= 15 || render.Scale().X+mod <= 1 {
		mod *= -1
	}

	newScale := render.Scale()
	newScale.AddScalar(mod)
	render.SetScale(newScale)
}

func main() {
	iconScene = &IconScene{}
	rockScene = &RockScene{}

	// Register other Scenes for later use, this can be done from anywhere, as long as it
	// happens before calling engi.SetSceneByName
	engi.RegisterScene(rockScene)

	opts := engi.RunOptions{
		Title:  "Scenes Demo",
		Width:  1024,
		Height: 640,
	}

	engi.Open(opts, iconScene)
}
