// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netgo

package engi

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

func init() {
	rafPolyfill()
}

var canvas js.Object
var gl *webgl.Context

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
	gl, err = webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	GL = newgl2()

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
		responder.Key(Key(ev.Get("keyCode").Int()), 0, PRESS)
	}, false)

	js.Global.Call("addEventListener", "keyup", func(ev js.Object) {
		responder.Key(Key(ev.Get("keyCode").Int()), 0, RELEASE)
	}, false)

	GL.Viewport(0, 0, width, height)

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
	GL.Clear(GL.COLOR_BUFFER_BIT)
	responder.Render()
	Time.Tick()
}

func exit() {
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

type TextureObject struct{ js.Object }
type BufferObject struct{ js.Object }
type FramebufferObject struct{ js.Object }
type ProgramObject struct{ js.Object }
type UniformObject struct{ js.Object }
type ShaderObject struct{ js.Object }

func newgl2() *gl2 {
	return &gl2{
		gl.ELEMENT_ARRAY_BUFFER,
		gl.ARRAY_BUFFER,
		gl.STATIC_DRAW,
		gl.DYNAMIC_DRAW,
		gl.SRC_ALPHA,
		gl.ONE_MINUS_SRC_ALPHA,
		gl.BLEND,
		gl.TEXTURE_2D,
		gl.TEXTURE0,
		gl.UNSIGNED_SHORT,
		gl.UNSIGNED_BYTE,
		gl.FLOAT,
		gl.TRIANGLES,
		gl.LINEAR,
		gl.CLAMP_TO_EDGE,
		gl.FRAMEBUFFER,
		gl.COLOR_ATTACHMENT0,
		gl.FRAMEBUFFER_COMPLETE,
		gl.COLOR_BUFFER_BIT,
		gl.VERTEX_SHADER,
		gl.FRAGMENT_SHADER,
		gl.TEXTURE_WRAP_S,
		gl.TEXTURE_WRAP_T,
		gl.TEXTURE_MIN_FILTER,
		gl.TEXTURE_MAG_FILTER,
		gl.LINEAR_MIPMAP_LINEAR,
		gl.NEAREST,
		gl.RGBA,
	}
}

type gl2 struct {
	ELEMENT_ARRAY_BUFFER int
	ARRAY_BUFFER         int
	STATIC_DRAW          int
	DYNAMIC_DRAW         int
	SRC_ALPHA            int
	ONE_MINUS_SRC_ALPHA  int
	BLEND                int
	TEXTURE_2D           int
	TEXTURE0             int
	UNSIGNED_SHORT       int
	UNSIGNED_BYTE        int
	FLOAT                int
	TRIANGLES            int
	LINEAR               int
	CLAMP_TO_EDGE        int
	FRAMEBUFFER          int
	COLOR_ATTACHMENT0    int
	FRAMEBUFFER_COMPLETE int
	COLOR_BUFFER_BIT     int
	VERTEX_SHADER        int
	FRAGMENT_SHADER      int
	TEXTURE_WRAP_S       int
	TEXTURE_WRAP_T       int
	TEXTURE_MIN_FILTER   int
	TEXTURE_MAG_FILTER   int
	LINEAR_MIPMAP_LINEAR int
	NEAREST              int
	RGBA                 int
}

func (z *gl2) BindTexture(target int, texture *TextureObject) {
	gl.BindTexture(target, texture.Object)
}

func (z *gl2) DeleteTexture(tex *TextureObject) {
	gl.DeleteTexture(tex)
}

func (z *gl2) TexParameteri(target int, param int, value int) {
	gl.TexParameteri(target, param, value)
}

func (z *gl2) UseProgram(program *ProgramObject) {
	if program != nil {
		gl.UseProgram(program.Object)
	}
}

func (z *gl2) GetUniformLocation(program *ProgramObject, uniform string) *UniformObject {
	return &UniformObject{gl.GetUniformLocation(program.Object, uniform)}
}

func (z *gl2) GetAttribLocation(program *ProgramObject, attrib string) int {
	return gl.GetAttribLocation(program.Object, attrib)
}

func (z *gl2) Disable(flag int) {
	gl.Disable(flag)
}

func (z *gl2) BindBuffer(typ int, buf *BufferObject) {
	if buf != nil {
		gl.BindBuffer(typ, buf.Object)
	}
}

func (z *gl2) Enable(flag int) {
	gl.Enable(flag)
}

func (z *gl2) BlendFunc(src, dst int) {
	gl.BlendFunc(src, dst)
}

func (z *gl2) ActiveTexture(flag int) {
	gl.ActiveTexture(flag)
}

func (z *gl2) Uniform2f(uf *UniformObject, x, y float32) {
	gl.Uniform2f(uf.Object, x, y)
}

func (z *gl2) BufferSubData(flag, offset, size int, data interface{}) {
	gl.BufferSubData(flag, offset, data)
}

func (z *gl2) EnableVertexAttribArray(pos int) {
	gl.EnableVertexAttribArray(pos)
}

func (z *gl2) VertexAttribPointer(pos, size, typ int, n bool, stride, offset int) {
	gl.VertexAttribPointer(pos, size, typ, n, stride, offset)
}

func (z *gl2) DrawElements(typ, size, flag, offset int) {
	gl.DrawElements(typ, size, flag, offset)
}

func (z *gl2) CreateBuffer() *BufferObject {
	return &BufferObject{gl.CreateBuffer()}
}

func (z *gl2) BufferData(typ int, data interface{}, flag int) {
	gl.BufferData(typ, data, flag)
}

func (z *gl2) Viewport(x, y, width, height int) {
	gl.Viewport(x, y, width, height)
}

func (z *gl2) BindFramebuffer(typ int, buf *FramebufferObject) {
	gl.BindFramebuffer(typ, buf)
}

func (z *gl2) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (z *gl2) Clear(flags int) {
	gl.Clear(flags)
}

func (z *gl2) CreateFramebuffer() *FramebufferObject {
	return &FramebufferObject{gl.CreateFramebuffer()}
}

func (z *gl2) FramebufferTexture2D(target, attachment, textarget int, texture *TextureObject, level int) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

func (z *gl2) CheckFramebufferStatus(target int) int {
	return gl.CheckFramebufferStatus(target)
}

func (z *gl2) DeleteFramebuffer(buf *FramebufferObject) {
	gl.DeleteFramebuffer(buf)
}

func (z *gl2) CreateShader(flag int) *ShaderObject {
	return &ShaderObject{gl.CreateShader(flag)}
}

func (z *gl2) ShaderSource(shader *ShaderObject, src string) {
	gl.ShaderSource(shader.Object, src)
}

func (z *gl2) CompileShader(shader *ShaderObject) {
	gl.CompileShader(shader.Object)
}

func (z *gl2) DeleteShader(shader *ShaderObject) {
	gl.DeleteShader(shader.Object)
}

func (z *gl2) AttachShader(program *ProgramObject, shader *ShaderObject) {
	gl.AttachShader(program.Object, shader.Object)
}

func (z *gl2) CreateProgram() *ProgramObject {
	return &ProgramObject{gl.CreateProgram()}
}

func (z *gl2) LinkProgram(program *ProgramObject) {
	gl.LinkProgram(program.Object)
}

func (z *gl2) CreateTexture() *TextureObject {
	return &TextureObject{gl.CreateTexture()}
}

func (z *gl2) TexImage2D(target, level, internalFormat, width, height, border, format, kind int, data interface{}) {
	var pix js.Object
	if data == nil {
		pix = nil
	} else {
		pix = data.(js.Object)
	}
	gl.TexImage2D(target, level, internalFormat, format, kind, pix)
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
