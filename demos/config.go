package main

import (
	"github.com/ajhager/eng"
)

func main() {
	config := eng.NewConfig()
	config.Title = "Config"
	config.Width = 800
	config.Height = 600
	config.Fullscreen = false
	config.Vsync = true
	config.Resizable = true
	config.Fsaa = 4
	config.PrintFPS = true
	eng.RunConfig(config, new(eng.Game))
}
