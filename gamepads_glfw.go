// +build darwin,!arm,!arm64 linux windows
// +build !ios,!android,!js,!sdl,!headless,!vulkan

package engo

import (
	"errors"

	"github.com/go-gl/glfw/v3.3/glfw"
)

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

	joystick  glfw.Joystick
	id        string
	connected bool
}

var joys = []glfw.Joystick{
	glfw.Joystick1, glfw.Joystick2, glfw.Joystick3, glfw.Joystick4,
	glfw.Joystick5, glfw.Joystick6, glfw.Joystick7, glfw.Joystick8,
	glfw.Joystick9, glfw.Joystick10, glfw.Joystick11, glfw.Joystick12,
	glfw.Joystick13, glfw.Joystick14, glfw.Joystick15, glfw.Joystick16,
}

var usedjoys = []glfw.Joystick{}

func (gm *GamepadManager) registerGamepadImpl(name string) error {
	found := false
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
joyLoop:
	for _, joy := range joys {
		for _, u := range usedjoys {
			if joy == u {
				continue joyLoop
			}
		}
		if joy.IsGamepad() {
			gm.gamepads[name] = &Gamepad{
				joystick:  joy,
				id:        joy.GetGUID(),
				connected: true,
			}
			found = true
			usedjoys = append(usedjoys, joy)
			break joyLoop
		}
	}
	if !found {
		warning("Unable to locate any usable gamepads.")
		gm.gamepads[name] = &Gamepad{}
		return errors.New("unable to locate any usable gamepads \ngamepad will be added when a new one is plugged in")
	}
	return nil
}

func (gm *GamepadManager) updateImpl() {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	for name, gamepad := range gm.gamepads {
		if !gamepad.connected {
			warning("Gamepad " + name + " was not available for update!")
			continue
		}
		if gamepad.joystick.Present() {
			state := gamepad.joystick.GetGamepadState()

			if state.Buttons[glfw.ButtonA] == glfw.Press {
				gamepad.A.set(true)
			} else if state.Buttons[glfw.ButtonA] == glfw.Release {
				gamepad.A.set(false)
			}

			if state.Buttons[glfw.ButtonB] == glfw.Press {
				gamepad.B.set(true)
			} else if state.Buttons[glfw.ButtonB] == glfw.Release {
				gamepad.B.set(false)
			}

			if state.Buttons[glfw.ButtonX] == glfw.Press {
				gamepad.X.set(true)
			} else if state.Buttons[glfw.ButtonX] == glfw.Release {
				gamepad.X.set(false)
			}

			if state.Buttons[glfw.ButtonY] == glfw.Press {
				gamepad.Y.set(true)
			} else if state.Buttons[glfw.ButtonY] == glfw.Release {
				gamepad.Y.set(false)
			}

			if state.Buttons[glfw.ButtonBack] == glfw.Press {
				gamepad.Back.set(true)
			} else if state.Buttons[glfw.ButtonBack] == glfw.Release {
				gamepad.Back.set(false)
			}

			if state.Buttons[glfw.ButtonStart] == glfw.Press {
				gamepad.Start.set(true)
			} else if state.Buttons[glfw.ButtonStart] == glfw.Release {
				gamepad.Start.set(false)
			}

			if state.Buttons[glfw.ButtonGuide] == glfw.Press {
				gamepad.Guide.set(true)
			} else if state.Buttons[glfw.ButtonGuide] == glfw.Release {
				gamepad.Guide.set(false)
			}

			if state.Buttons[glfw.ButtonDpadUp] == glfw.Press {
				gamepad.DpadUp.set(true)
			} else if state.Buttons[glfw.ButtonDpadUp] == glfw.Release {
				gamepad.DpadUp.set(false)
			}

			if state.Buttons[glfw.ButtonDpadRight] == glfw.Press {
				gamepad.DpadRight.set(true)
			} else if state.Buttons[glfw.ButtonDpadRight] == glfw.Release {
				gamepad.DpadRight.set(false)
			}

			if state.Buttons[glfw.ButtonDpadDown] == glfw.Press {
				gamepad.DpadDown.set(true)
			} else if state.Buttons[glfw.ButtonDpadDown] == glfw.Release {
				gamepad.DpadDown.set(false)
			}

			if state.Buttons[glfw.ButtonDpadLeft] == glfw.Press {
				gamepad.DpadLeft.set(true)
			} else if state.Buttons[glfw.ButtonDpadLeft] == glfw.Release {
				gamepad.DpadLeft.set(false)
			}

			if state.Buttons[glfw.ButtonLeftBumper] == glfw.Press {
				gamepad.LeftBumper.set(true)
			} else if state.Buttons[glfw.ButtonLeftBumper] == glfw.Release {
				gamepad.LeftBumper.set(false)
			}

			if state.Buttons[glfw.ButtonRightBumper] == glfw.Press {
				gamepad.RightBumper.set(true)
			} else if state.Buttons[glfw.ButtonRightBumper] == glfw.Release {
				gamepad.RightBumper.set(false)
			}

			if state.Buttons[glfw.ButtonLeftThumb] == glfw.Press {
				gamepad.LeftThumb.set(true)
			} else if state.Buttons[glfw.ButtonLeftThumb] == glfw.Release {
				gamepad.LeftThumb.set(false)
			}

			if state.Buttons[glfw.ButtonRightThumb] == glfw.Press {
				gamepad.RightThumb.set(true)
			} else if state.Buttons[glfw.ButtonRightThumb] == glfw.Release {
				gamepad.RightThumb.set(false)
			}

			gamepad.LeftX.set(state.Axes[glfw.AxisLeftX])
			gamepad.LeftY.set(state.Axes[glfw.AxisLeftY])
			gamepad.RightX.set(state.Axes[glfw.AxisRightX])
			gamepad.RightY.set(state.Axes[glfw.AxisRightY])
			gamepad.LeftTrigger.set(state.Axes[glfw.AxisLeftTrigger])
			gamepad.RightTrigger.set(state.Axes[glfw.AxisRightTrigger])
		} else {
			gamepad.connected = false
			warning("Gamepad " + name + " was not available to update!")
		}
	}
}
