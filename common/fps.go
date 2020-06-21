package common

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"

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
	elapsed float32
	Font    *Font // Font used to display the FPS to the screen, defaults to gomonobold
}

// New is called when FPSSystem is added to the world
func (f *FPSSystem) New(w *ecs.World) {
	if f.Display {
		if f.Font == nil {
			if err := engo.Files.LoadReaderData("gomonobold_fps.ttf", bytes.NewReader(gomonobold.TTF)); err != nil {
				panic("unable to load gomonobold.ttf for the fps system! Error was: " + err.Error())
			}

			f.Font = &Font{
				URL:  "gomonobold_fps.ttf",
				FG:   color.White,
				BG:   color.Black,
				Size: 32,
			}

			if err := f.Font.CreatePreloaded(); err != nil {
				panic("unable to create gomonobold.ttf for the fps system! Error was: " + err.Error())
			}
		}

		txt := Text{
			Font: f.Font,
			Text: f.DisplayString(),
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
	text := f.DisplayString()
	if f.elapsed >= 1 {
		if f.Display {
			f.entity.Drawable = Text{
				Font: f.Font,
				Text: text,
			}
		}
		if f.Terminal {
			log.Println(text)
		}
		f.elapsed--
	}
}

// DisplayString returns the display string in the format FPS: 60
func (f *FPSSystem) DisplayString() string {
	if engo.Time == nil {
		return ""
	}
	return fmt.Sprintf("FPS: %g", engo.Time.FPS())
}
