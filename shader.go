// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

import (
	gl "github.com/chsc/gogl/gl33"
)

type Shader struct {
	id          gl.Uint
	InPosition  gl.Uint
	InColor     gl.Uint
	InTexCoords gl.Uint
	UfMatrix    gl.Int
}

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

	inPosition := gl.GLString("in_Position")
	defer gl.GLStringFree(inPosition)
	gl.BindAttribLocation(program, 0, inPosition)

	inColor := gl.GLString("in_Color")
	defer gl.GLStringFree(inColor)
	gl.BindAttribLocation(program, 1, inColor)

	inTexCoords := gl.GLString("in_TexCoords")
	defer gl.GLStringFree(inTexCoords)
	gl.BindAttribLocation(program, 2, inTexCoords)

	gl.LinkProgram(program)

	var link_status gl.Int
	gl.GetProgramiv(program, gl.LINK_STATUS, &link_status)
	if link_status == 0 {
		panic("Unable to link shader program.\n")
	}

	matrix := gl.GLString("uf_Matrix")
	defer gl.GLStringFree(matrix)
	ufMatrix := gl.GetUniformLocation(program, matrix)

	return &Shader{program, 0, 1, 2, ufMatrix}
}

func (s *Shader) Bind() {
	gl.UseProgram(s.id)
}

func (s *Shader) Get(attrib string) gl.Int {
	return gl.GetUniformLocation(s.id, gl.GLString(attrib))
}
