//+build demo

package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	iconScene *IconScene
	rockScene *RockScene
)

type Guy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Rock struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// IconScene is responsible for managing the icon
type IconScene struct{}

func (*IconScene) Preload() {
	err := engo.Files.Load("icon.png")
	if err != nil {
		log.Println(err)
	}
}

func (*IconScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "RockScene", WaitTime: time.Second * 3})

	// Retrieve a texture
	texture, err := common.LoadedSprite("icon.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy := Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * guy.RenderComponent.Scale.X,
		Height:   texture.Height() * guy.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *ScaleSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent)
		}
	}
}

func (*IconScene) Hide() {
	log.Println("IconScene is now hidden")
}

func (*IconScene) Show() {
	log.Println("IconScene is now shown")
}

func (*IconScene) Type() string { return "IconScene" }

// RockScene is responsible for managing the rock
type RockScene struct{}

func (*RockScene) Preload() {
	err := engo.Files.Load("rock.png")
	if err != nil {
		log.Println(err)
	}
}

func (game *RockScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&ScaleSystem{})
	w.AddSystem(&SceneSwitcherSystem{NextScene: "IconScene", WaitTime: time.Second * 3})

	// Retrieve a texture
	texture, err := common.LoadedSprite("rock.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	rock := Rock{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	rock.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{8, 8},
	}
	rock.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * rock.RenderComponent.Scale.X,
		Height:   texture.Height() * rock.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *ScaleSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent)
		}
	}
}

func (*RockScene) Hide() {
	log.Println("RockScene is now hidden")
}

func (*RockScene) Show() {
	log.Println("RockScens is now shown")
}

func (*RockScene) Type() string { return "RockScene" }

// SceneSwitcherSystem is a System that actually calls SetScene
type SceneSwitcherSystem struct {
	NextScene     string
	WaitTime      time.Duration
	secondsWaited float32
}

func (*SceneSwitcherSystem) Priority() int          { return 1 }
func (*SceneSwitcherSystem) Remove(ecs.BasicEntity) {}

func (s *SceneSwitcherSystem) Update(dt float32) {
	s.secondsWaited += dt
	if float64(s.secondsWaited) > s.WaitTime.Seconds() {
		s.secondsWaited = 0

		// Change the world to s.NextScene, and don't override / force World re-creation
		engo.SetSceneByName(s.NextScene, false)

		log.Println("Switched to", s.NextScene)
	}
}

type scaleEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
}

type ScaleSystem struct {
	entities []scaleEntity
}

func (s *ScaleSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent) {
	s.entities = append(s.entities, scaleEntity{basic, render})
}

func (s *ScaleSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *ScaleSystem) Update(dt float32) {
	for _, e := range s.entities {
		var mod float32

		if rand.Int()%2 == 0 {
			mod = 0.1
		} else {
			mod = -0.1
		}

		if e.RenderComponent.Scale.X+mod >= 15 || e.RenderComponent.Scale.X+mod <= 1 {
			mod *= -1
		}

		newScale := e.RenderComponent.Scale
		newScale.AddScalar(mod)
		e.RenderComponent.Scale = newScale
	}
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
