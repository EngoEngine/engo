//+build js

package engo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"engo.io/gl"
	"github.com/gopherjs/gopherwasm/js"
)

var (
	// Gl is the current OpenGL context
	Gl *gl.Context

	devicePixelRatio float64

	poll     = make(map[int]bool)
	pollLock sync.Mutex

	document = js.Global().Get("document")
	window   = js.Global().Get("window")
	canvas   js.Value
)

// CreateWindow creates a window with the specified parameters
func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	rafPolyfill()
	CurrentBackEnd = BackEndWeb
	canvas = document.Call("createElement", "canvas")

	devicePixelRatio = js.Global().Get("devicePixelRatio").Float()
	canvas.Set("width", int(float64(width)+0.5))   // Nearest non-negative int.
	canvas.Set("height", int(float64(height)+0.5)) // Nearest non-negative int.

	if document.Get("body") == js.Null() {
		document.Set("body", document.Call("createElement", "body"))
	}
	body := document.Get("body")

	body.Get("style").Set("margin", "0")
	body.Get("style").Set("padding", "0")

	canvas.Call("setAttribute", "tabindex", 1)
	canvas.Get("style").Set("outline", "none")
	body.Call("appendChild", canvas)

	document.Set("title", title)

	Gl, _ = gl.NewContext(canvas, nil)

	Gl.GetExtension("OES_texture_float")

	gameWidth = float32(width)
	gameHeight = float32(height)
	windowWidth = WindowWidth()
	windowHeight = WindowHeight()

	ResizeXOffset = gameWidth - CanvasWidth()
	ResizeYOffset = gameHeight - CanvasHeight()

	canvas.Call("addEventListener", "keypress", js.NewEventCallback(0, func(event js.Value) {
		// TODO: Not sure what to do here, come back
		//ke := ev.(*dom.KeyboardEvent)
		//responser.Type(rune(keyStates[Key(ke.KeyCode)]))
	}))

	canvas.Call("addEventListener", "keydown", js.NewEventCallback(0, func(event js.Value) {
		ke := event.Get("code") // TODO: preventdefault if it's a special key so it doesn't scroll when you press it!
		if ke.Type() == js.TypeUndefined {
			kc := event.Get("keyCode").Int()
			go func(i int) {
				pollLock.Lock()
				poll[i] = true
				pollLock.Unlock()
			}(kc)
			k := Key(kc)
			if k == KeyArrowUp || k == KeyArrowDown || k == KeyArrowLeft || k == KeyArrowRight || k == KeyTab || k == KeyBackspace || k == KeySpace {
				event.Call("preventDefault")
			}
			return
		}
		k := jsStrToKey[ke.String()]
		go func(i int) {
			pollLock.Lock()
			poll[i] = true
			pollLock.Unlock()
		}(int(k))
		if k == KeyArrowUp || k == KeyArrowDown || k == KeyArrowLeft || k == KeyArrowRight || k == KeyTab || k == KeyBackspace || k == KeySpace {
			event.Call("preventDefault")
		}
		char := event.Get("key").String()
		if len(char) == 1 {
			Mailbox.Dispatch(TextMessage{[]rune(char)[0]})
		}
	}))

	canvas.Call("addEventListener", "keyup", js.NewEventCallback(0, func(event js.Value) {
		ke := event.Get("code")
		if ke.Type() == js.TypeUndefined {
			kc := event.Get("keyCode").Int()
			go func(i int) {
				pollLock.Lock()
				poll[i] = false
				pollLock.Unlock()
			}(kc)
			k := Key(kc)
			if k == KeyArrowUp || k == KeyArrowDown || k == KeyArrowLeft || k == KeyArrowRight || k == KeyTab || k == KeyBackspace || k == KeySpace {
				event.Call("preventDefault")
			}
			return
		}
		k := jsStrToKey[ke.String()]
		go func(i int) {
			pollLock.Lock()
			poll[i] = false
			pollLock.Unlock()
		}(int(k))
		if k == KeyArrowUp || k == KeyArrowDown || k == KeyArrowLeft || k == KeyArrowRight || k == KeyTab || k == KeyBackspace || k == KeySpace {
			event.Call("preventDefault")
		}
	}))

	canvas.Call("addEventListener", "mousemove", js.NewEventCallback(0, func(event js.Value) {
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
	}))

	canvas.Call("addEventListener", "mousedown", js.NewEventCallback(0, func(event js.Value) {
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
		Input.Mouse.Action = Press
	}))

	canvas.Call("addEventListener", "mouseup", js.NewEventCallback(0, func(event js.Value) {
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
		Input.Mouse.Action = Release
	}))
}

