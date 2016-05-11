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

// Resource represents a game resource, such as an image or a sound.
type Resource interface {
	// URL returns the uniform resource locator of the given resource.
	URL() string
}

// Files manages global resource handling of registered file formats for game
// assets.
var Files = &Formats{formats: make(map[string]FileLoader)}

// Formats manages resource handling of registered file formats.
type Formats struct {
	// formats maps from file extensions to resource loaders.
	formats map[string]FileLoader

	// root is the directory which is prepended to every resource url internally
	root string
}

// SetRoot can be used to change the default directory from `assets` to whatever you want.
//
// Whenever `root` does not start with the directory `assets`, you will not be able to support mobile (Android/iOS)
// since they require you to put all resources within the `assets` directory. You can, however, use subfolders within
// the `assets` folder, and set those as `root`.
func (formats *Formats) SetRoot(root string) {
	formats.root = root
}

// Register registers a resource loader for the given file format.
func (formats *Formats) Register(ext string, loader FileLoader) {
	formats.formats[ext] = loader
}

// Load loads the given resource into memory.
func (formats *Formats) Load(url string) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		readCloser, err := openFile(filepath.Join(formats.root, url))
		if err != nil {
			return err
		}
		defer readCloser.Close()

		return loader.Load(url, readCloser)
	}
	return fmt.Errorf("no resource loader registered for file format %q", ext)
}

// LoadMany loads the given resources into memory, stopping at the first error.
func (formats *Formats) LoadMany(urls ...string) error {
	for _, url := range urls {
		err := formats.Load(url)
		if err != nil {
			return err
		}
	}
	return nil
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

/*
type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

// Region has been removed, and is now called core.Texture

func ImageToNRGBA(img image.Image, width, height int) *image.NRGBA {
	newm := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newm, newm.Bounds(), img, image.Point{0, 0}, draw.Src)

	return newm
}

*/
