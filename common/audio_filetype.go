package common

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common/internal/decode/mp3"
	"github.com/EngoEngine/engo/common/internal/decode/vorbis"
	"github.com/EngoEngine/engo/common/internal/decode/wav"
)

// audioLoader is responsible for managing audio files within `engo.Files`
type audioLoader struct {
	audios map[string]*Player
}

// Load processes the data stream and parses it as an audio file
func (a *audioLoader) Load(url string, data io.Reader) error {
	var err error
	audioBytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	audioBuffer := bytes.NewReader(audioBytes)

	var player *Player
	switch getExt(url) {
	case ".wav":
		d, err := wav.Decode(&readSeekCloserBuffer{audioBuffer}, SampleRate)
		if err != nil {
			return err
		}

		player, err = newPlayer(d, url)
		if err != nil {
			return err
		}
	case ".mp3":
		d, err := mp3.Decode(&readSeekCloserBuffer{audioBuffer}, SampleRate)
		if err != nil {
			return err
		}

		player, err = newPlayer(d, url)
		if err != nil {
			return err
		}
	case ".ogg":
		d, err := vorbis.Decode(&readSeekCloserBuffer{audioBuffer}, SampleRate)
		if err != nil {
			return err
		}

		player, err = newPlayer(d, url)
		if err != nil {
			return err
		}
	}

	a.audios[url] = player
	return nil
}

// Load removes the preloaded audio file from the cache
func (a *audioLoader) Unload(url string) error {
	delete(a.audios, url)
	return nil
}

// Resource retrieves the preloaded audio file, passed as a `AudioResource`
func (a *audioLoader) Resource(url string) (engo.Resource, error) {
	texture, ok := a.audios[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return texture, nil
}

// LoadedPlayer retrieves the *audio.Player created from the URL
func LoadedPlayer(url string) (*Player, error) {
	res, err := engo.Files.Resource(url)
	if err != nil {
		return nil, err
	}

	audioRes, ok := res.(*Player)
	if !ok {
		return nil, fmt.Errorf("resource not of type `*Player`: %s", url)
	}

	return audioRes, nil
}

// readSeekCloserBuffer is a wrapper to create a ReadSeekCloser
type readSeekCloserBuffer struct {
	inner *bytes.Reader
}

func (r *readSeekCloserBuffer) Close() error {
	r.inner = nil
	return nil
}

func (r *readSeekCloserBuffer) Read(p []byte) (n int, err error) {
	return r.inner.Read(p)
}

func (r *readSeekCloserBuffer) Seek(offset int64, whence int) (int64, error) {
	return r.inner.Seek(offset, whence)
}

func init() {
	engo.Files.Register(".wav", &audioLoader{audios: make(map[string]*Player)})
	engo.Files.Register(".mp3", &audioLoader{audios: make(map[string]*Player)})
	engo.Files.Register(".ogg", &audioLoader{audios: make(map[string]*Player)})
}

// getExt returns the extension of the file(including extensions with `.` in them) from the given url.
func getExt(path string) string {
	ext := ""
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			ext = path[i:]
		}
	}
	return ext
}
