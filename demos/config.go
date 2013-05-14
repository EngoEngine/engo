package main

import (
	"github.com/ajhager/eng"
)

type Config struct {
	*eng.Game
}

func (c *Config) Init(config *eng.Config) {
	config.Title = "Config"
	config.Width = 800
	config.Height = 600
	config.Fullscreen = false
	config.Vsync = true
	config.Resizable = true
	config.Fsaa = 4
	config.PrintFPS = true
}

func main() {
	eng.Run(new(Config))
}
