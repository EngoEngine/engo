// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"strings"
)

var (
	responder Responder
	config    *Config
	timing    *stats
	bgColor   *Color
	Files     *Loader
	GL        *gl2
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

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

type Resource struct {
	kind string
	name string
	url  string
}

type Loader struct {
	resources []Resource
	images    map[string]Image
	jsons     map[string]string
}

func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		images:    make(map[string]Image),
		jsons:     make(map[string]string),
	}
}

func (l *Loader) Add(name, url string) {
	parts := strings.Split(url, ".")
	kind := parts[len(parts)-1]
	l.resources = append(l.resources, Resource{kind, name, url})
}

func (l *Loader) Image(name string) Image {
	return l.images[name]
}

func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		switch r.kind {
		case "png":
			data, err := loadImage(r)
			if err == nil {
				l.images[r.name] = data
			}
		case "json":
			data, err := loadJson(r)
			if err == nil {
				l.jsons[r.name] = data
			}
		}
	}
	onFinish()
}

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
	Files = NewLoader()
	bgColor = NewColorA(0, 0, 0, 0)
	run()
}

// Exit closes the window and breaks out of the game loop.
func Exit() {
	exit()
}

// Width returns the current window width.
func Width() float32 {
	return float32(config.Width)
}

// Height returns the current window height.
func Height() float32 {
	return float32(config.Height)
}

// SetBgColor sets the default opengl clear color.
func SetBgColor(c *Color) {
	bgColor = c.Copy()
}

// Fps returns the number of frames being rendered each second.
func Fps() float32 {
	return float32(timing.Fps)
}
