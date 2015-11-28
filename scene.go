package engi

import "github.com/paked/engi/ecs"

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
}
