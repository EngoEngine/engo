// A cross-platform game engine written in Go following an interpretation
// of the Entity Component System paradigm. Engi is currently compilable for Mac OSX,
// Linux and Windows. Mobile and web(gopherjs) support are also in the works

package engi

import (
	"fmt"
	"image/color"

	"github.com/paked/engi/ecs"
	"github.com/paked/webgl"
)

var (
	Time        *Clock
	Files       *Loader
	Gl          *webgl.Context
	WorldBounds AABB

	currentWorld *ecs.World
	currentScene Scene
	Mailbox      *MessageManager
	cam          *cameraSystem

	scaleOnResize   = false
	fpsLimit        = 120
	headless        = false
	vsync           = true
	resetLoopTicker = make(chan bool, 1)
)

type RunOptions struct {
	// NoRun indicates the Open function should return immediately, without looping
	NoRun bool

	// Title is the Window title
	Title string

	// HeadlessMode indicates whether or not OpenGL calls should be made
	HeadlessMode bool

	Fullscreen    bool
	Width, Height int

	// VSync indicates whether or not OpenGL should wait for the monitor to swp the buffers
	VSync bool

	// ScaleOnResize indicates whether or not engi should make things larger/smaller whenever the screen resizes
	ScaleOnResize bool

	// FPSLimit indicates the maximum number of frames per second
	FPSLimit int
}

// Start up the engine with the specified configuration
func Open(opts RunOptions, defaultScene Scene) {
	// Save settings
	SetScaleOnResize(opts.ScaleOnResize)
	SetFPSLimit(opts.FPSLimit)
	vsync = opts.VSync

	if opts.HeadlessMode {
		headless = true

		if !opts.NoRun {
			runHeadless(defaultScene)
		}
	} else {
		CreateWindow(opts.Title, opts.Width, opts.Height, opts.Fullscreen)
		defer DestroyWindow()

		if !opts.NoRun {
			runLoop(defaultScene, false)
		}
	}
}

// Set the background color of the current scene to the specified color
func SetBg(c color.Color) {
	if !headless {
		r, g, b, a := c.RGBA()

		Gl.ClearColor(float32(r), float32(g), float32(b), float32(a))
	}
}

// When true, engi will automatically scale resources on screen resize
func SetScaleOnResize(b bool) {
	scaleOnResize = b
}

// Set the maximum frames per second (defaults to 120)
func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}
