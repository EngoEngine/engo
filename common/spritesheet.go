package common

import (
	"log"

	"engo.io/engo"
	"engo.io/gl"
)

// Spritesheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Spritesheet struct {
	texture               *gl.Texture     // The original texture
	width, height         float32         // The dimensions of the total texture
	cellWidth, cellHeight int             // The dimensions of the cells
	cache                 map[int]Texture // The cell cache cells
}

func NewSpritesheetFromTexture(tr *TextureResource, cellWidth, cellHeight int) *Spritesheet {
	return &Spritesheet{texture: tr.Texture,
		width: tr.Width, height: tr.Height,
		cellWidth: cellWidth, cellHeight: cellHeight,
		cache: make(map[int]Texture),
	}
}

// NewSpritesheetFromFile is a simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engo.Files.Add
func NewSpritesheetFromFile(textureName string, cellWidth, cellHeight int) *Spritesheet {
	res, err := engo.Files.Resource(textureName)
	if err != nil {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Received error:", err)
		return nil
	}

	img, ok := res.(TextureResource)
	if !ok {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Resource not of type `TextureResource`:", textureName)
		return nil
	}

	return NewSpritesheetFromTexture(&img, cellWidth, cellHeight)
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(index int) Texture {
	if r, ok := s.cache[index]; ok {
		return r
	}

	cellsPerRow := int(s.Width())
	var x float32 = float32((index % cellsPerRow) * s.cellWidth)
	var y float32 = float32((index / cellsPerRow) * s.cellHeight)
	s.cache[index] = Texture{id: s.texture, width: float32(s.cellWidth), height: float32(s.cellHeight), viewport: engo.AABB{
		engo.Point{x / s.width, y / s.height},
		engo.Point{(x + float32(s.cellWidth)) / s.width, (y + float32(s.cellHeight)) / s.height},
	}}

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

func (s *Spritesheet) Cells() []Texture {
	cellsNo := s.CellCount()
	cells := make([]Texture, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
}

// Width is the amount of tiles on the x-axis of the spritesheet
func (s Spritesheet) Width() float32 {
	return s.width / float32(s.cellWidth)
}

// Height is the amount of tiles on the y-axis of the spritesheet
func (s Spritesheet) Height() float32 {
	return s.height / float32(s.cellHeight)
}

/*
type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float32
	Color    color.Color
	Alpha    float32
	Region   *Region
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Anchor:   &Point{0, 0},
		Rotation: 0,
		Color:    color.White,
		Alpha:    1,
		Region:   region,
	}
}
*/
