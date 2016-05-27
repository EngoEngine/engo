# EntityScroller Demo

This demo is built upon the TileMap Demo

## What does it do?
It demonstrates how one can use the EntityScroller system to scroll the camera to the position of a entity 

## What are important aspects of the code?
These lines are key in this demo:

* 'w.AddSystem(&common.EntityScroller{SpaceComponent: &character.SpaceComponent, TrackingBounds: levelData.Bounds()})' to add the system