// DestroyWindow handles destroying the window when done
func DestroyWindow() {}

// CursorPos returns the current cursor position
func CursorPos() (x, y float32) {
	return Input.Mouse.X * opts.GlobalScale.X, Input.Mouse.Y * opts.GlobalScale.Y
}

// SetTitle changes the title of the page to the given string
func SetTitle(title string) {
	if opts.HeadlessMode {
		log.Println("Title set to:", title)
	} else {
		document.Set("title", title)
	}
}

// WindowSize returns the width and height of the current window
func WindowSize() (w, h int) {
	w = int(WindowWidth())
	h = int(WindowHeight())
	return
}

// WindowWidth returns the current window width
func WindowWidth() float32 {
	return float32(window.Get("innerWidth").Int())
}

// WindowHeight returns the current window height
func WindowHeight() float32 {
	return float32(window.Get("innerHeight").Int())
}

// CanvasWidth returns the current canvas width
func CanvasWidth() float32 {
	return float32(canvas.Get("width").Int())
}

// CanvasHeight returns the current canvas height
func CanvasHeight() float32 {
	return float32(canvas.Get("height").Int())
}

func CanvasScale() float32 {
	return 1
}

func rafPolyfill() {
	vendors := []string{"ms", "moz", "webkit", "o"}
	if window.Get("requestAnimationFrame").Type() == js.TypeUndefined {
		for i := 0; i < len(vendors) && window.Get("requestAnimationFrame").Type() == js.TypeUndefined; i++ {
			vendor := vendors[i]
			window.Set("requestAnimationFrame", window.Get(vendor+"RequestAnimationFrame"))
			window.Set("cancelAnimationFrame", window.Get(vendor+"CancelAnimationFrame"))
			if window.Get("cancelAnimationFrame").Type() == js.TypeUndefined {
				window.Set("cancelAnimationFrame", window.Get(vendor+"CancelRequestAnimationFrame"))
			}
		}
	}

	lastTime := 0.0
	if window.Get("requestAnimationFrame").Type() == js.TypeUndefined {
		window.Set("requestAnimationFrame", js.NewCallback(func(arg1 []js.Value) {
			currTime := js.Global().Get("Date").New().Call("getTime").Float()
			timeToCall := math.Max(0, 16-(currTime-lastTime))
			window.Call("setTimeout", js.NewCallback(func(arg2 []js.Value) {
				arg1[0].Invoke(currTime + timeToCall)
			}), timeToCall)
			lastTime = currTime + timeToCall
		}))
	}

	if window.Get("cancelAnimationFrame").Type() == js.TypeUndefined {
		window.Set("cancelAnimationFrame", js.NewCallback(func(arg1 []js.Value) {
			js.Global().Get("clearTimeout").Invoke(arg1[0])
		}))
	}
}

// RunIteration runs one iteration per frame
func RunIteration() {
	Time.Tick()
	Input.update()
	jsPollKeys()
	currentUpdater.Update(Time.Delta())
	Input.Mouse.Action = Neutral
	// TODO: this may not work, and sky-rocket the FPS
	//  requestAnimationFrame(func(dt float32) {
	// 	currentWorld.Update(Time.Delta())
	// 	keysUpdate()
	// 	if !headless {
	// 		// TODO: does this require !headless?
	// 		Mouse.ScrollX, Mouse.ScrollY = 0, 0
	// 	}
	// 	Time.Tick()
	// })
}

