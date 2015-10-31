# Pong Demo

## What does it do?
It combines a lot of features of `engi`, to demonstrate the playing of a game of pong. 

## What are important aspects of the code?
These things are key in this demo:

* Preloading files in `PongGame.Preload()`;
* The `SpeedSystem`, to make sure the ball increases speed gradually;
* The `ControlSystem`, to allow moving paddles up/down with WASD and arrow keys;
* The `BallSystem`, to keep track of the location of the ball, sending messages to the `ScoreSystem` when needed;
* The `ScoreSystem`, to keep track of the score of each player;
* Creating some `Entity`s and requiring the correct `System`s. 
