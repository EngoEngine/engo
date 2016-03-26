package engi

import (
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/golang/freetype/truetype"
	"github.com/luxengine/math"
	"github.com/paked/webgl"
)

type Resource struct {
	kind string
	name string
	url  string
}

// Asset Loader
type Loader struct {
	resources []Resource
	images    map[string]*Texture
	jsons     map[string]string
	levels    map[string]*Level
	sounds    map[string]string
	fonts     map[string]*truetype.Font
}

// Creates a new Asset Loader
func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		images:    make(map[string]*Texture),
		jsons:     make(map[string]string),
		levels:    make(map[string]*Level),
		sounds:    make(map[string]string),
		fonts:     make(map[string]*truetype.Font),
	}
}

// Creates a new resource from it's url
func NewResource(url string) Resource {
	kind := path.Ext(url)
	//name := strings.TrimSuffix(path.Base(url), kind)
	name := path.Base(url)
	return Resource{name: name, url: url, kind: kind[1:]}
}

// Add all resources from a specified directory
func (l *Loader) AddFromDir(url string, recurse bool) {
	files, err := ioutil.ReadDir(url)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		furl := url + "/" + f.Name()
		if !f.IsDir() {
			Files.Add(furl)
		} else if recurse {
			Files.AddFromDir(furl, recurse)
		}
	}
}

// Add resources from url(s)
func (l *Loader) Add(urls ...string) {
	for _, u := range urls {
		r := NewResource(u)
		l.resources = append(l.resources, r)
		log.Println(r)
	}
}

// Returns a texture from the specified name
func (l *Loader) Image(name string) *Texture {
	return l.images[name]
}

// Returns a JSON string from a specified name
func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

// Returns a level from a specified name
func (l *Loader) Level(name string) *Level {
	return l.levels[name]
}

// Returns audio from a specified name
func (l *Loader) Sound(name string) ReadSeekCloser {
	f, err := os.Open(l.sounds[name])
	if err != nil {
		return nil
	}
	return f
}

// Loads all the loader resources and calls the
// specified callback when done
func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		switch r.kind {
		case "png":
			data, err := loadImage(r)
			if err == nil {
				l.images[r.name] = NewTexture(data)
			}
		case "jpg":
			data, err := loadImage(r)
			if err == nil {
				l.images[r.name] = NewTexture(data)
			}
		case "json":
			data, err := loadJSON(r)
			if err == nil {
				l.jsons[r.name] = data
			}
		case "tmx":
			data, err := createLevelFromTmx(r)
			if err == nil {
				l.levels[r.name] = data
			}
		case "wav":
			l.sounds[r.name] = r.url
		case "ttf":
			f, err := loadFont(r)
			if err == nil {
				l.fonts[r.name] = f
			} else {
				log.Printf("Could not load font %s\n", r.name)
			}
		}
	}
	onFinish()
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

// LoadShader loads a shader and returns a *webgl.Program
func LoadShader(vertSrc, fragSrc string) *webgl.Program {
	vertShader := Gl.CreateShader(Gl.VERTEX_SHADER)
	Gl.ShaderSource(vertShader, vertSrc)
	Gl.CompileShader(vertShader)
	defer Gl.DeleteShader(vertShader)

	fragShader := Gl.CreateShader(Gl.FRAGMENT_SHADER)
	Gl.ShaderSource(fragShader, fragSrc)
	Gl.CompileShader(fragShader)
	defer Gl.DeleteShader(fragShader)

	program := Gl.CreateProgram()
	Gl.AttachShader(program, vertShader)
	Gl.AttachShader(program, fragShader)
	Gl.LinkProgram(program)

	return program
}

type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
	width, height float32
}

// NewRegion Creates  new region from
// a texture, coordinates, and dimensions
func NewRegion(texture *Texture, x, y, w, h float32) *Region {
	invTexWidth := 1.0 / texture.Width()
	invTexHeight := 1.0 / texture.Height()

	u := x * invTexWidth
	v := y * invTexHeight
	u2 := (x + w) * invTexWidth
	v2 := (y + h) * invTexHeight

	width := math.Abs(w)
	height := math.Abs(h)

	return &Region{texture, u, v, u2, v2, width, height}
}

// Width Returns the current region width
func (r *Region) Width() float32 {
	return float32(r.width)
}

// Height Returns the current region height
func (r *Region) Height() float32 {
	return float32(r.height)
}

// Texture Returns the current region texture
func (r *Region) Texture() *webgl.Texture {
	return r.texture.id
}

func (r *Region) View() (float32, float32, float32, float32) {
	return r.u, r.v, r.u2, r.v2
}

type Texture struct {
	id     *webgl.Texture
	width  float32
	height float32
}

// NewTexture returns a new texture from an image
func NewTexture(img Image) *Texture {
	var id *webgl.Texture
	if !headless {
		id = Gl.CreateTexture()

		Gl.BindTexture(Gl.TEXTURE_2D, id)

		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_S, Gl.CLAMP_TO_EDGE)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_T, Gl.CLAMP_TO_EDGE)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MIN_FILTER, Gl.LINEAR)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MAG_FILTER, Gl.NEAREST)

		if img.Data() == nil {
			panic("Texture image data is nil.")
		}

		Gl.TexImage2D(Gl.TEXTURE_2D, 0, Gl.RGBA, Gl.RGBA, Gl.UNSIGNED_BYTE, img.Data())
	}

	return &Texture{id, float32(img.Width()), float32(img.Height())}
}

// Width returns the width of the texture.
func (t *Texture) Width() float32 {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() float32 {
	return t.height
}

// Texture returns the webgl texture of the current texture
func (t *Texture) Texture() *webgl.Texture {
	return t.id
}

func (r *Texture) View() (float32, float32, float32, float32) {
	return 0.0, 0.0, 1.0, 1.0
}

type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float32
	Color    color.Color
	Alpha    float32
	Region   *Region
}

// NewSprite creates and returns a new sprite from a region and coordinates
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
