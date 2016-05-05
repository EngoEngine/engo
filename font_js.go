// +build netgo

package engo

//TODO go generate to convert fonts to ones readable by engo
import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/gopherjs/gopherjs/js"
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
	var ok bool
	f.TTF, ok = Files.fonts[f.URL]
	if !ok {
		return fmt.Errorf("could not find preloaded font: %s", f.URL)
	}

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

func (f *Font) RenderNRGBA(text string) *js.Object {
	log.Println("[WARNING] RenderNRGBA not implemented on Gopherjs")
	return &js.Object{}
}

func (f *Font) Render(text string) *Texture {
	log.Println("render called")
	nrgba := f.RenderNRGBA(text)

	// Create texture
	imObj := &ImageObject{nrgba}
	return NewTexture(imObj)
}
