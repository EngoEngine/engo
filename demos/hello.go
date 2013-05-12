package main

import (
	"github.com/ajhager/eng"
)

type Hello struct {
	*eng.Game
	batch *eng.Batch
}

func (g *Hello) Open() {
	g.batch = eng.NewBatch()
}

func (g *Hello) Draw() {
	g.batch.Begin()
	eng.DefaultFont().Print(g.batch, "Hello, world!", 430, 280, nil)
	g.batch.End()
}

func main() {
	eng.Run(new(Hello))
}
