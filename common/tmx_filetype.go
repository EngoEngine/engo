package common

import (
	"fmt"
	"io"
	"io/ioutil"

	"engo.io/engo"
)

// TMXResource contains a level created from a Tile Map XML
type TMXResource struct {
	// Level holds the reference to the parsed TMX level
	Level *Level
	url   string
}

// URL retrieves the url to the .tmx file
func (r TMXResource) URL() string {
	return r.url
}

// tmxLoader is responsible for managing '.tmx' files within 'engo.Files'.
// You can generate a TMX file with the Tiled map editor.
type tmxLoader struct {
	levels map[string]TMXResource
}

// Load will load the tmx file and any other image resources that are needed
func (t *tmxLoader) Load(url string, data io.Reader) error {
	tmxBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	lvl, err := createLevelFromTmx(tmxBytes, url)
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
