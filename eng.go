// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

var (
	responder   Responder
	config      *Config
	timing      *stats
	defaultFont *Font
	bgColor     *Color
	GL          *gl2
)

type Action int
type Key int
type Modifier int

var (
	MOVE    = Action(0)
	PRESS   = Action(1)
	RELEASE = Action(2)
	SHIFT   = Modifier(0x0001)
	CONTROL = Modifier(0x0002)
	ALT     = Modifier(0x0004)
	SUPER   = Modifier(0x0008)
)

// Run should be called with a type that satisfies the Responder
// interface. Windows will be setup using your Config and a runloop
// will start, blocking the main thread and calling methods on the
// given responder.
func Run(title string, width, height int, fullscreen bool, r Responder) {
	RunConfig(&Config{title, width, height, fullscreen, true, false, 1, false}, r)
}

// RunConfig allows you to run with a custom configuration.
func RunConfig(c *Config, r Responder) {
	config = c
	responder = r
	bgColor = NewColorA(0, 0, 0, 0)
	GL = newgl2()
	run()
}

// Exit closes the window and breaks out of the game loop.
func Exit() {
	exit()
}

// Width returns the current window width.
func Width() int {
	return config.Width
}

// Height returns the current window height.
func Height() int {
	return config.Height
}

// SetBgColor sets the default opengl clear color.
func SetBgColor(c *Color) {
	bgColor = c.Copy()
}

// Dt returns the time since the last frame.
func Dt() float32 {
	return float32(timing.Dt)
}

// Fps returns the number of frames being rendered each second.
func Fps() float32 {
	return float32(timing.Fps)
}

// DefaultFont returns eng's built in font, creating it on first use.
func DefaultFont() *Font {
	if defaultFont == nil {
		defaultFont = NewBitmapFont(dfontimg(), dfonttxt)
	}
	return defaultFont
}
