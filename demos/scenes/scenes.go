package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/engoengine/engo"
	"github.com/engoengine/ecs"
)

var (
	iconScene *IconScene
	rockScene *RockScene
)

// IconScene is responsible for managing the icon
type IconScene struct{}

func (game *IconScene) Preload() {
	engo.Files.Add("data/icon.png")
}

func (game *IconScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "RockScene", WaitTime: time.Second * 3})

	guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engo.Files.Image("icon.png")
	render := engo.NewRenderComponent(texture, engo.Point{8, 8}, "icon")
	collision := &engo.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engo.SpaceComponent{engo.Point{(engo.Width() - width) / 2, (engo.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	err := w.AddEntity(guy)
	if err != nil {
		log.Println(err)
	}
}

func (*IconScene) Hide()        {}
func (*IconScene) Show()        {}
func (*IconScene) Type() string { return "IconScene" }

// RockScene is responsible for managing the rock
type RockScene struct{}

func (*RockScene) Preload() {
	engo.Files.Add("data/rock.png")
}

func (game *RockScene) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "IconScene", WaitTime: time.Second * 3})

	guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})
	texture := engo.Files.Image("rock.png")
	render := engo.NewRenderComponent(texture, engo.Point{8, 8}, "rock")
	collision := &engo.CollisionComponent{Solid: true, Main: true}

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engo.SpaceComponent{engo.Point{(engo.Width() - width) / 2, (engo.Height() - height) / 2}, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)
	guy.AddComponent(collision)

	err := w.AddEntity(guy)
	if err != nil {
		log.Println(err)
	}
}

func (*RockScene) Hide()        {}
func (*RockScene) Show()        {}
func (*RockScene) Type() string { return "RockScene" }

// SceneSwitcherSystem is a System that actually calls SetScene
type SceneSwitcherSystem struct {
	NextScene     string
	WaitTime      time.Duration
	secondsWaited float32
}

func (*SceneSwitcherSystem) Type() string             { return "SceneSwitcherSystem" }
func (*SceneSwitcherSystem) Priority() int            { return 1 }
func (*SceneSwitcherSystem) AddEntity(*ecs.Entity)    {}
func (*SceneSwitcherSystem) RemoveEntity(*ecs.Entity) {}
func (*SceneSwitcherSystem) New(*ecs.World)           {}

func (s *SceneSwitcherSystem) Update(dt float32) {
	s.secondsWaited += dt
	if float64(s.secondsWaited) > s.WaitTime.Seconds() {
		s.secondsWaited = 0

		// Change the world to s.NextScene, and don't override / force World re-creation
		engo.SetSceneByName(s.NextScene, false)

		fmt.Println("Switched to", s.NextScene)
	}
}

// ScaleSystem is the System which scales the Entities inside
type ScaleSystem struct {
	ecs.LinearSystem
}

func (*ScaleSystem) Type() string { return "ScaleSystem" }

func (s *ScaleSystem) New(*ecs.World) {}

func (c *ScaleSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var render *engo.RenderComponent
	if !entity.Component(&render) {
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
	// happens before calling engo.SetSceneByName
	engo.RegisterScene(rockScene)

	opts := engo.RunOptions{
		Title:  "Scenes Demo",
		Width:  1024,
		Height: 640,
	}

	engo.Run(opts, iconScene)
}
