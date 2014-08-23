// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !netgo,!android

package engi

import (
	"github.com/go-gl/glow/gl/2.1/gl"
	glfw "github.com/slimsag/glfw3"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
)

var window *glfw.Window

func run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		defer glfw.Terminate()
	}

	if !config.Resizable {
		glfw.WindowHint(glfw.Resizable, 0)
	}
	glfw.WindowHint(glfw.Samples, config.Fsaa)

	width := config.Width
	height := config.Height

	monitor, err := glfw.GetPrimaryMonitor()
	if err != nil {
		log.Fatal(err)
	}
	mode, err := monitor.GetVideoMode()
	if err != nil {
		log.Fatal(err)
	}

	if config.Fullscreen {
		width = mode.Width
		height = mode.Height
		glfw.WindowHint(glfw.Decorated, 0)
	} else {
		monitor = nil
	}

	title := config.Title

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	config.Width, config.Height = window.GetSize()

	if !config.Fullscreen {
		window.SetPosition((mode.Width-width)/2, (mode.Height-height)/2)
	}

	if config.Vsync {
		glfw.SwapInterval(1)
	}

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	GL = newgl2()

	GL.Viewport(0, 0, config.Width, config.Height)

	window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		config.Width, config.Height = window.GetSize()
		responder.Resize(w, h)
	})

	window.SetCursorPositionCallback(func(window *glfw.Window, x, y float64) {
		responder.Mouse(float32(x), float32(y), MOVE)
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPosition()
		if a == glfw.Press {
			responder.Mouse(float32(x), float32(y), PRESS)
		} else {
			responder.Mouse(float32(x), float32(y), RELEASE)
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		responder.Scroll(float32(yoff))
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		if a == glfw.Press {
			responder.Key(Key(k), Modifier(m), PRESS)
		} else {
			responder.Key(Key(k), Modifier(m), RELEASE)
		}
	})

	window.SetCharacterCallback(func(window *glfw.Window, char rune) {
		responder.Type(char)
	})

	responder.Preload()

	timing = NewStats(config.LogFPS)
	timing.Update()

	Files.Load(func() {})

	responder.Setup()

	for !window.ShouldClose() {
		responder.Update(float32(timing.Dt))
		GL.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.SwapBuffers()
		glfw.PollEvents()
		timing.Update()
	}
}

func exit() {
	window.SetShouldClose(true)
}

