# Animation Demo

## What does it do?
It demonstrates how one can create animations.  

For doing so, it loaded a spritesheet, and then created several moving animations.     

## What are important aspects of the code?
These lines are key in this demo:

## What can I do in this demo?
* You can walk right (Right Arrow)
* Use a skill action (Space)

* `w.AddSystem(&engo.AnimationSystem{})`, to add/enable animations;
* `RunAction = &engo.AnimationAction{Name: "run", Frames: []int{16, 17, 18, 19, 20, 21}}`, for defining which frames were responsible for the `run` animation;
* `entity.AnimationComponent = engo.NewAnimationComponent(spriteSheet.Renderables(), 0.1)`, to create the animation component;
* `entity.AnimationComponent.AddAnimationActions(actions)`, to define the possible animations;
* `entity.AnimationComponent.SelectAnimationByAction(action)`, to set it to a specific animation;
