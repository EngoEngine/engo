//+build !netgo,!android

package engo

import (
	"image"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"engo.io/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	window *glfw.Window
	// Gl is the current OpenGL context
	Gl *gl.Context

	cursorArrow     *glfw.Cursor
	cursorIBeam     *glfw.Cursor
	cursorCrosshair *glfw.Cursor
	cursorHand      *glfw.Cursor
	cursorHResize   *glfw.Cursor
	cursorVResize   *glfw.Cursor

	headlessWidth             = 800
	headlessHeight            = 800
	canvasWidth, canvasHeight float32
	scale                     = float32(1)
)

// fatalErr calls log.Fatal with the given error if it is non-nil.
func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	err := glfw.Init()
	fatalErr(err)

	cursorArrow = glfw.CreateStandardCursor(int(glfw.ArrowCursor))
	cursorIBeam = glfw.CreateStandardCursor(int(glfw.IBeamCursor))
	cursorCrosshair = glfw.CreateStandardCursor(int(glfw.CrosshairCursor))
	cursorHand = glfw.CreateStandardCursor(int(glfw.HandCursor))
	cursorHResize = glfw.CreateStandardCursor(int(glfw.HResizeCursor))
	cursorVResize = glfw.CreateStandardCursor(int(glfw.VResizeCursor))

	monitor := glfw.GetPrimaryMonitor()

	var mode *glfw.VidMode
	if monitor != nil {
		mode = monitor.GetVideoMode()
	} else {
		// Initialize default values if no monitor is found
		mode = &glfw.VidMode{
			Width:       1,
			Height:      1,
			RedBits:     8,
			GreenBits:   8,
			BlueBits:    8,
			RefreshRate: 60,
		}
	}

	gameWidth = float32(width)
	gameHeight = float32(height)

	if fullscreen {
		width = mode.Width
		height = mode.Height
		glfw.WindowHint(glfw.Decorated, 0)
	} else {
		monitor = nil
	}

	if opts.HeadlessMode {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Samples, msaa)

	window, err = glfw.CreateWindow(width, height, title, monitor, nil)
	fatalErr(err)

	window.MakeContextCurrent()

	if !fullscreen {
		window.SetPos((mode.Width-width)/2, (mode.Height-height)/2)
	}

	SetVSync(opts.VSync)

	Gl = gl.NewContext()

	width, height = window.GetSize()
	windowWidth, windowHeight = float32(width), float32(height)

	vp := Gl.GetViewport()
	canvasWidth, canvasHeight = float32(vp[2]), float32(vp[3])

	if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
		scale = canvasWidth / windowWidth
	}

	window.SetFramebufferSizeCallback(func(window *glfw.Window, w, h int) {
		windowWidth, windowHeight = float32(w), float32(h)
		Gl.Viewport(0, 0, w, h)

		vp := Gl.GetViewport()
		canvasWidth, canvasHeight = float32(vp[2]), float32(vp[3])

		if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
			scale = canvasWidth / windowWidth
		}
	})

	window.SetCursorPosCallback(func(window *glfw.Window, x, y float64) {
		Input.Mouse.X, Input.Mouse.Y = float32(x), float32(y)
		Input.Mouse.Action = Move
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPos()
		Input.Mouse.X, Input.Mouse.Y = float32(x), float32(y)

		// this is only valid because we use an internal structure that is
		// 100% compatible with glfw3.h
		Input.Mouse.Button = MouseButton(b)
		Input.Mouse.Modifer = Modifier(m)

		if a == glfw.Press {
			Input.Mouse.Action = Press
		} else {
			Input.Mouse.Action = Release
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		Input.Mouse.ScrollX = float32(xoff)
		Input.Mouse.ScrollY = float32(yoff)
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		key := Key(k)
		if a == glfw.Press {
			Input.keys.Set(key, true)
		} else if a == glfw.Release {
			Input.keys.Set(key, false)
		}
	})

	window.SetSizeCallback(func(w *glfw.Window, widthInt int, heightInt int) {
		message := WindowResizeMessage{
			OldWidth:  int(windowWidth),
			OldHeight: int(windowHeight),
			NewWidth:  widthInt,
			NewHeight: heightInt,
		}

		windowWidth = float32(widthInt)
		windowHeight = float32(heightInt)

		// TODO: verify these for retina displays & verify if needed here
		vp := Gl.GetViewport()
		canvasWidth, canvasHeight = float32(vp[2]), float32(vp[3])

		if !opts.ScaleOnResize {
			gameWidth, gameHeight = float32(widthInt), float32(heightInt)
		}

		Mailbox.Dispatch(message)
	})

	window.SetCharCallback(func(window *glfw.Window, char rune) {
		// TODO: what does this do, when can we use it?
		// it's like KeyCallback, but for specific characters instead of keys...?
		// responder.Type(char)
	})
}

