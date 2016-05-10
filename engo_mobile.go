//+build android

package engo

import (
	"engo.io/gl"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
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

	msaaPreference int
)

func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	gameWidth = float32(width)
	gameHeight = float32(height)
	msaaPreference = msaa
}

func loadJSON(r Resource) (string, error) {
	return "", fmt.Errorf("loadJSON not yet impplemented")
}

func loadFont(r Resource) (*truetype.Font, error) {
	if strings.HasPrefix(r.URL, "assets/") {
		r.URL = r.URL[7:]
	}

	file, err := asset.Open(r.URL)
	if err != nil {
		return nil, err
	}

	ttfBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(ttfBytes)
}

func WindowSize() (w, h int) {
	return sz.WidthPx, sz.HeightPx
}

func CursorPos() (x, y float64) {
	log.Println("warning: not yet implemented CursorPos")
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
				if Mouse.Action == RELEASE {
					Mouse.Action = NEUTRAL
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
					Mouse.Action = PRESS
				case touch.TypeMove:
					Mouse.Action = MOVE
				case touch.TypeEnd:
					Mouse.Action = RELEASE
				}
			}
		}
	})
}

// RunPreparation is called only once, and is called automatically when calling Open
// It is only here for benchmarking in combination with OpenHeadlessNoRun
func RunPreparation(defaultScene Scene) {
	Time = NewClock()
	Files = NewLoader()

	// Default WorldBounds values
	WorldBounds.Max = Point{GameWidth(), GameHeight()}

	SetScene(defaultScene, false)
}

// RunIteration runs one iteration / frame
func RunIteration() {
	if !headless {
		Input.update()
	}

	// Then update the world and all Systems
	currentWorld.Update(Time.Delta())

	Time.Tick()
}

// SetCursor changes the cursor - not yet implemented
func SetCursor(c Cursor) {
	log.Println("SetCursor: not yet implemented for mobile")
}
