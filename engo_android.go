//+build android

package engo

import (
	"engo.io/gl"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
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
	"io/ioutil"
)

var (
	Gl *gl.Context

	gameWidth, gameHeight     float32
	windowWidth, windowHeight float32
)

func CreateWindow(title string, width, height int, fullscreen bool) {}

func loadImage(r Resource) (Image, error) {
	if strings.HasPrefix(r.url, "assets/") {
		r.url = r.url[7:]
	}

	file, err := asset.Open(r.url)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return &ImageObject{newm}, nil
}

func loadJSON(r Resource) (string, error) {
	return "", fmt.Errorf("loadJSON not yet impplemented")
}

func loadFont(r Resource) (*truetype.Font, error) {
	if strings.HasPrefix(r.url, "assets/") {
		r.url = r.url[7:]
	}

	file, err := asset.Open(r.url)
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
	log.Println("warning: not yet implemented WindowSize")
	return 0, 0
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

func NewImageObject(img *image.NRGBA) *ImageObject {
	return &ImageObject{img}
}

type ImageObject struct {
	data *image.NRGBA
}

func (i *ImageObject) Data() interface{} {
	return i.data
}

func (i *ImageObject) Width() int {
	return i.data.Rect.Max.X
}

func (i *ImageObject) Height() int {
	return i.data.Rect.Max.Y
}

func runLoop(defaultScene Scene, headless bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		closeEvent()
	}()

	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					Gl = gl.NewContext(e.DrawContext)

					RunPreparation(defaultScene)

					// Let the device know we want to start painting :-)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					closeEvent()
				}
			case size.Event:
			//sz = e
			//touchX = float32(sz.WidthPx / 2)
			//touchY = float32(sz.HeightPx / 2)
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

				a.Publish() // same as SwapBuffers

				// Drive the animation by preparing to paint the next frame
				// after this one is shown. - FPS is ignored here!
				a.Send(paint.Event{})
			case touch.Event:
				//touchX = e.X
				//touchY = e.Y
			}
		}
	})
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

func Exit() {
	closeGame = true
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
	// First check for new keypresses
	if !headless {
		Input.update()
		//glfw.PollEvents()
	}

	// Then update the world and all Systems
	currentWorld.Update(Time.Delta())

	// Lastly, forget keypresses and swap buffers
	if !headless {
		// reset values to avoid catching the same "signal" twice
		Mouse.ScrollX, Mouse.ScrollY = 0, 0
		Mouse.Action = NEUTRAL

		//window.SwapBuffers()
	}

	Time.Tick()
}

func SetBackground(c color.Color) {
	if !headless {
		r, g, b, a := c.RGBA()

		Gl.ClearColor(float32(r), float32(g), float32(b), float32(a))
	}
}
