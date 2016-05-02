package engo

import (
	"image/color"

	"engo.io/gl"
)

type Triangle struct {
	Color color.Color
}

func (Triangle) Texture() *gl.Texture                       { return nil }
func (Triangle) Width() float32                             { return 0 }
func (Triangle) Height() float32                            { return 0 }
func (Triangle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Triangle) Close()                                     {}