const (
	Dash         = Key(glfw.KeyMinus)
	Apostrophe   = Key(glfw.KeyApostrophe)
	Semicolon    = Key(glfw.KeySemicolon)
	Equals       = Key(glfw.KeyEqual)
	Comma        = Key(glfw.KeyComma)
	Period       = Key(glfw.KeyPeriod)
	Slash        = Key(glfw.KeySlash)
	Backslash    = Key(glfw.KeyBackslash)
	Backspace    = Key(glfw.KeyBackspace)
	Tab          = Key(glfw.KeyTab)
	CapsLock     = Key(glfw.KeyCapsLock)
	Space        = Key(glfw.KeySpace)
	Enter        = Key(glfw.KeyEnter)
	Escape       = Key(glfw.KeyEscape)
	Insert       = Key(glfw.KeyInsert)
	PrintScreen  = Key(glfw.KeyPrintScreen)
	Delete       = Key(glfw.KeyDelete)
	PageUp       = Key(glfw.KeyPageUp)
	PageDown     = Key(glfw.KeyPageDown)
	Home         = Key(glfw.KeyHome)
	End          = Key(glfw.KeyEnd)
	Pause        = Key(glfw.KeyPause)
	ScrollLock   = Key(glfw.KeyScrollLock)
	ArrowLeft    = Key(glfw.KeyLeft)
	ArrowRight   = Key(glfw.KeyRight)
	ArrowDown    = Key(glfw.KeyDown)
	ArrowUp      = Key(glfw.KeyUp)
	LeftBracket  = Key(glfw.KeyLeftBracket)
	LeftShift    = Key(glfw.KeyLeftShift)
	LeftControl  = Key(glfw.KeyLeftControl)
	LeftSuper    = Key(glfw.KeyLeftSuper)
	LeftAlt      = Key(glfw.KeyLeftAlt)
	RightBracket = Key(glfw.KeyRightBracket)
	RightShift   = Key(glfw.KeyRightShift)
	RightControl = Key(glfw.KeyRightControl)
	RightSuper   = Key(glfw.KeyRightSuper)
	RightAlt     = Key(glfw.KeyRightAlt)
	Zero         = Key(glfw.Key0)
	One          = Key(glfw.Key1)
	Two          = Key(glfw.Key2)
	Three        = Key(glfw.Key3)
	Four         = Key(glfw.Key4)
	Five         = Key(glfw.Key5)
	Six          = Key(glfw.Key6)
	Seven        = Key(glfw.Key7)
	Eight        = Key(glfw.Key8)
	Nine         = Key(glfw.Key9)
	F1           = Key(glfw.KeyF1)
	F2           = Key(glfw.KeyF2)
	F3           = Key(glfw.KeyF3)
	F4           = Key(glfw.KeyF4)
	F5           = Key(glfw.KeyF5)
	F6           = Key(glfw.KeyF6)
	F7           = Key(glfw.KeyF7)
	F8           = Key(glfw.KeyF8)
	F9           = Key(glfw.KeyF9)
	F10          = Key(glfw.KeyF10)
	F11          = Key(glfw.KeyF11)
	F12          = Key(glfw.KeyF12)
	A            = Key(glfw.KeyA)
	B            = Key(glfw.KeyB)
	C            = Key(glfw.KeyC)
	D            = Key(glfw.KeyD)
	E            = Key(glfw.KeyE)
	F            = Key(glfw.KeyF)
	G            = Key(glfw.KeyG)
	H            = Key(glfw.KeyH)
	I            = Key(glfw.KeyI)
	J            = Key(glfw.KeyJ)
	K            = Key(glfw.KeyK)
	L            = Key(glfw.KeyL)
	M            = Key(glfw.KeyM)
	N            = Key(glfw.KeyN)
	O            = Key(glfw.KeyO)
	P            = Key(glfw.KeyP)
	Q            = Key(glfw.KeyQ)
	R            = Key(glfw.KeyR)
	S            = Key(glfw.KeyS)
	T            = Key(glfw.KeyT)
	U            = Key(glfw.KeyU)
	V            = Key(glfw.KeyV)
	W            = Key(glfw.KeyW)
	X            = Key(glfw.KeyX)
	Y            = Key(glfw.KeyY)
	Z            = Key(glfw.KeyZ)
	NumLock      = Key(glfw.KeyNumLock)
	NumMultiply  = Key(glfw.KeyKpMultiply)
	NumDivide    = Key(glfw.KeyKpDivide)
	NumAdd       = Key(glfw.KeyKpAdd)
	NumSubtract  = Key(glfw.KeyKpSubtract)
	NumZero      = Key(glfw.KeyKp0)
	NumOne       = Key(glfw.KeyKp1)
	NumTwo       = Key(glfw.KeyKp2)
	NumThree     = Key(glfw.KeyKp3)
	NumFour      = Key(glfw.KeyKp4)
	NumFive      = Key(glfw.KeyKp5)
	NumSix       = Key(glfw.KeyKp6)
	NumSeven     = Key(glfw.KeyKp7)
	NumEight     = Key(glfw.KeyKp8)
	NumNine      = Key(glfw.KeyKp9)
	NumDecimal   = Key(glfw.KeyKpDecimal)
	NumEnter     = Key(glfw.KeyKpEnter)
)

type TextureObject uint32
type BufferObject uint32
type FramebufferObject uint32
type ProgramObject uint32
type UniformObject uint32
type ShaderObject uint32

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

func (z *gl2) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (z *gl2) Clear(flags int) {
	gl.Clear(uint32(flags))
}

func (z *gl2) CreateBuffer() *BufferObject {
	var loc uint32
	gl.GenBuffers(1, &loc)
	buffer := BufferObject(loc)
	return &buffer
}

func (z *gl2) BindBuffer(typ int, buf *BufferObject) {
	if buf == nil {
		gl.BindBuffer(uint32(typ), 0)
		return
	}
	gl.BindBuffer(uint32(typ), uint32(*buf))
}

func (z *gl2) BufferData(target int, data interface{}, flag int) {
	s := uintptr(reflect.ValueOf(data).Len()) * reflect.TypeOf(data).Elem().Size()
	gl.BufferData(uint32(target), int(s), gl.Ptr(data), uint32(flag))
}

func (z *gl2) Enable(flag int) {
	gl.Enable(uint32(flag))
}

func (z *gl2) Disable(flag int) {
	gl.Disable(uint32(flag))
}

func (z *gl2) BlendFunc(src, dst int) {
	gl.BlendFunc(uint32(src), uint32(dst))
}

func (z *gl2) ActiveTexture(flag int) {
	gl.ActiveTexture(uint32(flag))
}

func (z *gl2) Uniform2f(uf *UniformObject, x, y float32) {
	gl.Uniform2f(int32(*uf), x, y)
}

func (z *gl2) BufferSubData(flag, offset, size int, data interface{}) {
	gl.BufferSubData(uint32(flag), offset, size, gl.Ptr(data))
}

func (z *gl2) EnableVertexAttribArray(pos int) {
	gl.EnableVertexAttribArray(uint32(pos))
}

func (z *gl2) VertexAttribPointer(pos, size, typ int, n bool, stride, offset int) {
	gl.VertexAttribPointer(uint32(pos), int32(size), uint32(typ), n, int32(stride), gl.PtrOffset(offset))
}

func (z *gl2) DrawElements(typ, size, flag, offset int) {
	gl.DrawElements(uint32(typ), int32(size), uint32(flag), gl.PtrOffset(offset))
}

