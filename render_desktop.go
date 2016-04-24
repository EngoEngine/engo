// +build !netgo

package engo

type Drawable interface {
	Texture() *gl.Texture
	Width() float32
	Height() float32
	View() (float32, float32, float32, float32)
}

type RenderComponent struct {
	// Hidden is used to prevent drawing by OpenGL
	Hidden bool

	// Transparency is the level of transparency that is used to draw the texture
	Transparency float32

	scale  Point
	Color  color.Color
	shader Shader
	zIndex float32

	drawable      Drawable
	buffer        *gl.Buffer
	bufferContent []float32
}
