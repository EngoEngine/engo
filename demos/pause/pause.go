package main

import (
	"github.com/paked/engi"
)

var World *GameWorld

type GameWorld struct {
	engi.World
	RUN_ACTION   *engi.AnimationAction
	WALK_ACTION  *engi.AnimationAction
	STOP_ACTION  *engi.AnimationAction
	SKILL_ACTION *engi.AnimationAction
	DIE_ACTION   *engi.AnimationAction
	actions      []*engi.AnimationAction
	activeAction *engi.AnimationAction
}

func (game *GameWorld) Preload() {
	engi.Files.Add("assets/hero.png")
	game.STOP_ACTION = &engi.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	game.RUN_ACTION = &engi.AnimationAction{Name: "run", Frames: []int{16, 17, 18, 19, 20, 21}}
	game.WALK_ACTION = &engi.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	game.SKILL_ACTION = &engi.AnimationAction{Name: "skill", Frames: []int{44, 45, 46, 47, 48, 49, 50, 51, 52, 53}}
	game.DIE_ACTION = &engi.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	game.actions = []*engi.AnimationAction{game.DIE_ACTION, game.STOP_ACTION, game.WALK_ACTION, game.RUN_ACTION, game.SKILL_ACTION}
	game.activeAction = game.RUN_ACTION
}

func (game *GameWorld) Setup() {
	engi.SetBg(0xFFFFFF)

	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&engi.AnimationSystem{})
	game.AddSystem(&engi.PauseSystem{World: &game.World})

	spriteSheet := engi.NewSpritesheetFromFile("hero.png", 150, 150)

	game.AddEntity(game.CreateEntity(&engi.Point{0, 0}, spriteSheet, game.RUN_ACTION))
	game.AddEntity(game.CreateEntity(&engi.Point{300, 0}, spriteSheet, game.WALK_ACTION))
	game.AddEntity(game.CreateEntity(&engi.Point{600, 0}, spriteSheet, game.STOP_ACTION))
	game.AddEntity(game.CreateEntity(&engi.Point{900, 0}, spriteSheet, game.SKILL_ACTION))

	// This animation is special
	death := game.CreateEntity(&engi.Point{1200, 0}, spriteSheet, game.DIE_ACTION)
	// ... because now, it's not affected by pausing
	death.AddComponent(&engi.UnpauseComponent{})
	game.AddEntity(death)
}

func (game *GameWorld) CreateEntity(point *engi.Point, spriteSheet *engi.Spritesheet, action *engi.AnimationAction) *engi.Entity {
	entity := engi.NewEntity([]string{"AnimationSystem", "RenderSystem"})

	space := &engi.SpaceComponent{*point, 0, 0}
	render := engi.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engi.Point{3, 3}, "hero")
	animation := engi.NewAnimationComponent(spriteSheet.Renderables(), 0.1)
	animation.AddAnimationActions(game.actions)
	animation.SelectAnimationByAction(action)
	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

func (game *GameWorld) Scroll(amount float32) {
	// Pause the game if we're scrolling up; else, unpause
	engi.Mailbox.Dispatch(engi.PauseMessage{amount > 0})
}

func main() {
	World = &GameWorld{}
	engi.Open("Pause Demo", 1024, 640, false, World)
}
