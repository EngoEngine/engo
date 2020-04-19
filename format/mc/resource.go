package mc

import (
	"fmt"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

func LoadResource(url string) (*MovieClipResource, error) {
	resource, err := engo.Files.Resource(url)
	if err != nil {
		return nil, fmt.Errorf("[MovieClip] [LoadResource] load Resource %q: %s", url, err.Error())
	}

	mcr, ok := resource.(*MovieClipResource)
	if !ok {
		return nil, fmt.Errorf("[MovieClip] [LoadResource] Resource not of type `MovieClipResource` for %q", url)
	}

	return mcr, nil
}

type MovieClipResource struct {
	url           string
	SpriteSheet   *common.Spritesheet
	Actions       []*common.Animation
	DefaultAction *common.Animation
	Drawable      common.Drawable
}

func (r MovieClipResource) URL() string {
	return r.url
}
