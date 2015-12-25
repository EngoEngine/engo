package engi

import (
	"fmt"

	"github.com/paked/engi/ecs"
	"github.com/paked/webgl"
)

var (
	Time        *Clock
	Files       *Loader
	Gl          *webgl.Context
	WorldBounds AABB

	currentWorld *ecs.World
	currentScene Scene
	Mailbox      *MessageManager
	cam          *cameraSystem

	scaleOnResize   = false
	fpsLimit        = 120
	headless        = false
	vsync           = true
	resetLoopTicker = make(chan bool, 1)
)

type RunOptions struct {
	// NoRun indicates the Open function should return immediately, without looping
	NoRun bool

	// Title is the Window title
	Title string

	// HeadlessMode indicates whether or not OpenGL calls should be made
	HeadlessMode bool

	Fullscreen    bool
	Width, Height int

	// VSync indicates whether or not OpenGL should wait for the monitor to swp the buffers
	VSync bool

	// ScaleOnResize indicates whether or not engi should make things larger/smaller whenever the screen resizes
	ScaleOnResize bool

	// FPSLimit indicates the maximum number of frames per second
	FPSLimit int
}

func Open(opts RunOptions, defaultScene Scene) {
	// Save settings
	SetScaleOnResize(opts.ScaleOnResize)
	SetFPSLimit(opts.FPSLimit)
	vsync = opts.VSync

	if opts.HeadlessMode {
		headless = true

		if !opts.NoRun {
			runHeadless(defaultScene)
		}
	} else {
		CreateWindow(opts.Title, opts.Width, opts.Height, opts.Fullscreen)
		defer DestroyWindow()

		if !opts.NoRun {
			runLoop(defaultScene, false)
		}
	}
}

func SetBg(color uint32) {
	if !headless {
		r := float32((color>>16)&0xFF) / 255.0
		g := float32((color>>8)&0xFF) / 255.0
		b := float32(color&0xFF) / 255.0
		Gl.ClearColor(r, g, b, 1.0)
	}
}

func SetScaleOnResize(b bool) {
	scaleOnResize = b
}

func SetFPSLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("FPS Limit out of bounds. Requires > 0")
	}
	fpsLimit = limit
	resetLoopTicker <- true
	return nil
}
