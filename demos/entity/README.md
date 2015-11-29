# Entity Demo

## What does it do?
It demonstrates how one can create an `Entity`.   

## What are important aspects of the code?
These lines are key in this demo:

* `guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})`, to define the `Entity`;
* `render := engi.NewRenderComponent(texture, engi.Point{8, 8}, "guy")`, to add something renderable to the `Entity`;
* `guy.AddComponent(render)`, to combine the renderable component and the `Entity`;
* `game.AddEntity(guy)`, to actually add the `Entity` to the game.
