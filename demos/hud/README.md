# HUD Demo

## What does it do?
It demonstrates how one can have `Entity`s that aren't affected by the camera movement.    

For doing so, it created a green background -- which will be affected, and a violet sidebar, which won't be affected. 

## What are important aspects of the code?
These lines are key in this demo:

* `fieldRender.SetPriority(hudBackgroundPriority)`, to make sure the violet sidebar is rendered as HUD; 
* `hudBackgroundPriority = engo.PriorityLevel(engo.HUDGround)`, which defined `hudBackgroundPriority`. 
    * any value between `engo.HUDGround` and `engo.HighestGround` is rendered as HUD. 
