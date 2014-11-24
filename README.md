#Engi 
A cross-platform game engine written in Go following an interpretation of the Entity Component System paradigm

More documentation will be added later, but here is a little getting started guide!
##Installation
```go get -u github.com/paked/engi```

*TODO* Write about the needed dependencies
##Getting Started
```
package main
   
import (
	"github.com/paked/engi"
)

type Game struct {
	engi.World
}

func (game *Game) Setup() {
	engi.SetBg(0xffffff)
	game.AddSystem(&engi.RenderSystem{})
}

func main() {
	engi.Open("Title", 800, 800, false, &Game{})
}

```

First we start off by declaring that it is a runnable file, then import the engi library. Inside the ```main()``` function we finish off by opening the window, the four parameters that are passed in are ```Window Title```, ```Window Width```, ```Window Height```, ```Fullscreen Mode (as a bool)``` and finally an instance of ```Game```.

If you were to run this code, a white 800x800 window would appear on your screen.


*TODO* Write about entities

*TODO* Write about components

*TODO* Write about systems




