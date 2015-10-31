# EdgeScroller Demo

## What does it do?
It demonstrates how one can move the camera round, by keeping the cursor close to the edges of the window. 

For doing so, it created a green background. This way, you'll notice the moving.  

## What are important aspects of the code?
This line is key in this demo:

* `game.AddSystem(engi.NewEdgeScroller(scrollSpeed, edgeMargin))`, to enable moving the camera around by using the edges of the window.
