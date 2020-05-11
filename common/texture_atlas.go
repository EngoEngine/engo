package common

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/gl"
)

// TextureAtlas is a collection of small textures grouped into a big image
type TextureAtlas struct {
	XMLName     xml.Name     `xml:"TextureAtlas"`
	Text        string       `xml:",chardata"`
	ImagePath   string       `xml:"imagePath,attr"`
	SubTextures []SubTexture `xml:"SubTexture"`
}

// SubTexture represents a texture from a region in the TextureAtlas
type SubTexture struct {
	Text   string  `xml:",chardata"`
	Name   string  `xml:"name,attr"`
	X      float32 `xml:"x,attr"`
	Y      float32 `xml:"y,attr"`
	Width  float32 `xml:"width,attr"`
	Height float32 `xml:"height,attr"`
}

type TextureAtlasResource struct {
	texture *gl.Texture     // The original texture
	cache   map[int]Texture // The cell cache cells
	url     string
	Atlas   *TextureAtlas
}

// URL retrieves the url to the .xml file
func (r TextureAtlasResource) URL() string {
	return r.url
}

type textureAtlasLoader struct {
	atlases map[string]*TextureAtlasResource
}

func (t *textureAtlasLoader) Load(url string, data io.Reader) error {
	atlas, err := createAtlasFromXML(data, url)
	if err != nil {
		return err
	}

	t.atlases[url] = atlas
	return nil
}

func (t *textureAtlasLoader) Unload(url string) error {
	delete(t.atlases, url)
	return nil
}

func (t *textureAtlasLoader) Resource(url string) (engo.Resource, error) {
	atlas, ok := t.atlases[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return atlas, nil
}

func createAtlasFromXML(r io.Reader, url string) (*TextureAtlasResource, error) {
	var atlas *TextureAtlas
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, &atlas)
	if err != nil {
		return nil, err
	}

	imgURL := path.Join(path.Dir(url), atlas.ImagePath)
	if err := engo.Files.Load(imgURL); err != nil {
		return nil, fmt.Errorf("failed load texture atlas image: %v", err)
	}

	res, err := engo.Files.Resource(imgURL)
	if err != nil {
		return nil, err
	}

	img, ok := res.(TextureResource)
	if !ok {
		return nil, fmt.Errorf("resource not of type `TextureResource`: %v", url)
	}

	for _, subTexture := range atlas.SubTextures {
		texture := &Texture{
			id:     img.Texture,
			width:  subTexture.Width,
			height: subTexture.Height,
		}

		viewport := engo.AABB{
			Min: engo.Point{
				X: subTexture.X / img.Width,
				Y: subTexture.Y / img.Height,
			},
			Max: engo.Point{
				X: (subTexture.X + subTexture.Width) / img.Width,
				Y: (subTexture.Y + subTexture.Height) / img.Height,
			},
		}

		imgLoader.images[subTexture.Name] = TextureResource{Texture: texture.id, Width: texture.width, Height: texture.height, Viewport: &viewport}
	}

	return &TextureAtlasResource{
		Atlas:   atlas,
		url:     url,
		texture: img.Texture,
	}, nil
}

func init() {
	engo.Files.Register(".xml", &textureAtlasLoader{atlases: make(map[string]*TextureAtlasResource)})
}
