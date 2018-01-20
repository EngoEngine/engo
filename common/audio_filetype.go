//+build !windows,!android

package common

import (
	"bytes"
	"fmt"

	"engo.io/audio"
	"engo.io/engo"
)

// AudioResource is a wrapper for `*Player` which is being passed by the the `engo.Files.Resource` method in the
// case of `.wav` files.
type AudioResource struct {
	Player *audio.Player
	url    string
}

func (f AudioResource) URL() string {
	return f.url
}

// audioLoader is responsible for managing `.wav` files within `engo.Files`
type audioLoader struct {
	audios map[string]AudioResource
}

// Load removes the preloaded audio file from the cache
func (l *audioLoader) Unload(url string) error {
	delete(l.audios, url)
	return nil
}

// Resource retrieves the preloaded audio file, passed as a `AudioResource`
func (l *audioLoader) Resource(url string) (engo.Resource, error) {
	texture, ok := l.audios[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return texture, nil
}

// readSeekCloserBuffer is a wrapper to create a ReadSeekCloser
type readSeekCloserBuffer struct {
	inner *bytes.Reader
}

func (r *readSeekCloserBuffer) Close() error {
	return nil
}

func (r *readSeekCloserBuffer) Read(p []byte) (n int, err error) {
	return r.inner.Read(p)
}

func (r *readSeekCloserBuffer) Seek(offset int64, whence int) (int64, error) {
	return r.inner.Seek(offset, whence)
}

func init() {
	engo.Files.Register(".wav", &audioLoader{audios: make(map[string]AudioResource)})
}
