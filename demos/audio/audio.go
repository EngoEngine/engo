package main

import (
	"fmt"
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
)

type Game struct{}

func (game *Game) Preload() {
	engo.Files.Add("assets/326488.wav")
}

func (game *Game) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AudioSystem{})
	w.AddSystem(&WhoopSystem{})

	backgroundMusic := ecs.NewEntity("AudioSystem", "WhoopSystem")
	backgroundMusic.AddComponent(&engo.AudioComponent{File: "326488.wav", Repeat: true, Background: true, RawVolume: 1})

	err := w.AddEntity(backgroundMusic)
	if err != nil {
		log.Println(err)
	}
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Exit()        {}
func (*Game) Type() string { return "Game" }

type WhoopSystem struct {
	goingUp bool
}

func (WhoopSystem) Type() string             { return "WhoopSystem" }
func (WhoopSystem) Priority() int            { return 0 }
func (WhoopSystem) New(w *ecs.World)         {}
func (WhoopSystem) AddEntity(*ecs.Entity)    {}
func (WhoopSystem) RemoveEntity(*ecs.Entity) {}

func (ws *WhoopSystem) Update(dt float32) {
	d := float64(dt * 0.1)
	if ws.goingUp {
		engo.MasterVolume += d
	} else {
		engo.MasterVolume -= d
	}

	if engo.MasterVolume < 0 {
		engo.MasterVolume = 0
		ws.goingUp = true
	} else if engo.MasterVolume > 1 {
		engo.MasterVolume = 1
		ws.goingUp = false
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Audio Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &Game{})
}
