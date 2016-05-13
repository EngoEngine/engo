//+build !windows,!netgo,!android

package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"engo.io/engo"
)

// AudioResource is a wrapper for `*Player` which is being passed by the the `engo.Files.Resource` method in the
// case of `.wav` files.
type AudioResource struct {
	Player *Player
	url    string
}

func (f AudioResource) URL() string {
	return f.url
}

// audioLoader is responsible for managing `.wav` files within `engo.Files`
type audioLoader struct {
	audios map[string]AudioResource
}

// Load processes the data stream and parses it as an audio file
func (i *audioLoader) Load(url string, data io.Reader) error {
	audioBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	audioBuffer := bytes.NewReader(audioBytes)
	player, err := NewPlayer(&readSeekCloserBuffer{audioBuffer}, 0, 0)
	if err != nil {
		return fmt.Errorf("%s (are you running `core.AudioSystemPreload()` before preloading .wav files?)", err.Error())
	}

	i.audios[url] = AudioResource{Player: player, url: url}
	return nil
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
