package engi

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
)

var (
	dpi = float64(72)
)

type Color struct {
	R, G, B, A uint8
}

// TODO FG and BG color config
type Font struct {
	URL  string
	Size float64
	BG   Color
	FG   Color
	ttf  *truetype.Font
}

func (f *Font) Create() {
	url := f.URL

	// Read and parse the font
	ttfBytes, err := ioutil.ReadFile(url)
	if err != nil {
		log.Println(err)
		return
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		log.Println(err)
		return
	}
	f.ttf = ttf
}

func (f *Font) TextDimensions(text string) (int, int, int) {
	font := f.ttf
	size := f.Size
	var (
		totalWidth  = int32(0)
		totalHeight = int32(size)
		maxYBearing = int32(0)
	)
	fupe := font.FUnitsPerEm()
	for _, char := range text {
		idx := font.Index(char)
		hm := font.HMetric(fupe, idx)
		vm := font.VMetric(fupe, idx)
		g := truetype.NewGlyphBuf()
		err := g.Load(font, fupe, idx, truetype.NoHinting)
		if err != nil {
			log.Println(err)
			return 0, 0, 0
		}
		totalWidth += hm.AdvanceWidth
		yB := (vm.TopSideBearing * int32(size)) / fupe
		if yB > maxYBearing {
			maxYBearing = yB
		}
	}

	// Scale to actual pixel size
	totalWidth *= int32(size)
	totalWidth /= fupe

	return int(totalWidth), int(totalHeight), int(maxYBearing)
}

func (f *Font) Render(text string) *Texture {
	width, height, yBearing := f.TextDimensions(text)
	font := f.ttf
	size := f.Size

	// Colors
	fg := image.NewUniform(color.NRGBA{f.FG.R, f.FG.G, f.FG.B, f.FG.A})
	bg := image.NewUniform(color.NRGBA{f.BG.R, f.BG.G, f.BG.B, f.BG.A})

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