func (z *gl2) CreateFramebuffer() *FramebufferObject {
	var loc uint32
	gl.GenFramebuffers(1, &loc)
	fb := FramebufferObject(loc)
	return &fb
}

func (z *gl2) BindFramebuffer(typ int, buf *FramebufferObject) {
	if buf == nil {
		gl.BindFramebuffer(uint32(typ), 0)
		return
	}
	gl.BindFramebuffer(uint32(typ), uint32(*buf))
}

func (z *gl2) DeleteFramebuffer(buf *FramebufferObject) {
	buffer := uint32(*buf)
	gl.DeleteFramebuffers(1, &buffer)
}

func (z *gl2) FramebufferTexture2D(target, attachment, textarget int, texture *TextureObject, level int) {
	gl.FramebufferTexture2D(uint32(target), uint32(attachment), uint32(textarget), uint32(*texture), int32(level))
}

func (z *gl2) CheckFramebufferStatus(target int) int {
	return int(gl.CheckFramebufferStatus(uint32(target)))
}

func (z *gl2) Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

func (z *gl2) ShaderSource(shader *ShaderObject, src string) {
	source := gl.Str(src + "\x00")
	gl.ShaderSource(uint32(*shader), 1, &source, nil)
}

func (z *gl2) CreateShader(flag int) *ShaderObject {
	shader := ShaderObject(gl.CreateShader(uint32(flag)))
	return &shader
}

func (z *gl2) CompileShader(shader *ShaderObject) {
	gl.CompileShader(uint32(*shader))
}

func (z *gl2) DeleteShader(shader *ShaderObject) {
	gl.DeleteShader(uint32(*shader))
}

func (z *gl2) CreateProgram() *ProgramObject {
	program := ProgramObject(gl.CreateProgram())
	return &program
}

func (z *gl2) AttachShader(program *ProgramObject, shader *ShaderObject) {
	gl.AttachShader(uint32(*program), uint32(*shader))
}

func (z *gl2) LinkProgram(program *ProgramObject) {
	gl.LinkProgram(uint32(*program))
}

func (z *gl2) UseProgram(program *ProgramObject) {
	if program == nil {
		gl.UseProgram(0)
		return
	}
	gl.UseProgram(uint32(*program))
}

func (z *gl2) GetUniformLocation(program *ProgramObject, uniform string) *UniformObject {
	uo := UniformObject(gl.GetUniformLocation(uint32(*program), gl.Str(uniform+"\x00")))
	return &uo
}

func (z *gl2) GetAttribLocation(program *ProgramObject, attrib string) int {
	return int(gl.GetAttribLocation(uint32(*program), gl.Str(attrib+"\x00")))
}

func (z *gl2) CreateTexture() *TextureObject {
	var loc uint32
	gl.GenTextures(1, &loc)
	texture := TextureObject(loc)
	return &texture
}

func (z *gl2) BindTexture(target int, texture *TextureObject) {
	if texture == nil {
		gl.BindTexture(uint32(target), 0)
		return
	}
	gl.BindTexture(uint32(target), uint32(*texture))
}

func (z *gl2) DeleteTexture(tex *TextureObject) {
	texture := uint32(*tex)
	gl.DeleteTextures(1, &texture)
}

func (z *gl2) TexParameteri(target int, param int, value int) {
	gl.TexParameteri(uint32(target), uint32(param), int32(value))
}

func (z *gl2) TexImage2D(target, level, internalFormat, width, height, border, format, kind int, data interface{}) {
	var pix []uint8
	if data == nil {
		pix = nil
	} else {
		pix = data.(*image.NRGBA).Pix
	}
	gl.TexImage2D(uint32(target), int32(level), int32(internalFormat), int32(width), int32(height), int32(border), uint32(format), uint32(kind), gl.Ptr(pix))
}

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

func loadImage(r Resource) (Image, error) {
	file, err := os.Open(r.url)
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

func loadJson(r Resource) (string, error) {
	file, err := ioutil.ReadFile(r.url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

type Assets struct {
	queue  []string
	cache  map[string]Image
	loads  int
	errors int
}

func NewAssets() *Assets {
	return &Assets{make([]string, 0), make(map[string]Image), 0, 0}
}

func (a *Assets) Image(path string) {
	a.queue = append(a.queue, path)
}

func (a *Assets) Get(path string) Image {
	return a.cache[path]
}

func (a *Assets) Load(onFinish func()) {
	if len(a.queue) == 0 {
		onFinish()
	} else {
		for _, path := range a.queue {
			img := LoadImage(path)
			a.cache[path] = img
		}
	}
}

func LoadImage(data interface{}) Image {
	var m image.Image

	switch data := data.(type) {
	default:
		log.Fatal("NewTexture needs a string or io.Reader")
	case string:
		file, err := os.Open(data)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case io.Reader:
		img, _, err := image.Decode(data)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case image.Image:
		m = data
	}

	b := m.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	return &ImageObject{newm}
}
