// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fx

import (
	"github.com/ajhager/eng"
	gl "github.com/chsc/gogl/gl32"
)

const vignetteFrag = `
uniform float size;
uniform float amount;

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main() {
  vec4 color = texture2D(uf_Texture, var_TexCoords);

  float dist = distance(var_TexCoords, vec2(0.5, 0.5));
  color.rgb *= smoothstep(0.8, size*0.799, dist * (amount + size));

  gl_FragColor = var_Color * color;
}
`

type Vignette struct {
	Size     float32
	Amount   float32
	shader   *eng.Shader
	ufSize   gl.Int
	ufAmount gl.Int
}

func NewVignette(size, amount float32) *Vignette {
	vignette := new(Vignette)
	vignette.Size = size
	vignette.Amount = amount
	vignette.shader = eng.NewShader(filmVert, vignetteFrag)
	vignette.ufSize = vignette.shader.GetUniform("size")
	vignette.ufAmount = vignette.shader.GetUniform("amount")
	return vignette
}

// Shader returns the underlying shader of the effect.
func (v *Vignette) Shader() *eng.Shader {
	return v.shader
}

// Setup binds the uniform values need to run the effect.
func (v *Vignette) Setup() {
	gl.Uniform1f(v.ufSize, gl.Float(v.Size))
	gl.Uniform1f(v.ufAmount, gl.Float(v.Amount))
}
