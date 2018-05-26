//+build !appveyor,!jstesting

package common

import (
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/stretchr/testify/assert"
)

type testScene struct{}

func (*testScene) Preload() {}

func (t *testScene) Setup(u engo.Updater) {}

func (*testScene) Type() string { return "testScene" }

// TestShadersInitialization tests whether all registered `Shader`s will `Setup` without any errors
func TestShadersInitialization(t *testing.T) {
	engo.Run(engo.RunOptions{
		NoRun:        true,
		HeadlessMode: true,
	}, &testScene{})
	engo.CreateWindow("", 100, 100, false, 1)
	defer engo.DestroyWindow()

	w := &ecs.World{}
	w.AddSystem(&CameraSystem{})
	err := initShaders(w)
	assert.NoError(t, err)
}

// TestShaderCompilation tests whether the `LoadShader` method will indeed report errors iff
// (one of) the GLSL-shaders is incorrect.
func TestShaderCompilation(t *testing.T) {
	engo.Run(engo.RunOptions{
		NoRun:        true,
		HeadlessMode: true,
	}, &testScene{})
	engo.CreateWindow("", 100, 100, false, 1)
	defer engo.DestroyWindow()

	var err error

	_, err = LoadShader(correctVertShader, correctFragShader)
	assert.NoError(t, err)

	_, err = LoadShader(correctVertShader, incorrectFragShader)
	assert.IsType(t, FragmentShaderCompilationError{}, err)

	_, err = LoadShader(incorrectVertShader, correctFragShader)
	assert.IsType(t, VertexShaderCompilationError{}, err)

	_, err = LoadShader(incorrectVertShader, incorrectFragShader)
	assert.Error(t, err) // don't really care which one it is
}

var correctVertShader = `
attribute vec2 in_Position;
attribute vec4 in_Color;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;

void main() {

  vec3 matr = matrixProjection * matrixView * matrixModel * vec3(in_Position, 1.0);
  gl_Position = vec4(matr.xy, 0, matr.z);
}
`
var correctFragShader = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

void main (void) {
  gl_FragColor = vec4(1.0); // all-white
}
`
var incorrectVertShader = `
this is incorrect GLSL syntax
`
var incorrectFragShader = `
this is incorrect GLSL syntax
`
