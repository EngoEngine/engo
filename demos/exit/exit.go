package main

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
)

type Game struct{}

func (game *Game) Preload()           {}
func (game *Game) Setup(w *ecs.World) {}
func (*Game) Hide()                   {}
func (*Game) Show()                   {}

func (*Game) Exit() {
	log.Println("[GAME] Exit event called")
	//Here if you want you can prompt the user if they're sure they want to close
	log.Println("[GAME] Manually closing")
	engo.Exit()
}

func (*Game) Type() string { return "Game" }

func main() {
	opts := engo.RunOptions{
		Title:  "Exit Demo",
		Width:  1024,
		Height: 640,
	}
	engo.OverrideCloseAction()
	engo.Run(opts, &Game{})
}
