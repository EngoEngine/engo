// +build headless

package engo

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"engo.io/gl"
)

var (
	// Gl is the current OpenGL context
	Gl *gl.Context

	scale = float32(1)
)

// CreateWindow sets up the GLFW window and prepares the OpenGL surface for rendering
func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	CurrentBackEnd = BackEndHeadless

	gameWidth = float32(width)
	gameHeight = float32(height)

	Gl = gl.NewContext()

	windowWidth, windowHeight = float32(width), float32(height)
	canvasWidth, canvasHeight = float32(width), float32(height)
}

// DestroyWindow handles the termination of windows
func DestroyWindow() {}

// SetTitle sets the title of the window
func SetTitle(title string) {
	log.Println("Title set to:", title)
}

// RunIteration runs one iteration per frame
func RunIteration() {
	Time.Tick()
	currentUpdater.Update(Time.Delta())
}

// RunPreparation is called automatically when calling Open. It should only be called once.
func RunPreparation(defaultScene Scene) {
	Time = NewClock()
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

	for {
		select {
		case <-ticker.C:
			RunIteration()
		case <-resetLoopTicker:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))
		case <-c:
			fallthrough
		case <-closeGame:
			ticker.Stop()
			closeEvent()
			return
		}
	}
}

// CursorPos returns (0, 0) because there is no cursor
func CursorPos() (x, y float32) {
	return float32(0), float32(0)
}

// WindowSize gets the current window size
func WindowSize() (w, h int) {
	return int(windowWidth), int(windowHeight)
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

// SetCursor does nothing since there's no headless cursor
func SetCursor(c Cursor) {}

// SetVSync does nothing since there's no monitor to synchronize with
func SetVSync(enabled bool) {}

//SetCursorVisibility does nothing since there's no headless cursor
func SetCursorVisibility(visible bool) {}

// openFile is the desktop-specific way of opening a file
func openFile(url string) (io.ReadCloser, error) {
	return os.Open(url)
}

// IsAndroidChrome tells if the browser is Chrome for android
func IsAndroidChrome() bool {
	return false
}
