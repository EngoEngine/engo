package engi

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi = float64(72)
)

// TODO FG and BG color config
type Font struct {
	URL  string
	Size float64
	BG   color.Color
	FG   color.Color
	ttf  *truetype.Font
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
	f.ttf = ttf

	return nil
}

// CreatePreloaded is for loading fonts which have already been defined (and loaded) within Preload
func (f *Font) CreatePreloaded() error {
	var ok bool
	f.ttf, ok = Files.fonts[f.URL]
	if !ok {
		return fmt.Errorf("could not find preloaded font: %s", f.URL)
	}

	return nil
}

func (f *Font) TextDimensions(text string) (int, int, int) {
	fnt := f.ttf
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

func (f *Font) Render(text string) *Texture {
	width, height, yBearing := f.TextDimensions(text)
	font := f.ttf
	size := f.Size

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
	pt := freetype.Pt(0, int(yBearing))
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Create texture
	imObj := &ImageObject{nrgba}
	return NewTexture(imObj)

}
