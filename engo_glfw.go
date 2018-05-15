// +build darwin,!arm,!arm64 linux windows
// +build !ios,!android,!netgo

package engo

import (
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

	scale = float32(1)
)

func init() {
	runtime.LockOSThread()
}

// fatalErr calls log.Fatal with the given error if it is non-nil.
func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CreateWindow sets up the GLFW window and prepares the OpenGL surface for rendering
func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	CurrentBackEnd = BackEndGLFW
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
	if opts.NotResizable {
		glfw.WindowHint(glfw.Resizable, glfw.False)
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

	fw, fh := window.GetFramebufferSize()
	canvasWidth, canvasHeight = float32(fw), float32(fh)

	if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
		scale = canvasWidth / windowWidth
	}

	window.SetFramebufferSizeCallback(func(window *glfw.Window, w, h int) {
		Gl.Viewport(0, 0, w, h)
		width, height = window.GetSize()
		windowWidth, windowHeight = float32(width), float32(width)

		oldCanvasW, oldCanvasH := canvasWidth, canvasHeight

		canvasWidth, canvasHeight = float32(w), float32(h)

		ResizeXOffset += oldCanvasW - canvasWidth
		ResizeYOffset += oldCanvasH - canvasHeight

		if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
			scale = canvasWidth / windowWidth
		}
	})

	window.SetCursorPosCallback(func(window *glfw.Window, x, y float64) {
		Input.Mouse.X, Input.Mouse.Y = float32(x)/opts.GlobalScale.X, float32(y)/opts.GlobalScale.Y
		if Input.Mouse.Action != Release && Input.Mouse.Action != Press {
			Input.Mouse.Action = Move
		}
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPos()
		Input.Mouse.X, Input.Mouse.Y = float32(x)/(opts.GlobalScale.X), float32(y)/(opts.GlobalScale.Y)

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
		fw, fh := window.GetFramebufferSize()
		canvasWidth, canvasHeight = float32(fw), float32(fh)

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

// DestroyWindow handles the termination of windows
func DestroyWindow() {
	glfw.Terminate()
}

// SetTitle sets the title of the window
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
	currentUpdater.Update(Time.Delta())

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

// CursorPos returns the current cursor position
func CursorPos() (x, y float32) {
	w, h := window.GetCursorPos()
	return float32(w), float32(h)
}

// WindowSize gets the current window size
func WindowSize() (w, h int) {
	return window.GetSize()
}

// WindowWidth gets the current window width
func WindowWidth() float32 {
	return windowWidth
}

// WindowHeight gets the current window height
func WindowHeight() float32 {
	return windowHeight
}

// CanvasWidth gets the width of the current OpenGL Framebuffer
func CanvasWidth() float32 {
	return canvasWidth
}

// CanvasHeight gets the height of the current OpenGL Framebuffer
func CanvasHeight() float32 {
	return canvasHeight
}

// CanvasScale gets the ratio of the canvas to the window sizes
func CanvasScale() float32 {
	return scale
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

// SetVSync sets whether or not to use VSync
func SetVSync(enabled bool) {
	opts.VSync = enabled
	if opts.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

//SetCursorVisibility sets the visibility of the cursor.
//If true the cursor is visible, if false the cursor is not.
func SetCursorVisibility(visible bool) {
	if visible {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

// openFile is the desktop-specific way of opening a file
func openFile(url string) (io.ReadCloser, error) {
	return os.Open(url)
}

// IsAndroidChrome tells if the browser is Chrome for android
func IsAndroidChrome() bool {
	return false
}
