package common

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi = float64(72)
)

var fontCache []*Font

// Font keeps track of a specific Font. Fonts are explicit instances of a font file,
// including the Size and Color. A separate font will have to be generated to get
// different sizes and colors of the same font file.
type Font struct {
	URL  string
	Size float64
	BG   color.Color
	FG   color.Color
	TTF  *truetype.Font
	face font.Face
}

// LoadedFont returns a Font that was previously loaded via CreatePreloaded
func LoadedFont(url string, size float64, bg, fg color.Color) (*Font, error) {
	idx := -1
	for i, fnt := range fontCache {
		if fnt.URL == url && fnt.Size == size && fnt.BG == bg && fnt.FG == fg {
			idx = i
		}
	}
	if idx < 0 {
		return nil, errors.New("No font matching that description was found.")
	}
	return fontCache[idx], nil
}

// Create is for loading fonts from the disk, given a location
func (f *Font) Create() error {
	// Read and parse the font
	ttfBytes, err := ioutil.ReadFile(f.URL)
	if err != nil {
		return err
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return err
	}
	f.TTF = ttf
	f.face = truetype.NewFace(f.TTF, &truetype.Options{
		Size:    f.Size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	return nil
}

// CreatePreloaded is for loading fonts which have already been defined (and loaded) within Preload
func (f *Font) CreatePreloaded() error {
	fontres, err := engo.Files.Resource(f.URL)
	if err != nil {
		return err
	}

	fnt, ok := fontres.(FontResource)
	if !ok {
		return fmt.Errorf("preloaded font is not of type `*truetype.Font`: %s", f.URL)
	}

	f.TTF = fnt.Font
	f.face = truetype.NewFace(f.TTF, &truetype.Options{
		Size:    f.Size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	fontCache = append(fontCache, f)
	return nil
}

// TextDimensions returns the total width, total height and total line size
// of the input string written out in the Font.
func (f *Font) TextDimensions(text string) (int, int, int) {
	fnt := f.TTF
	size := f.Size
	var (
		totalWidth  = fixed.Int26_6(0)
		totalHeight = fixed.Int26_6(size)
		maxYBearing = fixed.Int26_6(0)
	)
	fupe := fixed.Int26_6(fnt.FUnitsPerEm())

	for _, char := range text {
		idx := fnt.Index(char)
		hm := fnt.HMetric(fupe, idx)
		vm := fnt.VMetric(fupe, idx)
		g := truetype.GlyphBuf{}
		err := g.Load(fnt, fupe, idx, font.HintingNone)
		if err != nil {
			log.Println(err)
			return 0, 0, 0
		}
		totalWidth += hm.AdvanceWidth

		yB := (vm.TopSideBearing * fixed.Int26_6(size)) / fupe
		if yB > maxYBearing {
			maxYBearing = yB
		}
		dY := (vm.AdvanceHeight * fixed.Int26_6(size)) / fupe
		if dY > totalHeight {
			totalHeight = dY
		}
	}

	// Scale to actual pixel size
	totalWidth *= fixed.Int26_6(size)
	totalWidth /= fupe

	return int(totalWidth), int(totalHeight), int(maxYBearing)
}

// RenderNRGBA returns an *image.NRGBA in the Font based on the input string.
func (f *Font) RenderNRGBA(text string) *image.NRGBA {
	width, height, yBearing := f.TextDimensions(text)
	font := f.TTF
	size := f.Size

	if size <= 0 {
		panic("Font size cannot be <= 0")
	}

	// Default colors
	if f.FG == nil {
		f.FG = color.NRGBA{0, 0, 0, 0}
	}
	if f.BG == nil {
		f.BG = color.NRGBA{0, 0, 0, 0}
	}

	// Colors
	fg := image.NewUniform(f.FG)
	bg := image.NewUniform(f.BG)

	// Create the font context
	c := freetype.NewContext()

	nrgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(nrgba, nrgba.Bounds(), bg, image.ZP, draw.Src)

	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(nrgba.Bounds())
	c.SetDst(nrgba)
	c.SetSrc(fg)

	// Draw the text.
	pt := fixed.P(0, yBearing)
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return nil
	}

	return nrgba
}

// Render returns a Texture in the Font based on the input string.
func (f *Font) Render(text string) Texture {
	nrgba := f.RenderNRGBA(text)

	// Create texture
	imObj := NewImageObject(nrgba)
	return NewTextureSingle(imObj)
}

// generateFontAtlas generates the font atlas for this given font, using the first `c` Unicode characters.
func (f *Font) generateFontAtlas(c int) FontAtlas {
	atlas := FontAtlas{
		XLocation: make([]float32, c),
		YLocation: make([]float32, c),
		Width:     make([]float32, c),
		Height:    make([]float32, c),
		OffsetX:   make([]float32, c),
		RightSide: make([]float32, c),
		OffsetY:   make([]float32, c),
	}

	currentX := float32(0)
	currentY := float32(0)

	// Default colors
	if f.FG == nil {
		f.FG = color.NRGBA{0, 0, 0, 0}
	}
	if f.BG == nil {
		f.BG = color.NRGBA{0, 0, 0, 0}
	}

	d := &font.Drawer{}
	d.Src = image.NewUniform(f.FG)
	d.Face = truetype.NewFace(f.TTF, &truetype.Options{
		Size:    f.Size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})

	lineHeight := d.Face.Metrics().Height
	ascent := d.Face.Metrics().Ascent
	prev := 0

	for i := 0; i < c; i++ {
		bounds, adv, ok := d.Face.GlyphBounds(rune(i))
		if !ok {
			continue
		}
		advance := float32(adv.Ceil())

		atlas.Width[i] = float32((bounds.Max.X - bounds.Min.X).Ceil())
		atlas.Height[i] = float32((bounds.Max.Y - bounds.Min.Y).Ceil())
		atlas.OffsetX[i] = float32(bounds.Min.X.Ceil())
		atlas.OffsetY[i] = float32((ascent + bounds.Min.Y).Ceil())
		atlas.RightSide[i] = advance - float32(bounds.Max.X.Ceil())

		if prev > 0 {
			currentX += float32(f.face.Kern(rune(prev), rune(i)).Ceil())
		}

		//overlapping characters
		if atlas.Width[i] > advance {
			if atlas.OffsetX[i] < 0 {
				advance -= atlas.OffsetX[i]
			} else if advance < float32(bounds.Max.X.Ceil()) {
				advance = float32(bounds.Max.X.Ceil())
			}
		}

		// position correction of overlapped characters
		if atlas.OffsetX[i] < 0 {
			currentX -= atlas.OffsetX[i]
		}

		if currentX+advance > 1024 {
			currentX = 0
			currentY += float32(lineHeight.Ceil())
			atlas.TotalHeight += float32(lineHeight.Ceil())
			prev = 0
		}

		if currentX+advance > atlas.TotalWidth {
			atlas.TotalWidth = currentX + advance
		}

		atlas.XLocation[i] = currentX
		atlas.YLocation[i] = currentY
		currentX += advance
		prev = i
	}

	// Create texture
	actual := image.NewNRGBA(image.Rect(0, 0, int(atlas.TotalWidth), int(atlas.TotalHeight)))
	draw.Draw(actual, actual.Bounds(), image.NewUniform(f.BG), image.ZP, draw.Src)
	d.Dst = actual

	for i := 0; i < c; i++ {
		_, _, ok := d.Face.GlyphBounds(rune(i))
		if !ok {
			continue
		}
		d.Dot = fixed.P(int(atlas.XLocation[i]), int(atlas.YLocation[i]+float32(ascent.Ceil())))
		d.DrawBytes([]byte{byte(i)})
		// position correction
		atlas.XLocation[i] += atlas.OffsetX[i]
		atlas.YLocation[i] += atlas.OffsetY[i]
	}

	imObj := NewImageObject(actual)
	atlas.Texture = NewTextureSingle(imObj).id

	return atlas
}

// GenerateFontAtlas generates the font atlas for this given font, using the first `c` Unicode characters.
// This should only be used if you are writing your own custom text shader.
func (f *Font) GenerateFontAtlas(c int) FontAtlas {
	return f.generateFontAtlas(c)
}

// A FontAtlas is a representation of some of the Font characters, as an image
type FontAtlas struct {
	Texture *gl.Texture
	// XLocation contains the X-coordinate of the starting position of all characters
	XLocation []float32
	// YLocation contains the Y-coordinate of the starting position of all characters
	YLocation []float32
	// Width contains the width in pixels of all the characters, including the spacing between characters
	Width []float32
	// Height contains the height in pixels of all the characters
	Height []float32
	// OffsetX is the offset of the glyph and the x-coordinate
	OffsetX []float32
	// RightSide is the space left to the right of the glyph
	RightSide []float32
	// OffsetY is the offset of the glyph and the y-coordinate
	OffsetY []float32
	// TotalWidth is the total amount of pixels the `FontAtlas` is wide; useful for determining the `Viewport`,
	// which is relative to this value.
	TotalWidth float32
	// TotalHeight is the total amount of pixels the `FontAtlas` is high; useful for determining the `Viewport`,
	// which is relative to this value.
	TotalHeight float32
}

// Text represents a string drawn onto the screen, as used by the `TextShader`.
type Text struct {
	// Font is the reference to the font you're using to render this. This includes the color, as well as the font size.
	Font *Font
	// Text is the actual text you want to draw. This may include newlines (\n).
	Text string
	// LineSpacing is the amount of additional spacing there is between the lines (when `Text` consists of multiple lines),
	// relative to the `Size` of the `Font`.
	LineSpacing float32
	// LetterSpacing is the amount of additional spacing there is between the characters, relative to the `Size` of
	// the `Font`.
	LetterSpacing float32
	// RightToLeft is an experimental variable used to indicate that subsequent characters come to the left of the
	// previous character.
	RightToLeft bool
}

// Texture returns nil because the Text is generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Texture() *gl.Texture { return nil }

// Width returns the width of the Text generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Width() float32 {
	atlas, ok := atlasCache[*t.Font]
	if !ok {
		// Generate texture first
		atlas = t.Font.generateFontAtlas(200)
		atlasCache[*t.Font] = atlas
	}

	var currentX float32
	var greatestX float32

	for _, char := range t.Text {
		// TODO: this might not work for all characters
		switch {
		case char == '\n':
			if currentX > greatestX {
				greatestX = currentX
			}
			currentX = 0
			continue
		case char < 32: // all system stuff should be ignored
			continue
		}

		currentX += atlas.Width[char] + float32(t.Font.Size)*t.LetterSpacing
	}
	if currentX > greatestX {
		return currentX
	}
	return greatestX
}

// Height returns the height the Text generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Height() float32 {
	atlas, ok := atlasCache[*t.Font]
	if !ok {
		// Generate texture first
		atlas = t.Font.generateFontAtlas(200)
		atlasCache[*t.Font] = atlas
	}

	var currentY float32
	var totalY float32
	var tallest float32

	for _, char := range t.Text {
		// TODO: this might not work for all characters
		switch {
		case char == '\n':
			totalY += tallest
			tallest = float32(0)
			continue
		case char < 32: // all system stuff should be ignored
			continue
		}
		currentY = atlas.Height[char] + t.LineSpacing*atlas.Height[char]
		if currentY > tallest {
			tallest = currentY
		}
	}
	return totalY + tallest
}

// View returns 0, 0, 1, 1 because the Text is generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }

// Close does nothing because the Text is generated from a FontAtlas. There is no underlying texture to close.
// This implements the common.Drawable interface.
func (t Text) Close() {}
