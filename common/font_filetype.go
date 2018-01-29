package common

import (
	"fmt"
	"io"
	"io/ioutil"

	"engo.io/engo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// FontResource is a wrapper for `*truetype.Font` which is being passed by the the `engo.Files.Resource` method in the
// case of `.ttf` files.
type FontResource struct {
	Font *truetype.Font
	url  string
}

// URL returns the file path for the FontResource.
func (f FontResource) URL() string {
	return f.url
}

// fontLoader is responsible for managing `.ttf` files within `engo.Files`
type fontLoader struct {
	fonts map[string]FontResource
}

// Load processes the data stream and parses it as a freetype font
func (i *fontLoader) Load(url string, data io.Reader) error {
	ttfBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return err
	}

	i.fonts[url] = FontResource{Font: ttf, url: url}
	return nil
}

// Load removes the preloaded font from the cache
func (i *fontLoader) Unload(url string) error {
	delete(i.fonts, url)
	return nil
}

// Resource retrieves the preloaded font, passed as a `FontResource`
func (i *fontLoader) Resource(url string) (engo.Resource, error) {
	texture, ok := i.fonts[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return texture, nil
}

func init() {
	engo.Files.Register(".ttf", &fontLoader{fonts: make(map[string]FontResource)})
}
