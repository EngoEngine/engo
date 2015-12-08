package engi

import (
	"fmt"
	"github.com/paked/engi/ecs"
)

var (
	worlds = make(map[string]*ecs.World)
	scenes = make(map[string]Scene)
)

// Scene represents a screen ingame.
// i.e.: main menu, settings, but also the game itself
type Scene interface {
	// Preload is called before loading resources
	Preload()

	// Setup is called before the main loop
	Setup(*ecs.World)

	// Show is called whenever the other Scene becomes inactive, and this one becomes the active one
	Show()

	// Hide is called when an other Scene becomes active
	Hide()

	// Type returns a unique string representation of the Scene, used to identify it
	Type() string
}

// CurrentScene returns the SceneWorld that is currently active
func CurrentScene() Scene {
	return currentScene
}

// SetScene sets the currentScene to the given Scene, and
// optionally forcing to create a new ecs.World that goes with it.
func SetScene(s Scene, forceNewWorld bool) {
	// Break down currentScene
	if currentScene != nil {
		currentScene.Hide()
	}

	// Register Scene if needed
	_, registered := scenes[s.Type()]
	if !registered {
		RegisterScene(s)
	}

	// Initialize new Scene / World if needed
	var newWorld *ecs.World
	if w, ok := worlds[s.Type()]; !ok || forceNewWorld {
		s.Preload()
		Files.Load(func() {})

		newWorld = &ecs.World{}
		newWorld.New()
		worlds[s.Type()] = newWorld

		s.Setup(newWorld)
	} else {
		newWorld = w
	}

	// Do the switch
	currentWorld = newWorld
	currentScene = s
}

// RegisterScene registers the `Scene`, so it can later be used by `SetSceneByName`
func RegisterScene(s Scene) {
	scenes[s.Type()] = s
}

// SetSceneByName does a lookup for the `Scene` where its `Type()` equals `name`, and then sets it as current `Scene`
func SetSceneByName(name string, forceNewWorld bool) error {
	scene, ok := scenes[name]
	if !ok {
		return fmt.Errorf("scene not registered: %s", name)
	}

	SetScene(scene, forceNewWorld)

	return nil
}
