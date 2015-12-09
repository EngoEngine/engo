# Mouse Demo

## What does it do?
It demonstrates how one can use mouse events in `engi`. 

## What are important aspects of the code?
These lines are key in this demo:

* `entity.AddComponent(&engi.MouseComponent{})` which adds a `MouseComponent` to the `Entity` of which we want to know mouse stuff; 
* `w.AddSystem(&engi.MouseSystem{})` which adds the `MouseSystem` to this `Scene`;
