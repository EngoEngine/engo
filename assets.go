package engo

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
)

// FileLoader implements support for loading and releasing file resources.
type FileLoader interface {
	// Load loads the given resource into memory.
	Load(url string, data io.Reader) error

	// Unload releases the given resource from memory.
	Unload(url string) error

	// Resource returns the given resource, and a boolean indicating whether the
	// resource was loaded.
	Resource(url string) (Resource, bool)
}

type Resource interface {
	URL() string
}

// Files manages global resource handling of registered file formats for game
// assets.
var Files = &Formats{formats: make(map[string]FileLoader)}

// Formats manages resource handling of registered file formats.
type Formats struct {
	// formats maps from file extensions to resource loaders.
	formats map[string]FileLoader
}

// Register registers a resource loader for the given file format.
func (formats *Formats) Register(ext string, loader FileLoader) {
	formats.formats[ext] = loader
}

// Load loads the given resource into memory.
func (formats *Formats) Load(url string) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Load(url, nil) // TODO: io.Reader instead of nil
	}
	return fmt.Errorf("no resource loader registered for file format %q", ext)
}

// Unload releases the given resource from memory.
func (formats *Formats) Unload(url string) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Unload(url)
	}
	return fmt.Errorf("no resource loader registered for file format %q", ext)
}

// Resource returns the given resource, and a boolean indicating whether the
// resource was loaded.
func (formats *Formats) Resource(url string) (Resource, bool) {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Resource(url)
	}
	log.Printf("no resource loader registered for file format %q", ext)
	return nil, false
}

/* desktop
func loadImage(r Resource) (Image, error) {
	file, err := os.Open(r.URL)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return &ImageObject{newm}, nil
}
*/

/* gopherjs
func loadImage(r Resource) (Image, error) {
	ch := make(chan error, 1)

	img := js.Global.Get("Image").New()
	img.Call("addEventListener", "load", func(*js.Object) {
		go func() { ch <- nil }()
	}, false)
	img.Call("addEventListener", "error", func(o *js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	img.Set("src", r.URL +"?"+strconv.FormatInt(rand.Int63(), 10))

	err := <-ch
	if err != nil {
		return nil, err
	}

	return NewHtmlImageObject(img), nil
}
*/

/* mobile
func loadImage(r Resource) (Image, error) {
	if strings.HasPrefix(r.URL, "assets/") {
		r.URL = r.URL[7:]
	}

	file, err := asset.Open(r.URL)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return &ImageObject{newm}, nil
}
*/

/*

type Loader struct {
	resources []Resource
	images    map[string]*Texture
	jsons     map[string]string
	sounds    map[string]string
	fonts     map[string]*truetype.Font
}

func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		images:    make(map[string]*Texture),
		jsons:     make(map[string]string),
		sounds:    make(map[string]string),
		fonts:     make(map[string]*truetype.Font),
	}
}

func NewResource(url string) Resource {
	kind := path.Ext(url)
	name := path.Base(url)

	if len(kind) == 0 {
		log.Println("WARNING: Cannot load extensionless resource.")
		return Resource{}
	}

	return Resource{name: name, url: url, kind: kind[1:]}
}

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

func (l *Loader) Add(urls ...string) {
	for _, u := range urls {
		r := NewResource(u)
		l.resources = append(l.resources, r)
		log.Println(r)
	}
}

func (l *Loader) Image(name string) *Texture {
	return l.images[name]
}

func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

func (l *Loader) Sound(name string) ReadSeekCloser {
	f, err := os.Open(l.sounds[name])
	if err != nil {
		return nil
	}
	return f
}

func (l *Loader) Font(name string) (*truetype.Font, bool) {
	font, ok := l.fonts[name]
	return font, ok
}


func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		switch r.kind {
		case "png":
			if _, ok := l.images[r.name]; ok {
				continue // with other resources
			}

			data, err := loadImage(r)
			if err != nil {
				log.Println("Error loading resource:", err)
				continue // with other resources
			}

			l.images[r.name] = NewTexture(data)
		case "jpg":
			if _, ok := l.images[r.name]; ok {
				continue // with other resources
			}

			data, err := loadImage(r)
			if err != nil {
				log.Println("Error loading resource:", err)
				continue // with other resources
			}

			l.images[r.name] = NewTexture(data)
		case "json":
			if _, ok := l.jsons[r.name]; ok {
				continue // with other resources
			}

			data, err := loadJSON(r)
			if err != nil {
				log.Println("Error loading resource:", err)
				continue // with other resources
			}

			l.jsons[r.name] = data
		case "wav":
			l.sounds[r.name] = r.url
		case "ttf":
			if _, ok := l.fonts[r.name]; ok {
				continue // with other resources
			}

			f, err := loadFont(r)
			if err != nil {
				log.Println("Error loading resource:", err)
				continue // with other resources
			}

			l.fonts[r.name] = f
		}
	}
	onFinish()
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
	width, height float32
}

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

func (r *Region) Width() float32 {
	return float32(r.width)
}

func (r *Region) Height() float32 {
	return float32(r.height)
}

func (r *Region) Texture() *gl.Texture {
	return r.texture.id
}

func (r *Region) View() (float32, float32, float32, float32) {
	return r.u, r.v, r.u2, r.v2
}

func (r *Region) Close() {
	r.texture.Close()
}

type Texture struct {
	id     *gl.Texture
	width  float32
	height float32
}

func NewTexture(img Image) *Texture {
	var id *gl.Texture
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

func (t *Texture) Texture() *gl.Texture {
	return t.id
}

func (r *Texture) View() (float32, float32, float32, float32) {
	return 0.0, 0.0, 1.0, 1.0
}

func (r *Texture) Close() {
	if !headless {
		Gl.DeleteTexture(r.id)
	}
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

func ImageToNRGBA(img image.Image, width, height int) *image.NRGBA {
	newm := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newm, newm.Bounds(), img, image.Point{0, 0}, draw.Src)

	return newm
}

*/
