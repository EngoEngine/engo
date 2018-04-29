//+build ios darwin,arm darwin,arm64
//+build mobilebind

package engo

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Foundation -framework UIKit
//
// #import <UIKit/UIKit.h>
import "C"

//TouchEvent handles the touch events sent from Android and puts them in the InputManager
func TouchEvent(x, y, id, action int) {
	Input.Mouse.X = float32(x) / opts.GlobalScale.X
	Input.Mouse.Y = float32(y) / opts.GlobalScale.Y
	switch action {
	case C.UITouchPhaseBegan, C.UITouchPhaseStationary:
		Input.Mouse.Action = Press
		Input.Touches[id] = Point{
			X: float32(x) / opts.GlobalScale.X,
			Y: float32(y) / opts.GlobalScale.Y,
		}
	case C.UITouchPhaseEnded, C.UITouchPhaseCancelled:
		Input.Mouse.Action = Release
		delete(Input.Touches, id)
	case C.UITouchPhaseMoved:
		Input.Mouse.Action = Move
		Input.Touches[id] = Point{
			X: float32(x) / opts.GlobalScale.X,
			Y: float32(y) / opts.GlobalScale.Y,
		}
	}
}
