# Scenes Demo

## What does it do?
It demonstrates how one can use multiple `Scene`s, and switch between them.  

## What are important aspects of the code?
These things are key in this demo:

* Defining two Scenes: `IconScene` and `RockScene`
* Giving one to `engi.Open` as the default `Scene`
* Registering the other with `engi.RegisterScene`, so we can later:
* Call `engi.SetSceneByName` to switch the `Scene`s. 
