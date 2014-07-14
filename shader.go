// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	"github.com/errcw/glow/gl/2.1/gl"
	"log"
	"strings"
)

// A Shader abstracts the loading, compiling, and linking of shader
// programs, which can directly modify the rendering of vertices and pixels.
type Shader struct {
	id uint32
}

// NewShader takes the source of a vertex and fragment shader and
// returns a compiled and linked shader program.
func NewShader(vertSrc, fragSrc string) *Shader {
	glVertSrc := gl.Str(vertSrc + "\x00")
	vertShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertShader, 1, &glVertSrc, nil)
	gl.CompileShader(vertShader)
	defer gl.DeleteShader(vertShader)

	var status int32
	gl.GetShaderiv(vertShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertShader, gl.INFO_LOG_LENGTH, &logLength)

		logGL := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertShader, logLength, nil, gl.Str(logGL))

		log.Fatal("failed to compile %v: %v", vertSrc, logGL)
	}

	glFragSrc := gl.Str(fragSrc + "\x00")
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragShader, 1, &glFragSrc, nil)
	gl.CompileShader(fragShader)
	defer gl.DeleteShader(fragShader)

	gl.GetShaderiv(fragShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragShader, gl.INFO_LOG_LENGTH, &logLength)

		logGL := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragShader, logLength, nil, gl.Str(logGL))

		log.Fatal("failed to compile %v: %v", fragSrc, logGL)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)

	gl.LinkProgram(program)

	var linkStatus int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkStatus)
	if linkStatus == 0 {
		log.Fatal("Unable to link shader program.")
	}

	gl.ValidateProgram(program)

	var validateStatus int32
	gl.GetProgramiv(program, gl.VALIDATE_STATUS, &validateStatus)
	if validateStatus == 0 {
		//log.Fatal("Unable to validate shader program.")
	}

	return &Shader{program}
}

// Bind turns the shader on to be used during rendering.
func (s *Shader) Bind() {
	gl.UseProgram(s.id)
}

// GetUniform returns the location of the named uniform.
func (s *Shader) GetUniform(uniform string) int32 {
	glUniform := gl.Str(uniform + "\x00")
	return gl.GetUniformLocation(s.id, glUniform)
}

// GetAttrib returns the location of the named attribute.
func (s *Shader) GetAttrib(attrib string) uint32 {
	glAttrib := gl.Str(attrib + "\x00")
	return uint32(gl.GetAttribLocation(s.id, glAttrib))
}
