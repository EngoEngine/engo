package common

import (
	"bytes"
	"image/color"
	"log"
	"strconv"

	"engo.io/ecs"
	"engo.io/engo"

	"golang.org/x/image/font/gofont/gomonobold"
)

// FPSSystem is a system for debugging that displays FPS to either the screen or
// the terminal.
type FPSSystem struct {
	Display, Terminal bool
	entity            struct {
		*ecs.BasicEntity
		*RenderComponent
		*SpaceComponent
	}
	elapsed, frames float32
	fnt             *Font
}

// New is called when FPSSystem is added to the world
func (f *FPSSystem) New(w *ecs.World) {
	if f.Display {
		if err := engo.Files.LoadReaderData("gomonobold_fps.ttf", bytes.NewReader(gomonobold.TTF)); err != nil {
			panic("unable to load gomonobold.ttf for the fps system! Error was: " + err.Error())
		}
		f.fnt = &Font{
			URL:  "gomonobold_fps.ttf",
			FG:   color.White,
			BG:   color.Black,
			Size: 32,
		}
		if err := f.fnt.CreatePreloaded(); err != nil {
			panic("unable to create gomonobold.ttf for the fps system! Error was: " + err.Error())
		}
		txt := Text{
			Font: f.fnt,
			Text: "Hello world!",
		}
		b := ecs.NewBasic()
		f.entity.BasicEntity = &b
		f.entity.RenderComponent = &RenderComponent{
			Drawable: txt,
		}
		f.entity.RenderComponent.SetShader(HUDShader)
		f.entity.RenderComponent.SetZIndex(1000)
		f.entity.SpaceComponent = &SpaceComponent{}
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *RenderSystem:
				sys.Add(f.entity.BasicEntity, f.entity.RenderComponent, f.entity.SpaceComponent)
			}
		}
	}
}

// Add doesn't do anything since New creates the only entity used
func (*FPSSystem) Add() {}

// Remove doesn't do anything since New creates the only entity used
func (*FPSSystem) Remove(b ecs.BasicEntity) {}

// Update changes the dipslayed text and prints to the terminal every second
// to report the FPS
func (f *FPSSystem) Update(dt float32) {
	f.elapsed += dt
	f.frames += 1
	text := "FPS: " + strconv.FormatFloat(float64(f.frames/f.elapsed), 'G', 5, 32)
	if f.elapsed >= 1 {
		if f.Display {
			f.entity.Drawable = Text{
				Font: f.fnt,
				Text: text,
			}
		}
		if f.Terminal {
			log.Println(text)
		}
		f.frames = 0
		f.elapsed -= 1
	}
}
