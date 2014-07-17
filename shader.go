// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

// A Shader abstracts the loading, compiling, and linking of shader
// programs, which can directly modify the rendering of vertices and pixels.
type Shader struct {
	id *ProgramObject
}

// NewShader takes the source of a vertex and fragment shader and
// returns a compiled and linked shader program.
func NewShader(vertSrc, fragSrc string) *Shader {
	vertShader := GL.CreateShader(GL.VERTEX_SHADER)
	GL.ShaderSource(vertShader, vertSrc)
	GL.CompileShader(vertShader)
	defer GL.DeleteShader(vertShader)

	/*
		var status int32
		gl.GetShaderiv(vertShader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(vertShader, gl.INFO_LOG_LENGTH, &logLength)

			logGL := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(vertShader, logLength, nil, gl.Str(logGL))

			log.Fatal("failed to compile %v: %v", vertSrc, logGL)
		}
	*/

	fragShader := GL.CreateShader(GL.FRAGMENT_SHADER)
	GL.ShaderSource(fragShader, fragSrc)
	GL.CompileShader(fragShader)
	defer GL.DeleteShader(fragShader)

	/*
		gl.GetShaderiv(fragShader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(fragShader, gl.INFO_LOG_LENGTH, &logLength)

			logGL := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(fragShader, logLength, nil, gl.Str(logGL))

			log.Fatal("failed to compile %v: %v", fragSrc, logGL)
		}
	*/

	program := GL.CreateProgram()
	GL.AttachShader(program, vertShader)
	GL.AttachShader(program, fragShader)

	GL.LinkProgram(program)

	/*
		var linkStatus int32
		gl.GetProgramiv(program, gl.LINK_STATUS, &linkStatus)
		if linkStatus == 0 {
			log.Fatal("Unable to link shader program.")
		}

		gl.ValidateProgram(program)

		var validateStatus int32
		gl.GetProgramiv(program, gl.VALIDATE_STATUS, &validateStatus)
		if validateStatus == 0 {
			log.Fatal("Unable to validate shader program.")
		}
	*/

	return &Shader{program}
}

// Bind turns the shader on to be used during rendering.
func (s *Shader) Bind() {
	GL.UseProgram(s.id)
}

// GetUniform returns the location of the named uniform.
func (s *Shader) GetUniform(uniform string) *UniformObject {
	return GL.GetUniformLocation(s.id, uniform)
}

// GetAttrib returns the location of the named attribute.
func (s *Shader) GetAttrib(attrib string) int {
	return GL.GetAttribLocation(s.id, attrib)
}
