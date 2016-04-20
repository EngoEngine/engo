# Exit Demo

## What does it do?
It demonstrates how one can detect the exit event and modify the default action

## What are important aspects of the code?
These things are key in this demo:

* `func (*DefaultScene) Exit() {}` Is called when the exit event is fired
* `engo.OverrideCloseAction()` overrides default close action for manual closing
