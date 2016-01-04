package engi

// Spritesheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Spritesheet struct {
	texture               *Texture        // The original texture
	CellWidth, CellHeight int             // The dimensions of the cells
	cache                 map[int]*Region // The cell cache cells
}

func NewSpritesheetFromTexture(texture *Texture, cellWidth, cellHeight int) *Spritesheet {
	return &Spritesheet{texture: texture, CellWidth: cellWidth, CellHeight: cellHeight, cache: make(map[int]*Region)}
}

// NewSpritesheetFromFile is a simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engi.Files.Add
func NewSpritesheetFromFile(textureName string, cellWidth, cellHeight int) *Spritesheet {
	return NewSpritesheetFromTexture(Files.Image(textureName), cellWidth, cellHeight)
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(index int) *Region {
	if r := s.cache[index]; r != nil {
		return r
	}
	s.cache[index] = regionFromSheet(s.texture, s.CellWidth, s.CellHeight, index)

	return s.cache[index]
}

func (s *Spritesheet) Drawable(index int) Drawable {
	return s.Cell(index)
}

func (s *Spritesheet) Drawables() []Drawable {
	drawables := make([]Drawable, s.CellCount())

	for i := 0; i < s.CellCount(); i++ {
		drawables[i] = s.Drawable(i)
	}

	return drawables
}

func (s *Spritesheet) CellCount() int {
	return int(s.Width()) * int(s.Height())
}

func (s *Spritesheet) Cells() []*Region {
	cellsNo := s.CellCount()
	cells := make([]*Region, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
}

// Width is the amount of tiles on the x-axis of the spritesheet
func (s Spritesheet) Width() float32 {
	return s.texture.Width() / float32(s.CellWidth)
}

// Height is the amount of tiles on the y-axis of the spritesheet
func (s Spritesheet) Height() float32 {
	return s.texture.Height() / float32(s.CellHeight)
}
