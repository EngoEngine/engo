package main

import (
	"github.com/ajhager/eng"
)

func main() {
	eng.Run("Empty", 1024, 640, false, new(eng.Game))
}
