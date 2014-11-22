//GAME IDEAS GUSY. GIVE ME SOME STUFF - QUICK as an example kinda trhing

package main

import (
	"github.com/paked/engi"
	"log"
)

type Game struct {
	engi.World
}

func (game Game) Preload() {
	engi.Files.Add("guy", "data/icon.png")
	engi.Files.Add("font", "data/font.go")
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&ControlSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ControlSystem"})
	texture := engi.Files.Image("guy")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "guy")

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}
	guy.AddComponent(&render)
	guy.AddComponent(&space)
	game.AddEntity(guy)
}

type ControlSystem struct {
	*engi.System
}

func (control *ControlSystem) New() {
	control.System = &engi.System{}
}

func (control ControlSystem) Name() string {
	return "ControlSystem"
}

func (control *ControlSystem) Update(entity *engi.Entity, dt float32) {
	space, hasSpace := entity.GetComponent("SpaceComponent").(*engi.SpaceComponent)
	if !hasSpace {
		return
	}

	speed := 200 * dt

	if engi.Keys.KEY_A.Down() {
		space.Position.X -= speed
	}

	if engi.Keys.KEY_D.Down() {
		space.Position.X += speed
	}

	if engi.Keys.KEY_W.Down() {
		space.Position.Y -= speed
	}

	if engi.Keys.KEY_S.Down() {
		space.Position.Y += speed
	}
}

func main() {
	log.Println("[Game] Says hello, written in github.com/paked/engi + Go")
	engi.Open("Stream Game", 800, 800, false, &Game{})
}
