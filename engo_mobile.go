//+build android darwin,arm darwin,arm64 ios
//+build !mobilebind

package engo

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"engo.io/gl"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
)

var (
	// Gl is the current OpenGL context
	Gl *gl.Context
	sz size.Event

	msaaPreference int
)

// CreateWindow creates a window with the specified parameters
func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	CurrentBackEnd = BackEndMobile
	gameWidth = float32(width)
	gameHeight = float32(height)
	msaaPreference = msaa
}

// WindowSize returns the width and height of the current window
func WindowSize() (w, h int) {
	return sz.WidthPx, sz.HeightPx
}

// CursorPos returns the current cursor position
func CursorPos() (x, y float32) {
	notImplemented("CursorPos")
	return 0, 0
}

// WindowWidth returns the current window width
func WindowWidth() float32 {
	return windowWidth
}

// WindowHeight returns the current window height
func WindowHeight() float32 {
	return windowHeight
}

// CanvasWidth returns the current canvas width
func CanvasWidth() float32 {
	return canvasWidth
}

// CanvasHeight returns the current canvas height
func CanvasHeight() float32 {
	return canvasHeight
}

// CanvasScale returns the current scale of the canvas from the original window
func CanvasScale() float32 {
	return CanvasWidth() / WindowWidth()
}

// DestroyWindow handles destroying the window
func DestroyWindow() { /* nothing to do here? */ }

func runLoop(defaultScene Scene, headless bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		closeEvent()
	}()

	app.Main(func(a app.App) {
		var (
			ticker *time.Ticker
		)

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					Gl = gl.NewContext(e.DrawContext)
					RunPreparation(defaultScene)

					ticker = time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))
					// Start tick, minimize the delta
					Time.Tick()

					// Let the device know we want to start painting :-)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					closeEvent()
					ticker.Stop()
					Gl = nil
				}

			case size.Event:
				sz = e
				windowWidth = float32(sz.WidthPx)
				windowHeight = float32(sz.HeightPx)
				canvasWidth = float32(sz.WidthPx)
				canvasHeight = float32(sz.HeightPx)
				Gl.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				ResizeXOffset = (gameWidth - canvasWidth)
				ResizeYOffset = (gameHeight - canvasHeight)
			case paint.Event:
				if e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				select {
				case <-ticker.C:
					RunIteration()
				case <-resetLoopTicker:
					ticker.Stop()
					ticker = time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))
				}

				Input.Mouse.Action = Neutral
				a.Publish() // same as SwapBuffers

				// Drive the animation by preparing to paint the next frame
				// after this one is shown. - FPS is ignored here!
				a.Send(paint.Event{})
			case touch.Event:
				Input.Mouse.X = e.X / opts.GlobalScale.X
				Input.Mouse.Y = e.Y / opts.GlobalScale.Y
				id := int(e.Sequence)
				switch e.Type {
				case touch.TypeBegin:
					Input.Mouse.Action = Press
					Input.Touches[id] = Point{
						X: float32(e.X) / opts.GlobalScale.X,
						Y: float32(e.Y) / opts.GlobalScale.Y,
					}
				case touch.TypeMove:
					Input.Mouse.Action = Move
					Input.Touches[id] = Point{
						X: float32(e.X) / opts.GlobalScale.X,
						Y: float32(e.Y) / opts.GlobalScale.Y,
					}
				case touch.TypeEnd:
					Input.Mouse.Action = Release
					delete(Input.Touches, id)
				}
			}
		}
	})
}

// RunPreparation is called only once, and is called automatically when calling Open
// It is only here for benchmarking in combination with OpenHeadlessNoRun
func RunPreparation(defaultScene Scene) {
	Time = NewClock()
	SetScene(defaultScene, false)
}

// RunIteration runs one iteration / frame
func RunIteration() {
	Time.Tick()

	if !opts.HeadlessMode {
		Input.update()
	}

	// Then update the world and all Systems
	currentUpdater.Update(Time.Delta())
}

// SetCursor changes the cursor - not yet implemented
func SetCursor(Cursor) {
	notImplemented("SetCursor")
}

//SetCursorVisibility sets the visibility of the cursor.
//If true the cursor is visible, if false the cursor is not.
//Does nothing in mobile since there's no visible cursor to begin with
func SetCursorVisibility(visible bool) {}

// SetTitle has no effect on mobile
func SetTitle(title string) {}

// openFile is the mobile-specific way of opening a file
func openFile(url string) (io.ReadCloser, error) {
	usedUrl := url
	if strings.HasPrefix(url, "assets/") {
		usedUrl = usedUrl[7:]
	}

	return asset.Open(usedUrl)
}

// IsAndroidChrome tells if the browser is Chrome for android
func IsAndroidChrome() bool {
	return false
}
