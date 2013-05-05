eng 2D game library for go
===

+NOTHING ABOUT THIS LIBRARY IS STABLE+

eng depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install libglfw-dev
* OSX: brew install glfw
* Windows: download the glfw binaries, then drop the GL directory into C:\MinGW\include and the files for your arch under libmingw into C:\MinGW\lib. You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
`go get github.com/ajhager/eng`

Try it!
-------
```go
package main

import (
    "github.com/ajhager/eng"
)

type Demo struct {
    *eng.Game
}

func (d *Demo) Draw() {
    eng.Print("Hello, world!", 500, 300)
}

func main() {
    eng.Run(&Demo{})
}
```
