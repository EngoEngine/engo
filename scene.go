package engo

import (
	"fmt"
	"reflect"
)

var scenes = make(map[string]*sceneWrapper)

// Scene represents a screen ingame.
// i.e.: main menu, settings, but also the game itself
type Scene interface {
	// Preload is called before loading resources
	Preload()

	// Setup is called before the main loop
	Setup(Updater)

	// Type returns a unique string representation of the Scene, used to identify it
	Type() string
}

// Shower is an optional interface a Scene can implement, indicating it'll have custom behavior
// whenever the Scene gets shown again after being hidden (due to switching to other Scenes)
type Shower interface {
	// Show is called whenever the other Scene becomes inactive, and this one becomes the active one
	Show()
}

// Hider is an optional interface a Scene can implement, indicating it'll have custom behavior
// whenever the Scene get hidden to make room fr other Scenes.
type Hider interface {
	// Hide is called when an other Scene becomes active
	Hide()
}

// Exiter is an optional interface a Scene can implement, indicating it'll have custom behavior
// whenever the game get closed.
type Exiter interface {
	// Exit is called when the user or the system requests to close the game
	// This should be used to cleanup or prompt user if they're sure they want to close
	// To prevent the default action (close/exit) make sure to set OverrideCloseAction in
	// your RunOpts to `true`. You should then handle the exiting of the program by calling
	//    engo.Exit()
	Exit()
}

// Updater is an interface for what handles your game's Update during each frame.
// typically, this will be an *ecs.World, but you can implement your own Upodater
// and use engo without using engo's ecs
type Updater interface {
	Update(float32)
}

type sceneWrapper struct {
	scene   Scene
	update  Updater
	mailbox *MessageManager
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
		if hider, ok := currentScene.(Hider); ok {
			hider.Hide()
		}
	}

	// Register Scene if needed
	wrapper, registered := scenes[s.Type()]
	if !registered {
		RegisterScene(s)
		wrapper = scenes[s.Type()]
	}

	// Initialize new Scene / World if needed
	var doSetup bool

	if wrapper.update == nil || forceNewWorld {
		t := reflect.Indirect(reflect.ValueOf(currentUpdater)).Type()
		v := reflect.New(t)
		wrapper.update = v.Interface().(Updater)
		wrapper.mailbox = &MessageManager{}

		doSetup = true
	}

	// Do the switch
	currentScene = s
	currentUpdater = wrapper.update
	Mailbox = wrapper.mailbox

	// doSetup is true whenever we're (re)initializing the Scene
	if doSetup {
		s.Preload()

		wrapper.mailbox.listeners = make(map[string][]MessageHandler)

		s.Setup(wrapper.update)
	} else {
		if shower, ok := currentScene.(Shower); ok {
			shower.Show()
		}
	}
}

// RegisterScene registers the `Scene`, so it can later be used by `SetSceneByName`
func RegisterScene(s Scene) {
	_, ok := scenes[s.Type()]
	if !ok {
		scenes[s.Type()] = &sceneWrapper{scene: s}
	}
}

// SetSceneByName does a lookup for the `Scene` where its `Type()` equals `name`, and then sets it as current `Scene`
func SetSceneByName(name string, forceNewWorld bool) error {
	scene, ok := scenes[name]
	if !ok {
		return fmt.Errorf("scene not registered: %s", name)
	}

	SetScene(scene.scene, forceNewWorld)

	return nil
}
