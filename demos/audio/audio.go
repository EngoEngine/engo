package main

import (
	"image/color"
	"log"

	"engo.io/engo"
	"engo.io/ecs"
)

type Game struct{}

func (game *Game) Preload() {
	engo.Files.Add("assets/326488.wav")
}

func (game *Game) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AudioSystem{})

	backgroundMusic := ecs.NewEntity("AudioSystem")
	backgroundMusic.AddComponent(&engo.AudioComponent{File: "326488.wav", Repeat: true, Background: true})

	err := w.AddEntity(backgroundMusic)
	if err != nil {
		log.Println(err)
	}
}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "Audio Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &Game{})
}
