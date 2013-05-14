// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	gl "github.com/chsc/gogl/gl33"
	"log"
)

// A Shader abstracts the loading, compiling, and linking of shader
// programs, which can directly modify the rendering of vertices and pixels.
type Shader struct {
	id gl.Uint
}

// NewShader takes the source of a vertex and fragment shader and
// returns a compiled and linked shader program.
func NewShader(vertSrc, fragSrc string) *Shader {
	glVertSrc := gl.GLString(vertSrc)
	defer gl.GLStringFree(glVertSrc)
	vertShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertShader, 1, &glVertSrc, nil)
	gl.CompileShader(vertShader)
	defer gl.DeleteShader(vertShader)

	glFragSrc := gl.GLString(fragSrc)
	defer gl.GLStringFree(glFragSrc)
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragShader, 1, &glFragSrc, nil)
	gl.CompileShader(fragShader)
	defer gl.DeleteShader(fragShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)

	gl.LinkProgram(program)

	var link_status gl.Int
	gl.GetProgramiv(program, gl.LINK_STATUS, &link_status)
	if link_status == 0 {
		log.Fatal("Unable to link shader program.")
	}

	return &Shader{program}
}

// Bind turns the shader on to be used during rendering.
func (s *Shader) Bind() {
	gl.UseProgram(s.id)
}

// GetUniform returns the location of the named uniform.
func (s *Shader) GetUniform(uniform string) gl.Int {
	glUniform := gl.GLString(uniform)
	defer gl.GLStringFree(glUniform)
	return gl.GetUniformLocation(s.id, glUniform)
}

// GetAttrib returns the location of the named attribute.
func (s *Shader) GetAttrib(attrib string) gl.Uint {
	glAttrib := gl.GLString(attrib)
	defer gl.GLStringFree(glAttrib)
	return gl.Uint(gl.GetAttribLocation(s.id, glAttrib))
}
