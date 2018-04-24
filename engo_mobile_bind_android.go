//+build android,mobilebind

package engo

//TouchEvent handles the touch events sent from Android and puts them in the InputManager
func TouchEvent(x, y, id, action int) {
	Input.Mouse.X = float32(x) / opts.GlobalScale.X
	Input.Mouse.Y = float32(y) / opts.GlobalScale.Y
	switch action {
	case 0, 5:
		Input.Mouse.Action = Press
		Input.Touches[id] = Point{
			X: float32(x) / opts.GlobalScale.X,
			Y: float32(y) / opts.GlobalScale.Y,
		}
	case 1, 6:
		Input.Mouse.Action = Release
		delete(Input.Touches, id)
	case 2:
		Input.Mouse.Action = Move
		Input.Touches[id] = Point{
			X: float32(x) / opts.GlobalScale.X,
			Y: float32(y) / opts.GlobalScale.Y,
		}
	}
}
