# Entity Demo

## What does it do?
It demonstrates how one can create an `Entity`.   

## What are important aspects of the code?
These lines are key in this demo:

* `guy := Guy{BasicEntity: ecs.NewBasic()}`, to define the basics (`ecs.NewBasic()` generates a new UID for the entity);
* `guy.RenderComponent = engo.NewRenderComponent(texture, engo.Point{8, 8}, "guy")`, to add something renderable;
* `guy.SpaceComponent = common.SpaceComponent{...}`, to define the location of the guy;
* The lines which add the `guy` to the appropriate systems. 
