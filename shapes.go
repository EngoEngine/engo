package engo

import (
	"engo.io/gl"
)

type Triangle struct{}

func (Triangle) Texture() *gl.Texture                       { return nil }
func (Triangle) Width() float32                             { return 0 }
func (Triangle) Height() float32                            { return 0 }
func (Triangle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Triangle) Close()                                     {}

type Rectangle struct{}

func (Rectangle) Texture() *gl.Texture                       { return nil }
func (Rectangle) Width() float32                             { return 0 }
func (Rectangle) Height() float32                            { return 0 }
func (Rectangle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Rectangle) Close()                                     {}

type Circle struct{}

func (Circle) Texture() *gl.Texture                       { return nil }
func (Circle) Width() float32                             { return 0 }
func (Circle) Height() float32                            { return 0 }
func (Circle) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Circle) Close()                                     {}
