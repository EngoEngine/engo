package engo // import "engo.io/engo"

import (
	"fmt"
	"log"

	"engo.io/ecs"
)

var (
	// Time is the active FPS counter
	Time *Clock

	// Input handles all input: mouse, keyboard and touch
	Input *InputManager

	// Mailbox is used by all Systems to communicate
	Mailbox *MessageManager

	currentWorld *ecs.World
	currentScene Scene

	opts                      RunOptions
	resetLoopTicker           = make(chan bool, 1)
	closeGame                 bool
	gameWidth, gameHeight     float32
	windowWidth, windowHeight float32
)

const (
	// DefaultVerticalAxis is the name of the default vertical axis, as used internally in `engo` when `StandardInputs`
	// is defined.
	DefaultVerticalAxis = "vertical"

	// DefaultHorizontalAxis is the name of the default horizontal axis, as used internally in `engo` when `StandardInputs`
	// is defined.
	DefaultHorizontalAxis = "horizontal"
	DefaultMouseXAxis     = "mouse x"
	DefaultMouseYAxis     = "mouse y"
)

type RunOptions struct {
	// NoRun indicates the Open function should return immediately, without looping
	NoRun bool

	// Title is the Window title
	Title string

	// HeadlessMode indicates whether or not OpenGL calls should be made
	HeadlessMode bool

	// Fullscreen indicates the game should run in fullscreen mode if run on a desktop
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

	// AssetsRoot is the path where all resources (images, audio files, fonts, etc.) can be found. Leaving this at
	// empty-string, will default this to `assets`.
	//
	// Whenever using any value that does not start with the directory `assets`, you will not be able to support
	// mobile (Android/iOS), because they **require** all assets to be within the `assets` directory. You may however
	// use any subfolder-structure within that `assets` directory.
	AssetsRoot string
}

// Run is called to create a window, initialize everything, and start the main loop. Once this function returns,
// the game window has been closed already. You can supply a lot of options within `RunOptions`, and your starting
// `Scene` should be defined in `defaultScene`.
func Run(o RunOptions, defaultScene Scene) {
	// Setting defaults
	if o.FPSLimit == 0 {
		o.FPSLimit = 60
	}

	if o.MSAA < 0 {
		panic("MSAA has to be greater or equal to 0")
	}

	if o.MSAA == 0 {
		o.MSAA = 1
	}

	if len(o.AssetsRoot) == 0 {
		o.AssetsRoot = "assets"
	}

	opts = o

	// Create input
	Input = NewInputManager()
	if opts.StandardInputs {
		log.Println("Using standard inputs")

		Input.RegisterButton("jump", Space)
		Input.RegisterButton("action", Enter)

		Input.RegisterAxis(DefaultHorizontalAxis, AxisKeyPair{A, D}, AxisKeyPair{ArrowLeft, ArrowRight})
		Input.RegisterAxis(DefaultVerticalAxis, AxisKeyPair{W, S}, AxisKeyPair{ArrowUp, ArrowDown})

		Input.RegisterAxis(DefaultMouseXAxis, NewAxisMouse(AxisMouseHori))
		Input.RegisterAxis(DefaultMouseYAxis, NewAxisMouse(AxisMouseVert))
	}

	Files.SetRoot(opts.AssetsRoot)

	// And run the game
	if opts.HeadlessMode {
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

// SetScaleOnResize can be used to change the value in the given `RunOpts` after already having called `engo.Run`.
func SetScaleOnResize(b bool) {
	opts.ScaleOnResize = b
}

// SetOverrideCloseAction can be used to change the value in the given `RunOpts` after already having called `engo.Run`.
func SetOverrideCloseAction(value bool) {
	opts.OverrideCloseAction = value
}

// SetFPSLimit can be used to change the value in the given `RunOpts` after already having called `engo.Run`.
func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	opts.FPSLimit = limit
	resetLoopTicker <- true
	return nil
}

// SetHeadless sets the headless-mode variable - should be used before calling `Run`, and will be overridden by the
// value within the `RunOpts` once you call `engo.Run`.
func SetHeadless(b bool) {
	opts.HeadlessMode = b
}

// Headless indicates whether or not OpenGL-calls should be made
func Headless() bool {
	return opts.HeadlessMode
}

// ScaleOnResizes indicates whether or not the screen should resize (i.e. make things look smaller/bigger) whenever
// the window resized. If `false`, then the size of the screen does not affect the size of the things drawn - it just
// makes less/more objects visible
func ScaleOnResize() bool {
	return opts.ScaleOnResize
}

// Exit is the safest way to close your game, as `engo` will correctly attempt to close all windows, handlers and contexts
func Exit() {
	closeGame = true
}

// GameWidth returns the current game width
func GameWidth() float32 {
	return gameWidth
}

// GameHeight returns the current game height
func GameHeight() float32 {
	return gameHeight
}

func closeEvent() {
	for _, scenes := range scenes {
		if exiter, ok := scenes.scene.(Exiter); ok {
			exiter.Exit()
		}
	}

	if !opts.OverrideCloseAction {
		Exit()
	} else {
		log.Println("[WARNING] default close action set to false, please make sure you manually handle this")
	}
}

func runHeadless(defaultScene Scene) {
	runLoop(defaultScene, true)
}
