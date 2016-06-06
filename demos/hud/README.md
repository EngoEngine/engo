# HUD Demo

## What does it do?
It demonstrates how one can have `Entity`s that aren't affected by the camera movement.    

For doing so, it created a green background -- which will be affected, and a violet sidebar, which won't be affected. 

## What are important aspects of the code?
These lines are key in this demo:

* `hudBg.RenderComponent.SetZIndex(1)`, to ensure it's rendered on top of the default background (with z-index 0)
* `hudBg.RenderComponent.SetShader(common.HUDShader)`, to enable the HUDShader instead of the DefaultShader
