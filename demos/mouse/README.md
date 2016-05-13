# Mouse Demo

## What does it do?
It demonstrates how one can use mouse events in `engo`. 

## What are important aspects of the code?
These lines are key in this demo:

* `common.MouseComponent` within the definition of `type Guy`, which adds a `MouseComponent` which stores information (such as hover events). 
* `w.AddSystem(&common.MouseSystem{})`, which adds the `MouseSystem` to this `Scene`;