// jsPollKeys polls the keys collected by the javascript callback
// this ensures the keys only get updated once per frame, since the
// callback has no information about the frames and is invoked several
// times between frames. This makes Input.Button.JustPressed and JustReleased
// able to return true properly.
func jsPollKeys() {
	pollLock.Lock()
	defer pollLock.Unlock()

	for key, state := range poll {
		Input.keys.Set(Key(key), state)
		delete(poll, key)
	}
}

func requestAnimationFrame(callback func(float32)) int {
	//return dom.GetWindow().RequestAnimationFrame(callback)
	return js.Global().Call("requestAnimationFrame", callback).Int()
}

func cancelAnimationFrame(id int) {
	window.Call("cancelAnimationFrame", id)
}

// RunPreparation is called automatically when calling Open. It should only be called once.
func RunPreparation() {
	Time = NewClock()

	if !opts.HeadlessMode {
		window.Call("addEventListener", "onbeforeunload", js.NewEventCallback(js.PreventDefault, func(event js.Value) {
			window.Call("alert", "You're closing")
		}))
	}
}

func runLoop(defaultScene Scene, headless bool) {
	SetScene(defaultScene, false)
	RunPreparation()
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
		case <-closeGame:
			ticker.Stop()
			closeEvent()
			return
		}
	}
}

func openFile(url string) (io.ReadCloser, error) {
	if Headless() { // Headless would be node.js
		return os.Open(url)
	}
	var resp js.Value
	var err error
	loadC := make(chan struct{})
	req := js.Global().Get("XMLHttpRequest").New()
	req.Call("open", "GET", url, true)
	req.Set("responseType", "arraybuffer")
	req.Call("addEventListener", "load", js.NewCallback(func([]js.Value) {
		status := req.Get("status").Int()
		if 200 <= status && status < 400 {
			resp = req.Get("response")
			loadC <- struct{}{}
			return
		}
		err = errors.New(fmt.Sprintf("http error: %d", status))
	}))
	req.Call("send")
	select {
	case <-loadC:
		if err != nil {
			return nil, err
		}
		uint8contentWrapper := js.Global().Get("Uint8Array").New(resp)
		data := make([]byte, uint8contentWrapper.Get("byteLength").Int())
		arr := js.TypedArrayOf(data)
		arr.Call("set", uint8contentWrapper)
		arr.Release()
		return noCloseReadCloser{bytes.NewReader(data)}, nil
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout while trying to fetch resource: " + url)
	}
}

type noCloseReadCloser struct {
	r io.Reader
}

func (n noCloseReadCloser) Close() error { return nil }
func (n noCloseReadCloser) Read(p []byte) (int, error) {
	return n.r.Read(p)
}

// SetCursor changes the cursor
func SetCursor(c Cursor) {
	switch c {
	case CursorNone:
		document.Get("body").Get("style").Set("cursor", "default")
	case CursorHand:
		document.Get("body").Get("style").Set("cursor", "hand")
	}
}

//SetCursorVisibility sets the visibility of the cursor.
//If true the cursor is visible, if false the cursor is not.
func SetCursorVisibility(visible bool) {
	if visible {
		document.Get("body").Get("style").Set("cursor", "default")
	} else {
		document.Get("body").Get("style").Set("cursor", "none")
	}
}

// IsAndroidChrome tells if the browser is Chrome for android
func IsAndroidChrome() bool {
	ua := js.Global().Get("navigator").Get("userAgent").String()
	if !strings.Contains(ua, "Android") {
		return false
	}
	if !strings.Contains(ua, "Chrome") {
		return false
	}
	return true
}
