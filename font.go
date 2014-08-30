// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

type glyph struct {
	region   *Region
	xoffset  float32
	yoffset  float32
	xadvance float32
}

type Font struct {
	glyphs map[rune]*glyph
}

func NewGridFont(texture *Texture, cellWidth, cellHeight int) *Font {
	i := 0
	glyphs := make(map[rune]*glyph)

	for y := 0; y < int(texture.Height())/cellHeight; y++ {
		for x := 0; x < int(texture.Width())/cellWidth; x++ {
			g := &glyph{xadvance: float32(cellWidth)}
			g.region = NewRegion(texture, x*cellWidth, y*cellHeight, cellWidth, cellHeight)
			glyphs[rune(i)] = g
			i += 1
		}
	}

	return &Font{glyphs}
}

func (f *Font) Remap(mapping string) {
	glyphs := make(map[rune]*glyph)

	i := 0
	for _, v := range mapping {
		glyphs[v] = f.glyphs[rune(i)]
		i++
	}

	f.glyphs = glyphs
}

func (f *Font) Put(batch *Batch, r rune, x, y float32, color uint32) {
	if g, ok := f.glyphs[r]; ok {
		batch.Draw(g.region, x+g.xoffset, y+g.yoffset, 0, 0, 1, 1, 0, color, 1)
	}
}

func (f *Font) Print(batch *Batch, text string, x, y float32, color uint32) {
	xx := x
	for _, r := range text {
		if g, ok := f.glyphs[r]; ok {
			batch.Draw(g.region, xx+g.xoffset, y+g.yoffset, 0, 0, 1, 1, 0, color, 1)
			xx += g.xadvance
		}
	}
}
