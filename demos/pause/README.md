# Pause Demo

## What does it do?
It demonstrates how one can pause most aspects of the game; it uses "Scrolling Up" as a command to pause, and "Scrolling Down" as a command to unpause. 

For doing so, it created a few animations which will be paused. The rightmost animation should not be paused. 

For information about the animations, see the Animation Demo.  

## What are important aspects of the code?
These lines are key in this demo:

* `game.AddSystem(&engi.PauseSystem{World: &game.World})`, to add/enable the `PauseSystem`;
* `engi.Mailbox.Dispatch(engi.PauseMessage{amount > 0})`, to send a `PauseMessage` to the `PauseSystem`, indicating (un)pausing;
* `d_entity.AddComponent(&engi.UnpauseComponent{})`, to prevent the rightmost animation from being paused.
