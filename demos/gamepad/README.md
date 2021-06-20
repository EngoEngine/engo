# Gamepad Demo

## What does it do?
It demonstrates how one use a gamepad.   

## What are important aspects of the code?
These lines are key in this demo:

* Registering the Gamepad
  ```
	err := engo.Input.RegisterGamepad("Player1")
	if err != nil {
		println("Unable to find suitable Gamepad. Error was: ", err.Error())
	}
  ```

* Retrieve the gamepad during Update of the InputSystem.
    !!! Make sure to check if the gamepad is not nil before trying to use it !!!
  ```
  // Retrieve the Gamepad
  gamepad := engo.Input.Gamepad("Player1")
  if gamepad == nil {
    println("No gamepad found for Player1.")
    return
  }
  ```

  * Using gamepad keys

  ```
  if gamepad.A.Up() {
    entity.Color = color.White
  } else if gamepad.A.JustPressed() {
    entity.Color = color.RGBA{0, 255, 0, 255}
  } else if gamepad.A.Down() {
    entity.Color = color.RGBA{255, 0, 0, 255}
  }
  ```

  * Using gamepad Axes

  ```
  gamepad.RightX.Value()
  ```
