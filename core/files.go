package core

import (
	"fmt"
	"path/filepath"
)

// Files manages global resource handling of registered file formats for game
// assets.
var Files *Formats

// FileLoader implements support for loading and releasing file resources.
type FileLoader interface {
	// Load loads the given resource into memory.
	Load(url string) error
	// Unload releases the given resource from memory.
	Unload(url string) error
	// Resource returns the given resource, and a boolean indicating whether the
	// resource was loaded.
	Resource(url string) (Resource, bool)
}

// A Resource represents a game resource, such as an image or a sound.
type Resource interface {
	// URL returns the uniform resource locator of the given resource.
	URL() string
}

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
		return loader.Load(url)
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
