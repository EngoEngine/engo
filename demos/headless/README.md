# Headless Pong Demo

## What does it do?
It copies 100% of the code at `pong/pong.go`, but attempts to run it in headless mode (without graphical display).   

## What are important aspects of the code?

Originally, `pong.go` has: 
```go
opts := engo.RunOptions{
    Title:         "Pong Demo",
    Width:         800,
    Height:        800,
    ScaleOnResize: true,
}
```

In order to turn it into headless mode, all we had to do, was change it to this:
```go
opts := engo.RunOptions{
    HeadlessMode: true,
}
```
