//go:build headless || ios || android || vulkan || sdl
// +build headless ios android vulkan sdl

package engo

import "errors"

// Gampad is a configuration of a joystick that is able to be mapped to the
// SDL_GameControllerDB.
// For more info See https://www.glfw.org/docs/3.3/input_guide.html#gamepad_mapping
// This is here for compatibility for other builds, it does not work on
// mobile or vulkan builds yet.
type Gamepad struct {
	A, B, X, Y                            GamepadButton
	Back, Start, Guide                    GamepadButton
	DpadUp, DpadRight, DpadDown, DpadLeft GamepadButton
	LeftBumper, RightBumper               GamepadButton
	LeftThumb, RightThumb                 GamepadButton
	LeftX, LeftY                          AxisGamepad
	RightX, RightY                        AxisGamepad
	LeftTrigger, RightTrigger             AxisGamepad
}

func (gm *GamepadManager) registerGamepadImpl(name string) error {
	return errors.New("Gamepads are not available on this platform!")
}

func (gm *GamepadManager) updateImpl() {}