func DestroyWindow() {
	glfw.Terminate()
}

func SetTitle(title string) {
	if opts.HeadlessMode {
		log.Println("Title set to:", title)
	} else {
		window.SetTitle(title)
	}
}

// RunIteration runs one iteration per frame
func RunIteration() {
	Time.Tick()

	// First check for new keypresses
	if !opts.HeadlessMode {
		Input.update()
		glfw.PollEvents()
	}

	// Then update the world and all Systems
	currentWorld.Update(Time.Delta())

	// Lastly, forget keypresses and swap buffers
	if !opts.HeadlessMode {
		// reset values to avoid catching the same "signal" twice
		Input.Mouse.ScrollX, Input.Mouse.ScrollY = 0, 0
		Input.Mouse.Action = Neutral

		window.SwapBuffers()
	}

}

// RunPreparation is called automatically when calling Open. It should only be called once.
func RunPreparation(defaultScene Scene) {
	Time = NewClock()

	// Default WorldBounds values
	//WorldBounds.Max = Point{GameWidth(), GameHeight()}
	// TODO: move this to appropriate location

	SetScene(defaultScene, false)
}

func runLoop(defaultScene Scene, headless bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		closeEvent()
	}()

	RunPreparation(defaultScene)
	ticker := time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))

	// Start tick, minimize the delta
	Time.Tick()

Outer:
	for {
		select {
		case <-ticker.C:
			RunIteration()
			if closeGame {
				break Outer
			}
			if !headless && window.ShouldClose() {
				closeEvent()
			}
		case <-resetLoopTicker:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))
		}
	}
	ticker.Stop()
}

func CursorPos() (x, y float32) {
	w, h := window.GetCursorPos()
	return float32(w) * scale, float32(h) * scale
}

func WindowSize() (w, h int) {
	return window.GetSize()
}

func WindowWidth() float32 {
	return windowWidth
}

func WindowHeight() float32 {
	return windowHeight
}

func CanvasWidth() float32 {
	return canvasWidth
}

func CanvasHeight() float32 {
	return canvasHeight
}

// SetCursor sets the pointer of the mouse to the defined standard cursor
func SetCursor(c Cursor) {
	var cur *glfw.Cursor
	switch c {
	case CursorNone:
		cur = nil // just for the documentation, this isn't required
	case CursorArrow:
		cur = cursorArrow
	case CursorCrosshair:
		cur = cursorCrosshair
	case CursorHand:
		cur = cursorHand
	case CursorIBeam:
		cur = cursorIBeam
	case CursorHResize:
		cur = cursorHResize
	case CursorVResize:
		cur = cursorVResize
	}
	window.SetCursor(cur)
}

