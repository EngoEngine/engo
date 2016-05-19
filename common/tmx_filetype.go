package common

import (
	"bytes"
	"fmt"
        "engo.io/engo"
	"io"
)

// TMXResource contains a level created from a Tile Map XML
type TMXResource struct {
	Level *Level
	url   string
}

func (r TMXResource) URL() string {
	return r.url
}

// tmxLoader is responsible for managing '.tmx' files within 'engo.Files'
type tmxLoader struct {
	levels map[string]TMXResource
}

// Load will load the tmx file and any other image resources that are needed
func (t *tmxLoader) Load(url string, data io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(data)
	if err != nil {
		return err
	}

	lvl, err := createLevelFromTmx(buf.String(), url)
	if err != nil {
		return err
	}

	t.levels[url] = TMXResource{Level: lvl, url: url}
	return nil
}

// Unload removes the preloaded level from the cache
func (t *tmxLoader) Unload(url string) error {
	delete(t.levels, url)
	return nil
}

// Resource retrieves and returns the preloaded level of type 'TMXResource'
func (t *tmxLoader) Resource(url string) (engo.Resource, error) {
	tmx, ok := t.levels[url]
	if !ok {
            return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return tmx, nil
}

func init() {
	engo.Files.Register(".tmx", &tmxLoader{levels: make(map[string]TMXResource)})
}
