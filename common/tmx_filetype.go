package common

import (
	"bytes"
	"fmt"
        "engo.io/engo"
	"io"
)

// TMXResource is a wrapper for a level that was created from a Tile Map XML
type TMXResource struct {
	Level *Level
	url   string
}

func (r TMXResource) URL() string {
	return r.url
}

// tmxLoader is responsible for managing '.tmx' files within 'engo.Files'
type tmxLoader struct {
	tmxs map[string]TMXResource
}

// Load will read the tmx file into a string to create a level using createLevelFromTmx.  Any image files required for creating the level will be loaded in as needed.
func (t *tmxLoader) Load(url string, data io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(data)
	if err != nil {
		return err
	}

	lvl, err := createLevelFromTmx(buf.String())
	if err != nil {
		return err
	}

	t.tmxs[url] = TMXResource{Level: lvl, url: url}
	return nil
}

// Unload removes the preloaded level from the cache
func (t *tmxLoader) Unload(url string) error {
	delete(t.tmxs, url)
	return nil
}

// Resource retrieves the preloaded level, passed as a 'TMXResource'
func (t *tmxLoader) Resource(url string) (engo.Resource, error) {
	tmx, ok := t.tmxs[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return tmx, nil
}

func init() {
	engo.Files.Register(".tmx", &tmxLoader{tmxs: make(map[string]TMXResource)})
}
