# Exit Demo

## What does it do?
It demonstrates how one can detect the exit event and modify the default action

## What are important aspects of the code?
These things are key in this demo:

* `func (*game) Exit() {}` Is called when the exit event is fired
* `DefaultCloseAction: false,` Doesn't close the game when exit is requested and lets use do
manual handling and closing
