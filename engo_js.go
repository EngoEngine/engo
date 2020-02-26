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
	"syscall/js"
	"time"

	"github.com/EngoEngine/gl"
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

	canvas.Call("addEventListener", "keypress", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// TODO: Not sure what to do here, come back
		//ke := ev.(*dom.KeyboardEvent)
		//responser.Type(rune(keyStates[Key(ke.KeyCode)]))
		return nil
	}))

	canvas.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
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
			return nil
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
		return nil
	}))

	canvas.Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
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
			return nil
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
		return nil
	}))

	canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
		return nil
	}))

	canvas.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
		Input.Mouse.Action = Press
		return nil
	}))

	canvas.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		mmX, mmY := event.Get("clientX").Int(), event.Get("clientY").Int()
		Input.Mouse.X = float32(mmX) / opts.GlobalScale.X
		Input.Mouse.Y = float32(mmY) / opts.GlobalScale.Y
		Input.Mouse.Action = Release
		return nil
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
		window.Set("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			currTime := js.Global().Get("Date").New().Call("getTime").Float()
			timeToCall := math.Max(0, 16-(currTime-lastTime))
			window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				args[0].Invoke(currTime + timeToCall)
				return nil
			}), timeToCall)
			lastTime = currTime + timeToCall
			return nil
		}))
	}

	if window.Get("cancelAnimationFrame").Type() == js.TypeUndefined {
		window.Set("cancelAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			js.Global().Get("clearTimeout").Invoke(args[0])
			return nil
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
		window.Call("addEventListener", "onbeforeunload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			window.Call("alert", "You're closing")
			return nil
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
	var err error
	var resp js.Value
	ch := make(chan struct{})
	req := js.Global().Get("XMLHttpRequest").New()
	req.Call("open", "GET", url, true)
	req.Set("responseType", "arraybuffer")
	loadf := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(ch)
		status := req.Get("status").Int()
		if 200 <= status && status < 400 {
			resp = req.Get("response")
			return nil
		}
		err = errors.New(fmt.Sprintf("http error: %d", status))
		return nil
	})
	defer loadf.Release()
	req.Call("addEventListener", "load", loadf)
	errorf := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(ch)
		err = errors.New(fmt.Sprintf("XMLHttpRequest error: %s", req.Get("statusText").String()))
		return nil
	})
	req.Call("addEventListener", "error", errorf)
	defer errorf.Release()
	req.Call("send")
	t := time.NewTicker(time.Duration(int(time.Second) * 10))
	select {
	case <-ch:
		if err != nil {
			return nil, err
		}
		buf := make([]byte, resp.Get("byteLength").Int())
		js.CopyBytesToGo(buf, js.Global().Get("Uint8Array").New(resp))
		f := &noCloseReadCloser{bytes.NewReader(buf)}
		return f, nil
	case <-t.C:
		return nil, errors.New("Did not recieve a response in 10 seconds.")
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

// GetKeyName returns the string returned from the given Key. So you can write "Press W to move forward"
// and get a W for QWERTY and Z for AZERTY
func GetKeyName(k Key) string {
	if 96 <= k && k <= 105 {
		k = k - 48
	}
	return js.Global().Get("String").Call("fromCharCode", k).String()
}
