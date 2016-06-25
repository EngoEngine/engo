# Mouse Follow Demo

## What does it do?
It demonstrates how one can implement a sprite that will follow the user's mouse movements.

## What are important aspects of the code?
These things are key in this demo:

* `FollowSystem`, which implements an update method that sets the SpaceComponent(sprite)'s position to that of the current mouse position. Like so:
```go
func (s *FollowSystem) Update(dt float32) {
    for _, e := range s.entities {
        e.SpaceComponent.Position.X += engo.Input.Axis(engo.DefaultMouseXAxis).Value()
        e.SpaceComponent.Position.Y += engo.Input.Axis(engo.DefaultMouseYAxis).Value()
    }
}
```
