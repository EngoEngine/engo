# Hide Demo

## What does it do?
It demonstrates how one can show and hide a character.

For doing so, it created an `Entity` (rock), that hides and shows as time passes by.  

## What are important aspects of the code?
These things are key in this demo:

* HideSystem that shows and hides an `Entity` (rock)

```go
if rand.Int()%10 == 0 {
  render.SetPriority(engi.Hidden)
} else {
  render.SetPriority(engi.MiddleGround)
}
```
