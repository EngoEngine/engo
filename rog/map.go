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
	states        map[int]int
}

func NewSparseMap(width, height int) *SparseMap {
	return &SparseMap{width, height, make(map[int]int)}
}

func (m *SparseMap) Map(x, y, state int) {
	y = y % m.Height
	if y < 0 {
		y = m.Height - y
	}
	x = x % m.Width
	if x < 0 {
		x = m.Width - x
	}
	if state == 0 {
		delete(m.states, y*m.Width+x)
	} else {
		m.states[y*m.Width+x] = state
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
	state, _ := m.states[y*m.Width+x]
	return state
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
