# MouseRotation Demo

## What does it do?
It demonstrates how one can rotate the camera, using the middle mouse button.  

For doing so, it created a green background. This way, you'll notice the rotation. 

## What are important aspects of the code?
These lines are key in this demo:

* `w.AddSystem(&common.MouseRotator{rotationSpeed})`, to enable the scrolling with the mouse wheel. 
