# Zoom Demo

## What does it do?
It demonstrates how one can zoom in/out, by using the mouse wheel. 

For doing so, it created a green background. This way, you'll notice the zooming. 

## What are important aspects of the code?
These lines are key in this demo:

```go
// Scroll is called whenever the mouse wheel scrolls
func (game *Game) Scroll(amount float32) {
	// Adding this line, allows for zooming on scrolling the mouse wheel
	engi.Mailbox.Dispatch(engi.CameraMessage{Axis: engi.ZAxis, Value: amount * zoomSpeed, Incremental: true})
}
```