func SetVSync(enabled bool) {
	opts.VSync = enabled
	if opts.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

func init() {
	runtime.LockOSThread()

	Grave = Key(glfw.KeyGraveAccent)
	Dash = Key(glfw.KeyMinus)
	Apostrophe = Key(glfw.KeyApostrophe)
	Semicolon = Key(glfw.KeySemicolon)
	Equals = Key(glfw.KeyEqual)
	Comma = Key(glfw.KeyComma)
	Period = Key(glfw.KeyPeriod)
	Slash = Key(glfw.KeySlash)
	Backslash = Key(glfw.KeyBackslash)
	Backspace = Key(glfw.KeyBackspace)
	Tab = Key(glfw.KeyTab)
	CapsLock = Key(glfw.KeyCapsLock)
	Space = Key(glfw.KeySpace)
	Enter = Key(glfw.KeyEnter)
	Escape = Key(glfw.KeyEscape)
	Insert = Key(glfw.KeyInsert)
	PrintScreen = Key(glfw.KeyPrintScreen)
	Delete = Key(glfw.KeyDelete)
	PageUp = Key(glfw.KeyPageUp)
	PageDown = Key(glfw.KeyPageDown)
	Home = Key(glfw.KeyHome)
	End = Key(glfw.KeyEnd)
	Pause = Key(glfw.KeyPause)
	ScrollLock = Key(glfw.KeyScrollLock)
	ArrowLeft = Key(glfw.KeyLeft)
	ArrowRight = Key(glfw.KeyRight)
	ArrowDown = Key(glfw.KeyDown)
	ArrowUp = Key(glfw.KeyUp)
	LeftBracket = Key(glfw.KeyLeftBracket)
	LeftShift = Key(glfw.KeyLeftShift)
	LeftControl = Key(glfw.KeyLeftControl)
	LeftSuper = Key(glfw.KeyLeftSuper)
	LeftAlt = Key(glfw.KeyLeftAlt)
	RightBracket = Key(glfw.KeyRightBracket)
	RightShift = Key(glfw.KeyRightShift)
	RightControl = Key(glfw.KeyRightControl)
	RightSuper = Key(glfw.KeyRightSuper)
	RightAlt = Key(glfw.KeyRightAlt)
	Zero = Key(glfw.Key0)
	One = Key(glfw.Key1)
	Two = Key(glfw.Key2)
	Three = Key(glfw.Key3)
	Four = Key(glfw.Key4)
	Five = Key(glfw.Key5)
	Six = Key(glfw.Key6)
	Seven = Key(glfw.Key7)
	Eight = Key(glfw.Key8)
	Nine = Key(glfw.Key9)
	F1 = Key(glfw.KeyF1)
	F2 = Key(glfw.KeyF2)
	F3 = Key(glfw.KeyF3)
	F4 = Key(glfw.KeyF4)
	F5 = Key(glfw.KeyF5)
	F6 = Key(glfw.KeyF6)
	F7 = Key(glfw.KeyF7)
	F8 = Key(glfw.KeyF8)
	F9 = Key(glfw.KeyF9)
	F10 = Key(glfw.KeyF10)
	F11 = Key(glfw.KeyF11)
	F12 = Key(glfw.KeyF12)
	A = Key(glfw.KeyA)
	B = Key(glfw.KeyB)
	C = Key(glfw.KeyC)
	D = Key(glfw.KeyD)
	E = Key(glfw.KeyE)
	F = Key(glfw.KeyF)
	G = Key(glfw.KeyG)
	H = Key(glfw.KeyH)
	I = Key(glfw.KeyI)
	J = Key(glfw.KeyJ)
	K = Key(glfw.KeyK)
	L = Key(glfw.KeyL)
	M = Key(glfw.KeyM)
	N = Key(glfw.KeyN)
	O = Key(glfw.KeyO)
	P = Key(glfw.KeyP)
	Q = Key(glfw.KeyQ)
	R = Key(glfw.KeyR)
	S = Key(glfw.KeyS)
	T = Key(glfw.KeyT)
	U = Key(glfw.KeyU)
	V = Key(glfw.KeyV)
	W = Key(glfw.KeyW)
	X = Key(glfw.KeyX)
	Y = Key(glfw.KeyY)
	Z = Key(glfw.KeyZ)
	NumLock = Key(glfw.KeyNumLock)
	NumMultiply = Key(glfw.KeyKPMultiply)
	NumDivide = Key(glfw.KeyKPDivide)
	NumAdd = Key(glfw.KeyKPAdd)
	NumSubtract = Key(glfw.KeyKPSubtract)
	NumZero = Key(glfw.KeyKP0)
	NumOne = Key(glfw.KeyKP1)
	NumTwo = Key(glfw.KeyKP2)
	NumThree = Key(glfw.KeyKP3)
	NumFour = Key(glfw.KeyKP4)
	NumFive = Key(glfw.KeyKP5)
	NumSix = Key(glfw.KeyKP6)
	NumSeven = Key(glfw.KeyKP7)
	NumEight = Key(glfw.KeyKP8)
	NumNine = Key(glfw.KeyKP9)
	NumDecimal = Key(glfw.KeyKPDecimal)
	NumEnter = Key(glfw.KeyKPEnter)
}

func NewImageRGBA(img *image.RGBA) *ImageRGBA {
	return &ImageRGBA{img}
}

type ImageRGBA struct {
	data *image.RGBA
}

func (i *ImageRGBA) Data() interface{} {
	return i.data
}

func (i *ImageRGBA) Width() int {
	return i.data.Rect.Max.X
}

func (i *ImageRGBA) Height() int {
	return i.data.Rect.Max.Y
}

// openFile is the desktop-specific way of opening a file
func openFile(url string) (io.ReadCloser, error) {
	return os.Open(url)
}
