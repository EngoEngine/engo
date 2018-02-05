package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct {
	audioSys *common.AudioSystem
}

type Whoop struct {
	ecs.BasicEntity
	common.AudioComponent
}

func (s *DefaultScene) Preload() {
	err := engo.Files.Load("326488.wav")
	if err != nil {
		log.Println(err)
	}
}

func (s *DefaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)
	s.audioSys = &common.AudioSystem{}

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(s.audioSys)
	//w.AddSystem(&WhoopSystem{})

	whoop := Whoop{BasicEntity: ecs.NewBasic()}
	whoopPlayer, err := common.LoadedPlayer("326488.wav")
	if err != nil {
		log.Fatalln(err)
	}
	whoop.AudioComponent = common.AudioComponent{Repeat: true, Player: whoopPlayer}
	whoop.AudioComponent.SetVolume(0.99)
	whoopPlayer.Play()

	// Let's add our whoop to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.AudioSystem:
			sys.Add(&whoop.BasicEntity, &whoop.AudioComponent)
		}
	}
}

func (d *DefaultScene) Hide() {
	d.audioSys.OtoPlayer.Close()
}

func (d *DefaultScene) Exit() {
	d.audioSys.OtoPlayer.Close()
}

func (*DefaultScene) Type() string { return "Game" }

type WhoopSystem struct {
	goingUp bool
}

// Remove is empty, because this system doesn't do anything with entities (note there's no `Add` method either)
func (w *WhoopSystem) Remove(basic ecs.BasicEntity) {}

func (w *WhoopSystem) Update(dt float32) {
	d := float64(dt * 0.1)
	if w.goingUp {
		common.MasterVolume += d
	} else {
		common.MasterVolume -= d
	}

	if common.MasterVolume < 0 {
		common.MasterVolume = 0
		w.goingUp = true
	} else if common.MasterVolume > 1 {
		common.MasterVolume = 1
		w.goingUp = false
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Audio Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &DefaultScene{})
}
