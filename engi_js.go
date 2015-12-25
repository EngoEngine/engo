// +build netgo

package engi

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/paked/webgl"
)

func init() {
	rafPolyfill()
}

var canvas js.Object

func run(title string, width, height int, fullscreen bool) {
	document := js.Global.Get("document")
	canvas = document.Call("createElement", "canvas")

	target := document.Call("getElementById", title)
	if target.IsNull() {
		target = document.Get("body")
	}
	target.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false
	attrs.Depth = false
	attrs.PremultipliedAlpha = false
	attrs.PreserveDrawingBuffer = false
	attrs.Antialias = false

	var err error
	Gl, err = webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	js.Global.Set("onunload", func() {
		responder.Close()
	})

	canvas.Get("style").Set("display", "block")
	winWidth := js.Global.Get("innerWidth").Int()
	winHeight := js.Global.Get("innerHeight").Int()
	if fullscreen {
		canvas.Set("width", winWidth)
		canvas.Set("height", winHeight)
	} else {
		canvas.Set("width", width)
		canvas.Set("height", height)
		canvas.Get("style").Set("marginLeft", toPx((winWidth-width)/2))
		canvas.Get("style").Set("marginTop", toPx((winHeight-height)/2))
	}

	canvas.Call("addEventListener", "mousemove", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		responder.Mouse(x, y, MOVE)
	}, false)

	canvas.Call("addEventListener", "mousedown", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		responder.Mouse(x, y, PRESS)
	}, false)

	canvas.Call("addEventListener", "mouseup", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		responder.Mouse(x, y, RELEASE)
	}, false)

	canvas.Call("addEventListener", "touchstart", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, PRESS)
		}
	}, false)

	canvas.Call("addEventListener", "touchcancel", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, RELEASE)
		}
	}, false)

	canvas.Call("addEventListener", "touchend", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, PRESS)
		}
	}, false)

	canvas.Call("addEventListener", "touchmove", func(ev js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, MOVE)
		}
	}, false)

	js.Global.Call("addEventListener", "keypress", func(ev js.Object) {
		responder.Type(rune(ev.Get("charCode").Int()))
	}, false)

	js.Global.Call("addEventListener", "keydown", func(ev js.Object) {
		key := Key(ev.Get("keyCode").Int())
		keyStates[key] = true
	}, false)

	js.Global.Call("addEventListener", "keyup", func(ev js.Object) {
		key := Key(ev.Get("keyCode").Int())
		keyStates[key] = false
		// responder.Key(Key(ev.Get("keyCode").Int()), 0, RELEASE)
	}, false)

	Gl.Viewport(0, 0, width, height)
	Wo.New()
	responder.Preload()
	Files.Load(func() {
		responder.Setup()
		RequestAnimationFrame(animate)
	})
}

func width() float32 {
	return float32(canvas.Get("width").Int())
}

func height() float32 {
	return float32(canvas.Get("height").Int())
}

func animate(dt float32) {
	RequestAnimationFrame(animate)
	responder.Update(Time.Delta())
	Time.Tick()
	keysUpdate()
}

func exit() {
	responder.Close()
}

func toPx(n int) string {
	return strconv.FormatInt(int64(n), 10) + "px"
}

func rafPolyfill() {
	window := js.Global
	vendors := []string{"ms", "moz", "webkit", "o"}
	if window.Get("requestAnimationFrame").IsUndefined() {
		for i := 0; i < len(vendors) && window.Get("requestAnimationFrame").IsUndefined(); i++ {
			vendor := vendors[i]
			window.Set("requestAnimationFrame", window.Get(vendor+"RequestAnimationFrame"))
			window.Set("cancelAnimationFrame", window.Get(vendor+"CancelAnimationFrame"))
			if window.Get("cancelAnimationFrame").IsUndefined() {
				window.Set("cancelAnimationFrame", window.Get(vendor+"CancelRequestAnimationFrame"))
			}
		}
	}

	lastTime := 0.0
	if window.Get("requestAnimationFrame").IsUndefined() {
		window.Set("requestAnimationFrame", func(callback func(float32)) int {
			currTime := js.Global.Get("Date").New().Call("getTime").Float()
			timeToCall := math.Max(0, 16-(currTime-lastTime))
			id := window.Call("setTimeout", func() { callback(float32(currTime + timeToCall)) }, timeToCall)
			lastTime = currTime + timeToCall
			return id.Int()
		})
	}

	if window.Get("cancelAnimationFrame").IsUndefined() {
		window.Set("cancelAnimationFrame", func(id int) {
			js.Global.Get("clearTimeout").Invoke(id)
		})
	}
}

func RequestAnimationFrame(callback func(float32)) int {
	return js.Global.Call("requestAnimationFrame", callback).Int()
}

func CancelAnimationFrame(id int) {
	js.Global.Call("cancelAnimationFrame")
}

func loadImage(r Resource) (Image, error) {
	ch := make(chan error, 1)

	img := js.Global.Get("Image").New()
	img.Call("addEventListener", "load", func(js.Object) {
		go func() { ch <- nil }()
	}, false)
	img.Call("addEventListener", "error", func(o js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	img.Set("src", r.url+"?"+strconv.FormatInt(rand.Int63(), 10))

	err := <-ch
	if err != nil {
		return nil, err
	}

	return &ImageObject{img}, nil
}

func loadJson(r Resource) (string, error) {
	ch := make(chan error, 1)

	req := js.Global.Get("XMLHttpRequest").New()
	req.Call("open", "GET", r.url, true)
	req.Call("addEventListener", "load", func(js.Object) {
		go func() { ch <- nil }()
	}, false)
	req.Call("addEventListener", "error", func(o js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	req.Call("send", nil)

	err := <-ch
	if err != nil {
		return "", err
	}

	return req.Get("responseText").Str(), nil
}

type ImageObject struct {
	data js.Object
}

func (i *ImageObject) Data() interface{} {
	return i.data
}

func (i *ImageObject) Width() int {
	return i.data.Get("width").Int()
}

func (i *ImageObject) Height() int {
	return i.data.Get("height").Int()
}
