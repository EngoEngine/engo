package common

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"engo.io/engo"
	"engo.io/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi = float64(72)
)

type Font struct {
	URL  string
	Size float64
	BG   color.Color
	FG   color.Color
	TTF  *truetype.Font
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

	return nil
}

// CreatePreloaded is for loading fonts which have already been defined (and loaded) within Preload
func (f *Font) CreatePreloaded() error {
	fontres, err := engo.Files.Resource(f.URL)
	if err != nil {
		return err
	}

	font, ok := fontres.(FontResource)
	if !ok {
		return fmt.Errorf("preloaded font is not of type `*truetype.Font`: %s", f.URL)
	}

	f.TTF = font.Font
	return nil
}

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
	}

	// Scale to actual pixel size
	totalWidth *= fixed.Int26_6(size)
	totalWidth /= fupe

	return int(totalWidth), int(totalHeight), int(maxYBearing)
}

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
	}

	var (
		int26Width  fixed.Int26_6
		int26Height fixed.Int26_6

		totalHeight fixed.Int26_6

		fupe        = fixed.Int26_6(f.TTF.FUnitsPerEm())
		maxYBearing = fixed.Int26_6(0)

		currentX float32
		maxX     float32

		totalString string // TODO; string is immutable, so this is relatively inefficient; see also where we use `totalString += string(char)`
		subString   string

		drawCurY int
	)

	// The "full image"
	nrgba := image.NewNRGBA(image.Rect(0, 0, 1200, 5000)) // way too big; hopefully

	for i := 0; i < c; i++ {
		char := rune(i)
		totalString += string(char)
		subString += string(char)
		atlas.XLocation[char] = currentX
		atlas.YLocation[char] = float32(drawCurY)

		idx := f.TTF.Index(char)
		hm := f.TTF.HMetric(fupe, idx)
		vm := f.TTF.VMetric(fupe, idx)

		g := truetype.GlyphBuf{}
		err := g.Load(f.TTF, fupe, idx, font.HintingNone)
		if err != nil {
			log.Println("Error creating font atlas:", err)
			return atlas
		}

		int26Width += hm.AdvanceWidth

		atlas.Width[char] = float32(g.AdvanceWidth * fixed.Int26_6(f.Size) / fupe)

		currentX = float32(int26Width * fixed.Int26_6(f.Size) / fupe)
		if currentX > maxX {
			maxX = currentX
		}

		yB := vm.TopSideBearing - g.Bounds.Min.Y
		if yB > maxYBearing {
			maxYBearing = yB
		}
		atlas.Height[char] = float32(yB * fixed.Int26_6(f.Size) / fupe)

		if int(int26Width*fixed.Int26_6(f.Size)/fupe) > 1024 {
			// Now let's draw these chars!
			subimg := f.RenderNRGBA(subString)
			// TODO: optimize this, because `f.RenderNRGBA` also allocates a new image, while we could be drawing to this one directly.
			draw.Draw(nrgba, image.Rect(0, drawCurY, subimg.Bounds().Max.X, drawCurY+subimg.Bounds().Max.Y), subimg, image.ZP, draw.Src)
			drawCurY += subimg.Bounds().Max.Y

			int26Height += maxYBearing
			totalHeight += maxYBearing
			maxYBearing = fixed.Int26_6(0)

			subString = ""
			int26Width = 0
			currentX = 0
		}
	}

	// Draw the last line as well
	subimg := f.RenderNRGBA(subString)
	// TODO: optimize this, because `f.RenderNRGBA` also allocates a new image, while we could be drawing to this one directly.
	draw.Draw(nrgba, image.Rect(0, drawCurY, subimg.Bounds().Max.X, drawCurY+subimg.Bounds().Max.Y), subimg, image.ZP, draw.Src)
	drawCurY += subimg.Bounds().Max.Y

	int26Height += maxYBearing
	totalHeight += maxYBearing
	maxYBearing = fixed.Int26_6(0)

	atlas.TotalWidth = maxX
	atlas.TotalHeight = float32(drawCurY)

	// Create texture
	actual := image.NewNRGBA(image.Rect(0, 0, int(atlas.TotalWidth), int(atlas.TotalHeight)))
	draw.Draw(actual, actual.Bounds(), nrgba, image.ZP, draw.Src)

	imObj := NewImageObject(actual)
	atlas.Texture = NewTextureSingle(imObj).id
	return atlas
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

func (Text) Texture() *gl.Texture                       { return nil }
func (Text) Width() float32                             { return 0 }
func (Text) Height() float32                            { return 0 }
func (Text) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (Text) Close()                                     {}
