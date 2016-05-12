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

	// Resource returns the given resource, and an error if it didn't succeed
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
			return ResourceOpenError{URL: url, Err: err}
		}
		defer readCloser.Close()

		return loader.Load(url, readCloser)
	}
	return ResourceLoaderNotFoundError{Format: ext, URL: url}
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
	return ResourceLoaderNotFoundError{Format: ext, URL: url}
}

// Resource returns the given resource, and an error if it didn't succeed
func (formats *Formats) Resource(url string) (Resource, error) {
	ext := filepath.Ext(url)
	if loader, ok := Files.formats[ext]; ok {
		return loader.Resource(url)
	}
	return nil, ResourceLoaderNotFoundError{Format: ext, URL: url}
}

// A ResourceLoaderNotFoundError is returned whenever the specified file format has no registered `FileLoader` to
// return the requested resource.
type ResourceLoaderNotFoundError struct {
	Format string
	URL    string
}

func (r ResourceLoaderNotFoundError) Error() string {
	return fmt.Sprintf("no `FileLoader` registered for file format %s (within url: %s)", r.Format, r.URL)
}

// A ResourceNotLoadedError is returned whenever the requested `Resource` was unable to be located within the memory.
// This usually indicates the `Load`-step failed.
type ResourceNotLoadedError struct {
	URL string
}

func (r ResourceNotLoadedError) Error() string {
	return fmt.Sprintf("the `FileLoader` was unable to find `Resource` in memory: %s", r.URL)
}

// ResourceOpenError is returned whenever the assets manager was unable to access the file. It was therefore also
// unable to send the file to the `FileLoader`.
//
// Possible reasons for this may be that your program does not have the right permissions to open the file,
// that the path (url) is false, or that the file simply does not exist.
type ResourceOpenError struct {
	URL string

	// Err is the internal error of the `io` method used. Which method this is, is OS-dependant (`io` for desktop,
	// `asset` for mobile, etc.)
	Err error
}

func (r ResourceOpenError) Error() string {
	return fmt.Sprintf("the assets manager was unable to open `Resource`: %s (%s)", r.URL, r.Err.Error())
}
