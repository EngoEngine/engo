package engo // import "engo.io/engo"

import (
	"fmt"
	"log"

	"engo.io/ecs"
)

var (
	Time  *Clock
	Files *Loader
	Input *InputManager

	closeGame          bool
	defaultCloseAction bool
	WorldBounds        AABB

	currentWorld *ecs.World
	currentScene Scene
	Mailbox      *MessageManager
	cam          *cameraSystem

	scaleOnResize   = false
	fpsLimit        = 60
	headless        = false
	vsync           = true
	resetLoopTicker = make(chan bool, 1)
)

const (
	DefaultVerticalAxis   = "vertical"
	DefaultHorizontalAxis = "horizontal"
)

type RunOptions struct {
	// NoRun indicates the Open function should return immediately, without looping
	NoRun bool

	// Title is the Window title
	Title string

	// HeadlessMode indicates whether or not OpenGL calls should be made
	HeadlessMode bool

	Fullscreen bool

	Width, Height int

	// VSync indicates whether or not OpenGL should wait for the monitor to swp the buffers
	VSync bool

	// ScaleOnResize indicates whether or not engo should make things larger/smaller whenever the screen resizes
	ScaleOnResize bool

	// FPSLimit indicates the maximum number of frames per second
	FPSLimit int

	// OverrideCloseAction indicates that (when true) engo will never close whenever the gamer wants to close the
	// game - that will be your responsibility
	OverrideCloseAction bool

	StandardInputs bool
}

func Run(opts RunOptions, defaultScene Scene) {
	// Save settings
	SetScaleOnResize(opts.ScaleOnResize)
	SetFPSLimit(opts.FPSLimit)
	vsync = opts.VSync
	defaultCloseAction = !opts.OverrideCloseAction

	// Create input
	Input = NewInputManager()
	if opts.StandardInputs {
		log.Println("Using standard inputs")

		Input.RegisterButton("jump", Space)
		Input.RegisterButton("action", Enter)

		Input.RegisterAxis(DefaultHorizontalAxis, AxisKeyPair{A, D}, AxisKeyPair{ArrowLeft, ArrowRight})
		Input.RegisterAxis(DefaultVerticalAxis, AxisKeyPair{W, S}, AxisKeyPair{ArrowUp, ArrowDown})
	}

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

func SetScaleOnResize(b bool) {
	scaleOnResize = b
}

func SetOverrideCloseAction(value bool) {
	defaultCloseAction = !value
}

func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}

func runHeadless(defaultScene Scene) {
	runLoop(defaultScene, true)
}
