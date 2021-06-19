//+build js

package engo

import "errors"

// Gampad is a configuration of a joystick that is able to be mapped to the
// SDL_GameControllerDB.
// For more info See https://www.glfw.org/docs/3.3/input_guide.html#gamepad_mapping
type Gamepad struct {
	A, B, X, Y                            GamepadButton
	Back, Start, Guide                    GamepadButton
	DpadUp, DpadRight, DpadDown, DpadLeft GamepadButton
	LeftBumper, RightBumper               GamepadButton
	LeftThumb, RightThumb                 GamepadButton
	LeftX, LeftY                          AxisGamepad
	RightX, RightY                        AxisGamepad
	LeftTrigger, RightTrigger             AxisGamepad

	id        string
	connected bool
}

var usedGpds []string

func (gm *GamepadManager) registerGamepadImpl(name string) error {
	gpds := window.Get("navigator").Call("getGamepads")
	found := false
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	for i := 0; i < gpds.Length(); i++ {
		if gpds.Index(i).Get("mapping").String() != "standard" {
			continue
		}
		if gpds.Index(i).IsNull() {
			continue
		}
		gpid := gpds.Index(i).Get("id").String()
		gm.gamepads[name] = &Gamepad{
			id:        gpid,
			connected: true,
		}
		found = true
	}
	if !found {
		warning("Unable to locate any usable gamepads.")
		gm.gamepads[name] = &Gamepad{}
		return errors.New("unable to locate any usable gamepads \ngamepad will be added when a new one is plugged in")
	}
	return nil
}

func (gm *GamepadManager) updateImpl() {
	gpds := window.Get("navigator").Call("getGamepads")
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	for name, gamepad := range gm.gamepads {
		if !gamepad.connected {
			warning("Gamepad " + name + " was not available for update!")
			continue
		}
		for i := 0; i < gpds.Length(); i++ {
			if gpds.Index(i).IsNull() {
				continue
			}
			gpid := gpds.Index(i).Get("id").String()
			if gpid == gamepad.id {
				if gpds.Index(i).Get("connected").Bool() {
					gamepad.A.set(gpds.Index(i).Get("buttons").Index(0).Get("pressed").Bool())
					gamepad.B.set(gpds.Index(i).Get("buttons").Index(1).Get("pressed").Bool())
					gamepad.X.set(gpds.Index(i).Get("buttons").Index(2).Get("pressed").Bool())
					gamepad.Y.set(gpds.Index(i).Get("buttons").Index(3).Get("pressed").Bool())
					gamepad.LeftBumper.set(gpds.Index(i).Get("buttons").Index(4).Get("pressed").Bool())
					gamepad.RightBumper.set(gpds.Index(i).Get("buttons").Index(5).Get("pressed").Bool())
					if gpds.Index(i).Get("buttons").Index(6).Get("pressed").Bool() {
						gamepad.LeftTrigger.set(1.0)
					} else {
						gamepad.LeftTrigger.set(0.0)
					}
					if gpds.Index(i).Get("buttons").Index(7).Get("pressed").Bool() {
						gamepad.RightTrigger.set(1.0)
					} else {
						gamepad.RightTrigger.set(0.0)
					}
					gamepad.Back.set(gpds.Index(i).Get("buttons").Index(8).Get("pressed").Bool())
					gamepad.Start.set(gpds.Index(i).Get("buttons").Index(9).Get("pressed").Bool())
					gamepad.LeftThumb.set(gpds.Index(i).Get("buttons").Index(10).Get("pressed").Bool())
					gamepad.RightThumb.set(gpds.Index(i).Get("buttons").Index(11).Get("pressed").Bool())
					gamepad.DpadUp.set(gpds.Index(i).Get("buttons").Index(12).Get("pressed").Bool())
					gamepad.DpadDown.set(gpds.Index(i).Get("buttons").Index(13).Get("pressed").Bool())
					gamepad.DpadLeft.set(gpds.Index(i).Get("buttons").Index(14).Get("pressed").Bool())
					gamepad.DpadRight.set(gpds.Index(i).Get("buttons").Index(15).Get("pressed").Bool())
					gamepad.Guide.set(gpds.Index(i).Get("buttons").Index(16).Get("pressed").Bool())
					gamepad.LeftX.set(float32(gpds.Index(i).Get("axes").Index(0).Float()))
					gamepad.LeftY.set(float32(gpds.Index(i).Get("axes").Index(1).Float()))
					gamepad.RightX.set(float32(gpds.Index(i).Get("axes").Index(2).Float()))
					gamepad.RightY.set(float32(gpds.Index(i).Get("axes").Index(3).Float()))
				} else {
					gamepad.connected = false
					warning("Gamepad " + name + " was not available to update!")
				}
			}
		}
	}
}
