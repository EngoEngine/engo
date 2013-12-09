// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rog

// State constants are arbitrary, but mapper functions in this
// module use the following by convention.
var (
	EMPTY = 0
	WALL  = 1
	FLOOR = 2
)

// Types that implement the Mapper interface can have the state
// at a coordinate in a 2d field set.
type Mapper interface {
	Map(x, y, state int)
}

// SparseMap reprsents a 2d field using a map.
type SparseMap struct {
	Width, Height int
	items         map[int]int
}

func NewSparseMap(width, height int) *SparseMap {
	return &SparseMap{width, height, make(map[int]int)}
}

func (m *SparseMap) Map(x, y int, item int) {
	y = y % m.Height
	if y < 0 {
		y = m.Height - y
	}
	x = x % m.Width
	if x < 0 {
		x = m.Width - x
	}
	if item == 0 {
		delete(m.items, y*m.Width+x)
	} else {
		m.items[y*m.Width+x] = item
	}
}

func (m *SparseMap) Get(x, y int) int {
	y = y % m.Height
	if y < 0 {
		y = m.Height - y
	}
	x = x % m.Width
	if x < 0 {
		x = m.Width - x
	}
	item, _ := m.items[y*m.Width+x]
	return item
}

// Maps out a rectangular room.
func MapArena(x, y, width, height int, m Mapper) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if i != x && j != y && i != x+width-1 && j != y+height-1 {
				m.Map(i, j, FLOOR)
			} else {
				m.Map(i, j, WALL)
			}
		}
	}
}
