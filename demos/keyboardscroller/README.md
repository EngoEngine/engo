# KeyboardScroller Demo

## What does it do?
It demonstrates how one can move the camera round, by pressing WASD.   

For doing so, it created a green background. This way, you'll notice the moving.  

## What are important aspects of the code?
This line is key in this demo:

* `game.AddSystem(engi.NewKeyboardScroller(scrollSpeed, engi.W, engi.D, engi.S, engi.A))`, to enable moving the camera by using the keyboard.
