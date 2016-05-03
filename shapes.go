package engo

import (
	"engo.io/gl"
)

type TriangleType uint8

const (
	// Indicates a Triangle where two sides have equal length
	TriangleIsosceles TriangleType = iota
	// Indicates a Triangles where one angle is at 90 degrees
	TriangleRight
)

// Triangle is a basic triangular form; the "point" of the triangle is pointing to the top
type Triangle struct {
	TriangleType TriangleType
}

func (Triangle) Texture() *gl.Texture                       { return nil }
func (Triangle) Width() float32                             { return 0 }
func (Triangle) Height() float32                            { return 0 }
func (Triangle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Triangle) Close()                                     {}

// Rectangle is a basic rectangular form; the dimensions are controlled via the `SpaceComponent`.
type Rectangle struct{}

func (Rectangle) Texture() *gl.Texture                       { return nil }
func (Rectangle) Width() float32                             { return 0 }
func (Rectangle) Height() float32                            { return 0 }
func (Rectangle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Rectangle) Close()                                     {}

// Circle is a basic circular form; the dimensions / radius are controlled via the `SpaceComponent`.
// This was made possible by the shared knowledge of Olivier Gagnon (@hydroflame).
type Circle struct{}

func (Circle) Texture() *gl.Texture                       { return nil }
func (Circle) Width() float32                             { return 0 }
func (Circle) Height() float32                            { return 0 }
func (Circle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Circle) Close()                                     {}
