//+build android

package engo

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"engo.io/gl"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	mobilegl "golang.org/x/mobile/gl"
)

var (
	Gl *gl.Context
	sz size.Event

	gameWidth, gameHeight     float32
	windowWidth, windowHeight float32
	canvasWidth, canvasHeight float32

	msaaPreference int
)

func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	gameWidth = float32(width)
	gameHeight = float32(height)
	msaaPreference = msaa
}

func WindowSize() (w, h int) {
	return sz.WidthPx, sz.HeightPx
}

func CursorPos() (x, y float64) {
	notImplemented("CursorPos")
	return 0, 0
}

func GameWidth() float32 {
	return gameWidth
}

func GameHeight() float32 {
	return gameHeight
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
			images *glutil.Images
			fps    *debug.FPS
		)

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					Gl = gl.NewContext(e.DrawContext)
					RunPreparation(defaultScene)

					images = glutil.NewImages(e.DrawContext.(mobilegl.Context))
					fps = debug.NewFPS(images)

					// Let the device know we want to start painting :-)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					closeEvent()
				}

			case size.Event:
				sz = e
				windowWidth = float32(sz.WidthPx)
				windowHeight = float32(sz.HeightPx)
				canvasWidth = float32(sz.WidthPx)
				canvasHeight = float32(sz.HeightPx)
				Gl.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
			case paint.Event:
				if e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				RunIteration()
				if closeGame {
					break
				}

				fps.Draw(sz)

				// Reset mouse if needed
				if Mouse.Action == Release {
					Mouse.Action = Neutral
				}

				a.Publish() // same as SwapBuffers

				// Drive the animation by preparing to paint the next frame
				// after this one is shown. - FPS is ignored here!
				a.Send(paint.Event{})
			case touch.Event:
				Mouse.X = e.X
				Mouse.Y = e.Y
				switch e.Type {
				case touch.TypeBegin:
					Mouse.Action = Press
				case touch.TypeMove:
					Mouse.Action = Move
				case touch.TypeEnd:
					Mouse.Action = Release
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
	if !opts.HeadlessMode {
		Input.update()
	}

	// Then update the world and all Systems
	currentWorld.Update(Time.Delta())

	Time.Tick()
}

// SetCursor changes the cursor - not yet implemented
func SetCursor(Cursor) {
	notImplemented("SetCursor")
}

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
