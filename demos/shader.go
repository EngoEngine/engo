package main

import (
	"github.com/ajhager/eng"
)

const vert = `
attribute vec4 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
	gl_Position = vec4(in_Position.x / uf_Projection.x - 1.0,
										 in_Position.y / -uf_Projection.y + 1.0,
										 0.0, 1.0);
}
`

const frag = `
varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = vec4(1, var_TexCoords.x, var_TexCoords.y, 1) * texture2D (uf_Texture, var_TexCoords);
}
`

var batch *eng.Batch

type Game struct {
	*eng.Game
}

func (g *Game) Setup() {
	batch = eng.NewBatch()
	batch.SetShader(eng.NewShader(vert, frag))
}

func (g *Game) Draw() {
	batch.Begin()
	eng.DefaultFont().Print(batch, "Hello, world!", 430, 280, nil)
	batch.End()
}

func main() {
	eng.Run("Shader", 1024, 640, false, new(Game))
}
