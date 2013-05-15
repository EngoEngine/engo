// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fx

import (
	"github.com/ajhager/eng"
	gl "github.com/chsc/gogl/gl33"
)

const filmVert = `
attribute vec4 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform mat4 uf_Matrix;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
  gl_Position = uf_Matrix * in_Position;
}
`

const filmFrag = `
varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

// control parameter
uniform float time;
uniform bool grayscale;
// noise effect intensity value (0 = no effect, 1 = full effect)
uniform float nIntensity;
// scanlines effect intensity value (0 = no effect, 1 = full effect)
uniform float sIntensity;
// scanlines effect count value (0 = no effect, 4096 = full effect)
uniform float sCount;

void main (void) {
  // sample the source
  vec4 cTextureScreen = texture2D(uf_Texture, var_TexCoords);
  // make some noise
  float x = var_TexCoords.x * var_TexCoords.y * time *  1000.0;
  x = mod( x, 13.0 ) * mod( x, 123.0 );
  float dx = mod( x, 0.01 );

  // add noise
  vec3 cResult = cTextureScreen.rgb + cTextureScreen.rgb * clamp( 0.1 + dx * 100.0, 0.0, 1.0 );
  // get us a sine and cosine
  vec2 sc = vec2( sin( var_TexCoords.y * sCount), cos( var_TexCoords.y * sCount ) );

  // add scanlines
  cResult += cTextureScreen.rgb * vec3( sc.x, sc.y, sc.x ) * sIntensity;

  // interpolate between source and result by intensity
  cResult = cTextureScreen.rgb + clamp( nIntensity, 0.0,1.0 ) * ( cResult - cTextureScreen.rgb );

  // convert to grayscale if desired
  if( grayscale ) {
    cResult = vec3( cResult.r * 0.3 + cResult.g * 0.59 + cResult.b * 0.11 );
  }

  gl_FragColor = var_Color * vec4( cResult, cTextureScreen.a );
}
`

type Film struct {
	nIntensity   float32
	sIntensity   float32
	sCount       float32
	grayscale    int
	time         float32
	shader       *eng.Shader
	ufTime       gl.Int
	ufNIntensity gl.Int
	ufSIntensity gl.Int
	ufSCount     gl.Int
	ufGrayscale  gl.Int
}

func DefaultFilm() *Film {
	return NewFilm(.5, .05, 1024, false)
}

// NewFilm returns an effect that produces noise and scanlines when
// rendering. nIntensity is the intensity of the noise and should be a
// number between 0 and 1. sIntensity is the intensity of the
// scanlines and should be a number between 0 and 1. sCount is the
// number of scanlines. grayscale is whether or not to turn everything
// rendered black and white.
func NewFilm(nIntensity, sIntensity, sCount float32, grayscale bool) *Film {
	film := new(Film)
	film.nIntensity = nIntensity
	film.sIntensity = sIntensity
	film.sCount = sCount
	if grayscale {
		film.grayscale = 1
	}
	film.shader = eng.NewShader(filmVert, filmFrag)
	film.ufTime = film.shader.GetUniform("time")
	film.ufNIntensity = film.shader.GetUniform("nIntensity")
	film.ufSIntensity = film.shader.GetUniform("sIntensity")
	film.ufSCount = film.shader.GetUniform("sCount")
	film.ufGrayscale = film.shader.GetUniform("grayscale")
	return film
}

// Shader returns the underlying shader of the effect.
func (f *Film) Shader() *eng.Shader {
	return f.shader
}

// Setup binds the uniform values need to run the effect.
func (f *Film) Setup() {
	f.time += eng.Dt()
	gl.Uniform1f(f.ufTime, gl.Float(f.time))
	gl.Uniform1f(f.ufNIntensity, gl.Float(f.nIntensity))
	gl.Uniform1f(f.ufSIntensity, gl.Float(f.sIntensity))
	gl.Uniform1f(f.ufSCount, gl.Float(f.sCount))
	gl.Uniform1i(f.ufGrayscale, gl.Int(f.grayscale))
}
