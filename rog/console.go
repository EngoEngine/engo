package rog

import (
	"fmt"
	"github.com/ajhager/eng"
)

type Console struct {
	bg, fg [][]*eng.Color
	ch     [][]rune
	w, h   int
}

func NewConsole(width, height int) *Console {
	bg := make([][]*eng.Color, height)
	fg := make([][]*eng.Color, height)
	ch := make([][]rune, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]*eng.Color, width)
		fg[y] = make([]*eng.Color, width)
		ch[y] = make([]rune, width)
	}

	con := &Console{bg, fg, ch, width, height}

	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.bg[y][x] = eng.NewColor(0, 0, 0, 1)
			con.fg[y][x] = eng.NewColor(1, 1, 1, 1)
			con.ch[y][x] = ' '
		}
	}

	return con
}

func (con *Console) put(x, y, i, t int, fg, bg eng.Blender, ch rune) {
	if x < 0 || x >= con.w || y < 0 || y >= con.h {
		return
	}

	if ch > 0 {
		con.ch[y][x] = ch
	}

	if bg != nil {
		con.bg[y][x] = bg.Blend(con.bg[y][x], i, t)
	}

	if fg != nil {
		con.fg[y][x] = fg.Blend(con.bg[y][x], i, t)
	}
}

func (con *Console) set(i, j, x, y, w, h int, fg, bg eng.Blender, data string, rest ...interface{}) {
	if len(rest) > 0 {
		data = fmt.Sprintf(data, rest...)
	}
	t := len(data)
	if t > 0 {
		if h == 0 {
			h = con.h - y
		}
		for k, r := range data {
			if i == x+w {
				j += 1
				i = x
			}
			if j == y+h {
				break
			}
			con.put(i, j, k, t, fg, bg, r)
			i += 1
		}
	} else {
		con.put(i, j, 0, 0, fg, bg, -1)
	}
}

// Blit draws con onto this console with top left starting at x, y.
func (con *Console) Blit(o *Console, x, y int) {
	for i := 0; i < o.Width(); i++ {
		for j := 0; j < o.Height(); j++ {
			fg, bg, ch := o.Get(i, j)
			con.Set(x+i, y+j, fg, bg, string(ch))
		}
	}
}

// Clear is a short hand to fill the entire screen with the given colors and rune.
func (con *Console) Clear(fg, bg eng.Blender, ch rune) {
	con.Fill(0, 0, con.w, con.h, fg, bg, ch)
}

// Fill draws a rect on the root console using ch.
func (con *Console) Fill(x, y, w, h int, fg, bg eng.Blender, ch rune) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			con.Set(i, j, fg, bg, string(ch))
		}
	}
}

// Get returns the fg, bg colors and rune of the cell.
func (con *Console) Get(x, y int) (*eng.Color, *eng.Color, rune) {
	return con.fg[y][x], con.bg[y][x], con.ch[y][x]
}

// Height returns the height of the console in cells.
func (con *Console) Height() int {
	return con.h
}

// Set draws a string starting at x,y onto the console, wrapping at the bounds if needed.
func (con *Console) Set(x, y int, fg, bg eng.Blender, data string, rest ...interface{}) {
	con.set(x, y, 0, 0, con.w, con.h, fg, bg, data, rest...)
}

// SetRect draws a string starting at x,y onto the console, wrapping at the bounds created by x, y, w, h if needed.
// If h is 0, the text will cut off at the bottom of the console, otherwise it will cut off after the y+h row.
func (con *Console) SetRect(x, y, w, h int, fg, bg eng.Blender, data string, rest ...interface{}) {
	con.set(x, y, x, y, w, h, fg, bg, data, rest...)
}

// Render takes a batch and font data, and renders the console using them.
func (con *Console) Render(batch *eng.Batch, font *eng.Font, w, h int) {
	for x := 0; x < con.Width(); x++ {
		for y := 0; y < con.Height(); y++ {
			fg, bg, ch := con.Get(x, y)
			font.Print(batch, "█", float32(x*w), float32(y*h), bg)
			font.Print(batch, fmt.Sprintf("%c", ch), float32(x*w), float32(y*h), fg)
		}
	}
}

// Width returns the width of the console in cells.
func (con *Console) Width() int {
	return con.w
}
