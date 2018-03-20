//+build demo

package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type DefaultScene struct{}

type Whoop struct {
	ecs.BasicEntity
	common.AudioComponent
}

func (s *DefaultScene) Preload() {
	err := engo.Files.Load("326488.wav")
	if err != nil {
		log.Println(err)
	}
	err = engo.Files.Load("326064.wav")
	if err != nil {
		log.Println(err)
	}

	engo.Input.RegisterButton("whoop", engo.KeySpace)
}

func (s *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AudioSystem{})
	w.AddSystem(&WhoopSystem{})

	birds := Whoop{BasicEntity: ecs.NewBasic()}
	birdPlayer, err := common.LoadedPlayer("326488.wav")
	if err != nil {
		log.Fatalln(err)
	}
	birds.AudioComponent = common.AudioComponent{Player: birdPlayer}
	birdPlayer.Play()
	birdPlayer.Repeat = true

	// Let's add our birds to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.AudioSystem:
			sys.Add(&birds.BasicEntity, &birds.AudioComponent)
		}
	}

	whoop := Whoop{BasicEntity: ecs.NewBasic()}
	whoopPlayer, err := common.LoadedPlayer("326064.wav")
	if err != nil {
		log.Fatalln(err)
	}
	whoop.AudioComponent = common.AudioComponent{Player: whoopPlayer}

	// Let's add our whoop to the appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.AudioSystem:
			sys.Add(&whoop.BasicEntity, &whoop.AudioComponent)
		case *WhoopSystem:
			sys.Add(&whoop.AudioComponent)
		}
	}
}

func (*DefaultScene) Type() string { return "Game" }

type WhoopSystem struct {
	goingUp bool
	player  *common.Player
}

func (w *WhoopSystem) Add(audio *common.AudioComponent) {
	w.player = audio.Player
}

func (w *WhoopSystem) Remove(basic ecs.BasicEntity) {}

func (w *WhoopSystem) Update(dt float32) {
	if btn := engo.Input.Button("whoop"); btn.JustPressed() {
		if !w.player.IsPlaying() {
			w.player.Rewind()
			w.player.Play()
		}
	}
	d := float64(dt * 0.1)
	volume := w.player.GetVolume()
	if w.goingUp {
		volume += d
	} else {
		volume -= d
	}

	if volume < 0 {
		w.player.SetVolume(0.0)
		w.goingUp = true
	} else if volume > 1 {
		w.player.SetVolume(1.0)
		w.goingUp = false
	} else {
		w.player.SetVolume(volume)
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
