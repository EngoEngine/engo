package mc

import (
	"fmt"
	"io"

	"github.com/EngoEngine/engo"
)

func init() {
	engo.Files.Register(".mc.json", NewMovieClipLoader())
}

// NewMovieClipLoader create loader for MovieClip files
func NewMovieClipLoader() engo.FileLoader {
	return &movieClipLoader{resource: make(map[string]*MovieClipResource)}
}

type movieClipLoader struct {
	resource map[string]*MovieClipResource
}

// Load loads the given resource into memory.
func (l *movieClipLoader) Load(url string, data io.Reader) error {
	resource, err := parseMC(url, data)
	if err != nil {
		return err
	}

	if resource.DefaultAction == nil && len(resource.Actions) > 0 {
		resource.DefaultAction = resource.Actions[0]
	}

	resource.Drawable = resource.SpriteSheet.Drawable(resource.DefaultAction.Frames[0])

	l.resource[url] = resource

	return nil
}

// Unload releases the given resource from memory.
func (l *movieClipLoader) Unload(url string) error {
	delete(l.resource, url)
	return nil
}

// Resource returns the given resource, and an error if it didn't succeed.
func (l *movieClipLoader) Resource(url string) (engo.Resource, error) {
	res, ok := l.resource[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return res, nil
}
