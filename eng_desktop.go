// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !netgo,!android

package eng

import (
	"github.com/errcw/glow/gl/2.1/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
)

var window *glfw.Window

func run() {
	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		log.Fatal("GLFW error %v: %v\n", err, desc)
	})

	if ok := glfw.Init(); ok {
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

	window.SetCharacterCallback(func(window *glfw.Window, char uint) {
		responder.Type(rune(char))
	})

	responder.Setup()

	timing = NewStats(config.LogFPS)
	timing.Update()

	for !window.ShouldClose() {
		responder.Update(float32(timing.Dt))
		GL.ClearColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
		GL.Clear(gl.COLOR_BUFFER_BIT)
		responder.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
		timing.Update()
	}
}

func exit() {
	window.SetShouldClose(true)
}

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
		gl.GENERATE_MIPMAP,
		gl.NEAREST,
		gl.TRUE,
		gl.FALSE,
		gl.RGBA,
	}
}

type gl2 struct {
	ELEMENT_ARRAY_BUFFER uint32
	ARRAY_BUFFER         uint32
	STATIC_DRAW          uint32
	DYNAMIC_DRAW         uint32
	SRC_ALPHA            uint32
	ONE_MINUS_SRC_ALPHA  uint32
	BLEND                uint32
	TEXTURE_2D           uint32
	TEXTURE0             uint32
	UNSIGNED_SHORT       uint32
	UNSIGNED_BYTE        uint32
	FLOAT                uint32
	TRIANGLES            uint32
	LINEAR               int32
	CLAMP_TO_EDGE        int32
	FRAMEBUFFER          uint32
	COLOR_ATTACHMENT0    uint32
	FRAMEBUFFER_COMPLETE uint32
	COLOR_BUFFER_BIT     uint32
	VERTEX_SHADER        uint32
	FRAGMENT_SHADER      uint32
	TEXTURE_WRAP_S       uint32
	TEXTURE_WRAP_T       uint32
	TEXTURE_MIN_FILTER   uint32
	TEXTURE_MAG_FILTER   uint32
	LINEAR_MIPMAP_LINEAR int32
	GENERATE_MIPMAP      uint32
	NEAREST              int32
	TRUE                 int32
	FALSE                int32
	RGBA                 uint32
}

func (z *gl2) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (z *gl2) Clear(flags uint32) {
	gl.Clear(flags)
}

func (z *gl2) GenBuffers(num int32, loc *uint32) {
	gl.GenBuffers(num, loc)
}

func (z *gl2) BindBuffer(typ, buf uint32) {
	gl.BindBuffer(typ, buf)
}

func (z *gl2) BufferData(typ uint32, size int, data interface{}, flag uint32) {
	gl.BufferData(typ, size, gl.Ptr(data), flag)
}

func (z *gl2) Enable(flag uint32) {
	gl.Enable(flag)
}

func (z *gl2) Disable(flag uint32) {
	gl.Disable(flag)
}

func (z *gl2) BlendFunc(src, dst uint32) {
	gl.BlendFunc(src, dst)
}

func (z *gl2) ActiveTexture(flag uint32) {
	gl.ActiveTexture(flag)
}

func (z *gl2) Uniform2f(uf int32, x, y float32) {
	gl.Uniform2f(uf, x, y)
}

func (z *gl2) BufferSubData(flag uint32, offset int, size int, data interface{}) {
	gl.BufferSubData(flag, offset, size, gl.Ptr(data))
}

func (z *gl2) EnableVertexAttribArray(pos uint32) {
	gl.EnableVertexAttribArray(pos)
}

func (z *gl2) VertexAttribPointer(pos uint32, size int32, typ uint32, b bool, stride int32, offset int) {
	gl.VertexAttribPointer(pos, size, typ, b, stride, gl.PtrOffset(offset))
}

func (z *gl2) DrawElements(typ uint32, size int32, flag uint32, offset int) {
	gl.DrawElements(typ, size, flag, gl.PtrOffset(offset))
}

func (z *gl2) GenFramebuffers(num int32, loc *uint32) {
	gl.GenFramebuffers(num, loc)
}

func (z *gl2) BindFramebuffer(typ, buf uint32) {
	gl.BindFramebuffer(typ, buf)
}

func (z *gl2) DeleteFramebuffers(num int32, loc *uint32) {
	gl.DeleteFramebuffers(num, loc)
}

