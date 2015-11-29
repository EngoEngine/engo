# Hello Demo

## What does it do?
It demonstrates how one can create a basic game. 

For doing so, it created an `Entity`, that is scaled up/down to show it's actually doing something.  

## What are important aspects of the code?
These things are key in this demo:

* `guy := ecs.NewEntity([]string{"RenderSystem", "ScaleSystem"})`, to create the guy and allow it to be rendered and scaled;
* The `ScaleSystem`, to enable the scaling up / down.
