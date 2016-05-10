package engo // import "engo.io/engo"

import (
	"fmt"
	"log"

	"engo.io/ecs"
	"image/color"
)

var (
	Time *Clock
	//Files *Loader
	Input *InputManager

	closeGame          bool
	defaultCloseAction bool
	WorldBounds        AABB

	currentWorld *ecs.World
	currentScene Scene
	Mailbox      *MessageManager

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

	// StandardInputs is an easy way to map common inputs to actions, such as "jump" being <SPACE>, and "action" being
	// <ENTER>.
	StandardInputs bool

	// MSAA indicates the amount of samples that should be taken. Leaving it blank will default to 1, and you may
	// use any positive value you wish. It may be possible that the operating system / environment doesn't support
	// the requested amount. In that case, GLFW will (hopefully) pick the highest supported sampling count. The higher
	// the value, the bigger the performance cost.
	//
	// Our `RenderSystem` automatically calls `gl.Enable(gl.MULTISAMPLE)` (which is required to make use of it), but
	// if you're going to use your own rendering `System` instead, you will have to call it yourself.
	//
	// Also note that this value is entirely ignored in WebGL - most browsers enable it by default when available, and
	// none of them (at time of writing) allow you to tune it.
	//
	// More info at https://www.opengl.org/wiki/Multisampling
	// "With multisampling, each pixel at the edge of a polygon is sampled multiple times."
	MSAA int
}

// Run is called to create a window, initialize everything, and start the main loop. Once this function returns,
// the game window has been closed already. You can supply a lot of options within `RunOptions`, and your starting
// `Scene` should be defined in `defaultScene`.
func Run(opts RunOptions, defaultScene Scene) {
	// Save settings
	SetScaleOnResize(opts.ScaleOnResize)
	SetFPSLimit(opts.FPSLimit)
	vsync = opts.VSync
	defaultCloseAction = !opts.OverrideCloseAction
	if opts.FPSLimit > 0 {
		fpsLimit = opts.FPSLimit
	}

	// Create input
	Input = NewInputManager()
	if opts.StandardInputs {
		log.Println("Using standard inputs")

		Input.RegisterButton("jump", Space)
		Input.RegisterButton("action", Enter)

		Input.RegisterAxis(DefaultHorizontalAxis, AxisKeyPair{A, D}, AxisKeyPair{ArrowLeft, ArrowRight})
		Input.RegisterAxis(DefaultVerticalAxis, AxisKeyPair{W, S}, AxisKeyPair{ArrowUp, ArrowDown})
	}

	if opts.MSAA < 0 {
		panic("MSAA has to be greater or equal to 0")
	}

	if opts.MSAA == 0 {
		opts.MSAA = 1
	}

	if opts.HeadlessMode {
		headless = true

		if !opts.NoRun {
			runHeadless(defaultScene)
		}
	} else {
		CreateWindow(opts.Title, opts.Width, opts.Height, opts.Fullscreen, opts.MSAA)
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

func SetBackground(c color.Color) {
	if !headless {
		r, g, b, a := c.RGBA()

		Gl.ClearColor(float32(r), float32(g), float32(b), float32(a))
	}
}

func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}

// Headless indicates whether or not OpenGL-calls should be made
func Headless() bool {
	return headless
}

// ScaleOnResizes indicates whether or not the screen should resize (i.e. make things look smaller/bigger) whenever
// the window resized. If `false`, then the size of the screen does not affect the size of the things drawn - it just
// makes less/more objects visible
func ScaleOnResize() bool {
	return scaleOnResize
}

func runHeadless(defaultScene Scene) {
	runLoop(defaultScene, true)
}

func Exit() {
	closeGame = true
}

func closeEvent() {
	for _, scenes := range scenes {
		if exiter, ok := scenes.scene.(Exiter); ok {
			exiter.Exit()
		}
	}

	if defaultCloseAction {
		Exit()
	} else {
		log.Println("Warning: default close action set to false, please make sure you manually handle this")
	}
}
