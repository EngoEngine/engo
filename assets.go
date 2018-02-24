package engo

import (
	"fmt"
	"io"
	"path/filepath"
)

// FileLoader implements support for loading and releasing file resources.
type FileLoader interface {
	// Load loads the given resource into memory.
	Load(url string, data io.Reader) error

	// Unload releases the given resource from memory.
	Unload(url string) error

	// Resource returns the given resource, and an error if it didn't succeed.
	Resource(url string) (Resource, error)
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

	// root is the directory which is prepended to every resource url internally.
	root string
}

// SetRoot can be used to change the default directory from `assets` to whatever you want.
//
// Whenever `root` does not start with the directory `assets`, you will not be able to support mobile (Android/iOS)
// since they require you to put all resources within the `assets` directory. More information about that is available
// here: https://godoc.org/golang.org/x/mobile/asset
//
// You can, however, use subfolders within the `assets` folder, and set those as `root`.
func (formats *Formats) SetRoot(root string) {
	formats.root = root
}

// Register registers a resource loader for the given file format.
func (formats *Formats) Register(ext string, loader FileLoader) {
	formats.formats[ext] = loader
}

// load loads the given resource into memory.
func (formats *Formats) load(url string) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		f, err := openFile(filepath.Join(formats.root, url))
		if err != nil {
			return fmt.Errorf("unable to open resource: %s", err)
		}
		defer f.Close()

		return loader.Load(url, f)
	}
	return fmt.Errorf("no `FileLoader` associated with this extension: %q in url %q", ext, url)
}

// Load loads the given resource(s) into memory, stopping at the first error.
func (formats *Formats) Load(urls ...string) error {
	for _, url := range urls {
		err := formats.load(url)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadReaderData loads a resource when you already have the reader for it.
func (formats *Formats) LoadReaderData(url string, f io.Reader) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Load(url, f)
	}
	return fmt.Errorf("no `FileLoader` associated with this extension: %q in url %q", ext, url)
}

// Unload releases the given resource from memory.
func (formats *Formats) Unload(url string) error {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Unload(url)
	}
	return fmt.Errorf("no `FileLoader` associated with this extension: %q in url %q", ext, url)
}

// Resource returns the given resource, and an error if it didn't succeed.
func (formats *Formats) Resource(url string) (Resource, error) {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Resource(url)
	}
	return nil, fmt.Errorf("no `FileLoader` associated with this extension: %q in url %q", ext, url)
}