func (z *gl2) FramebufferTexture2D(target, attachment, textarget, texture uint32, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

func (z *gl2) CheckFramebufferStatus(target uint32) uint32 {
	return gl.CheckFramebufferStatus(target)
}

func (z *gl2) Viewport(x, y, width, height int32) {
	gl.Viewport(x, y, width, height)
}

func (z *gl2) ShaderSource(shader uint32, num int32, src string, length *int32) {
	source := gl.Str(src + "\x00")
	gl.ShaderSource(shader, num, &source, length)
}

func (z *gl2) CreateShader(flag uint32) uint32 {
	return gl.CreateShader(flag)
}

func (z *gl2) CompileShader(shader uint32) {
	gl.CompileShader(shader)
}

func (z *gl2) DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

func (z *gl2) CreateProgram() uint32 {
	return gl.CreateProgram()
}

func (z *gl2) AttachShader(program uint32, shader uint32) {
	gl.AttachShader(program, shader)
}

func (z *gl2) LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

func (z *gl2) UseProgram(program uint32) {
	gl.UseProgram(program)
}

func (z *gl2) GetUniformLocation(shader uint32, uniform string) int32 {
	return gl.GetUniformLocation(shader, gl.Str(uniform+"\x00"))
}

func (z *gl2) GetAttribLocation(shader uint32, attrib string) int32 {
	return gl.GetAttribLocation(shader, gl.Str(attrib+"\x00"))
}

func (z *gl2) GenTextures(num int32, loc *uint32) {
	gl.GenTextures(num, loc)
}

func (z *gl2) BindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
}

func (z *gl2) DeleteTextures(num int32, loc *uint32) {
	gl.DeleteTextures(num, loc)
}

func (z *gl2) TexParameteri(target uint32, param uint32, value int32) {
	gl.TexParameteri(target, param, value)
}

func (z *gl2) TexImage2D(target uint32, level int32, internalFormat uint32, width, height, border int32, format, xtype uint32, pixels interface{}) {
	gl.TexImage2D(target, level, int32(internalFormat), width, height, border, format, xtype, gl.Ptr(pixels))
}

const (
	Escape       = Key(glfw.KeyEscape)
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
	F13          = Key(glfw.KeyF13)
	F14          = Key(glfw.KeyF14)
	F15          = Key(glfw.KeyF15)
	F16          = Key(glfw.KeyF16)
	F17          = Key(glfw.KeyF17)
	F18          = Key(glfw.KeyF18)
	F19          = Key(glfw.KeyF19)
	F20          = Key(glfw.KeyF20)
	F21          = Key(glfw.KeyF21)
	F22          = Key(glfw.KeyF22)
	F23          = Key(glfw.KeyF23)
	F24          = Key(glfw.KeyF24)
	F25          = Key(glfw.KeyF25)
	Up           = Key(glfw.KeyUp)
	Down         = Key(glfw.KeyDown)
	Left         = Key(glfw.KeyLeft)
	Right        = Key(glfw.KeyRight)
	LeftShift    = Key(glfw.KeyLeftShift)
	RightShift   = Key(glfw.KeyRightShift)
	LeftControl  = Key(glfw.KeyLeftControl)
	RightControl = Key(glfw.KeyRightControl)
	LeftAlt      = Key(glfw.KeyLeftAlt)
	RightAlt     = Key(glfw.KeyRightAlt)
	Tab          = Key(glfw.KeyTab)
	Space        = Key(glfw.KeySpace)
	Enter        = Key(glfw.KeyEnter)
	Backspace    = Key(glfw.KeyBackspace)
	Insert       = Key(glfw.KeyInsert)
	Delete       = Key(glfw.KeyDelete)
	PageUp       = Key(glfw.KeyPageUp)
	PageDown     = Key(glfw.KeyPageDown)
	Home         = Key(glfw.KeyHome)
	End          = Key(glfw.KeyEnd)
	Kp0          = Key(glfw.KeyKp0)
	Kp1          = Key(glfw.KeyKp1)
	Kp2          = Key(glfw.KeyKp2)
	Kp3          = Key(glfw.KeyKp3)
	Kp4          = Key(glfw.KeyKp4)
	Kp5          = Key(glfw.KeyKp5)
	Kp6          = Key(glfw.KeyKp6)
	Kp7          = Key(glfw.KeyKp7)
	Kp8          = Key(glfw.KeyKp8)
	Kp9          = Key(glfw.KeyKp9)
	KpDivide     = Key(glfw.KeyKpDivide)
	KpMultiply   = Key(glfw.KeyKpMultiply)
	KpSubtract   = Key(glfw.KeyKpSubtract)
	KpAdd        = Key(glfw.KeyKpAdd)
	KpDecimal    = Key(glfw.KeyKpDecimal)
	KpEqual      = Key(glfw.KeyKpEqual)
	KpEnter      = Key(glfw.KeyKpEnter)
	NumLock      = Key(glfw.KeyNumLock)
	CapsLock     = Key(glfw.KeyCapsLock)
	ScrollLock   = Key(glfw.KeyScrollLock)
	Pause        = Key(glfw.KeyPause)
	LeftSuper    = Key(glfw.KeyLeftSuper)
	RightSuper   = Key(glfw.KeyRightSuper)
	Menu         = Key(glfw.KeyMenu)
)
