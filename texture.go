package eng

import (
	gl "github.com/chsc/gogl/gl33"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"os"
)

const (
	FilterNearest      = gl.NEAREST
	FilterLinear       = gl.LINEAR
	WrapClampToEdge    = gl.CLAMP_TO_EDGE
	WrapRepeat         = gl.REPEAT
	WrapMirroredRepeat = gl.MIRRORED_REPEAT
)

type Texture struct {
	id        gl.Uint
	width     int
	height    int
	minFilter gl.Int
	maxFilter gl.Int
	uWrap     gl.Int
	vWrap     gl.Int
}

func NewTexture(data interface{}) *Texture {
	var m image.Image

	switch data := data.(type) {
	default:
		log.Fatal("NewTexture needs a string or io.Reader")
	case string:
		file, err := os.Open(data)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		m = img
	case io.Reader:
		img, _, err := image.Decode(data)
		if err != nil {
			panic(err)
		}
		m = img
	case image.Image:
		m = data
	}

	b := m.Bounds()
	newm := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	width := m.Bounds().Max.X
	height := m.Bounds().Max.Y

	var id gl.Uint
	gl.GenTextures(1, &id)

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.Sizei(width), gl.Sizei(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&newm.Pix[0]))

	gl.Disable(gl.TEXTURE_2D)

	return &Texture{id, width, height, gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE}
}

func (t *Texture) Split(w, h int) []*Region {
	x := 0
	y := 0
	width := t.Width()
	height := t.Height()

	rows := height / h
	cols := width / w

	startX := x
	tiles := make([]*Region, 0)
	for row := 0; row < rows; row++ {
		x = startX
		for col := 0; col < cols; col++ {
			tiles = append(tiles, NewRegion(t, x, y, w, h))
			x += w
		}
		y += h
	}

	return tiles
}

func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.id)
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Width() int {
	return t.width
}

func (t *Texture) Height() int {
	return t.height
}

func (t *Texture) SetFilter(min, max gl.Int) {
	t.minFilter = min
	t.maxFilter = max
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, max)
}

func (t *Texture) Filter() (gl.Int, gl.Int) {
	return t.minFilter, t.maxFilter
}

func (t *Texture) SetWrap(u, v gl.Int) {
	t.uWrap = u
	t.vWrap = v
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, u)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, v)
}

func (t *Texture) Wrap() (gl.Int, gl.Int) {
	return t.uWrap, t.vWrap
}
