package common

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
)

type testAudio struct {
	ecs.BasicEntity
	AudioComponent
}

type testAudioScene struct {
	w             *ecs.World
	audioSystem   *AudioSystem
	ogg, mp3, wav testAudio
}

func (*testAudioScene) Type() string { return "TestAudioScene" }

func (*testAudioScene) Preload() {
	if err := engo.Files.Load("1.ogg", "sfx_coin_double2.wav", "TripleShot.mp3"); err != nil {
		panic(err)
	}
}

func (t *testAudioScene) Setup(u engo.Updater) {
	var err error
	t.w = u.(*ecs.World)
	t.audioSystem = &AudioSystem{}
	t.w.AddSystem(t.audioSystem)

	t.ogg = testAudio{BasicEntity: ecs.NewBasic()}
	if t.ogg.AudioComponent.Player, err = LoadedPlayer("1.ogg"); err != nil {
		panic(err)
	}
	t.audioSystem.Add(&t.ogg.BasicEntity, &t.ogg.AudioComponent)

	t.wav = testAudio{BasicEntity: ecs.NewBasic()}
	if t.wav.AudioComponent.Player, err = LoadedPlayer("sfx_coin_double2.wav"); err != nil {
		panic(err)
	}
	t.audioSystem.Add(&t.wav.BasicEntity, &t.wav.AudioComponent)

	t.mp3 = testAudio{BasicEntity: ecs.NewBasic()}
	if t.mp3.AudioComponent.Player, err = LoadedPlayer("TripleShot.mp3"); err != nil {
		panic(err)
	}
	t.audioSystem.Add(&t.mp3.BasicEntity, &t.mp3.AudioComponent)
}

// TestAudioSystemIntegrationNormalUse tests using the AudioSystem as a part of
// engo. It doesn't fail on data because sometimes things "slip" and the data doesn't
// come out exactly right, but that's okay. If the audio system changes, check
// the logs and rerun accordingly.
func TestAudioSystemIntegrationNormalUse(t *testing.T) {
	s := testAudioScene{}
	engo.Run(engo.RunOptions{
		HeadlessMode: true,
		NoRun:        true,
		AssetsRoot:   "testdata",
	}, &s)
	p := s.audioSystem.otoPlayer.(*stepPlayer)

	s.w.Update(1)
	exp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	actual := p.Bytes()
	for i, b := range exp {
		if b != actual[i] {
			t.Logf("Audio byte values were incorrect. (First step)\nWanted: %v\nGot: %v\nIndex: %v", b, actual[i], i)
		}
	}
	s.ogg.Player.Play()
	s.w.Update(1)
	p.Step()
	exp = []byte{254, 255, 253, 255, 252, 255, 252, 255, 253, 255, 253, 255, 254, 255}
	actual = p.Bytes()
	for i, b := range exp {
		if b != actual[i] {
			t.Logf("Audio byte values were incorrect. (Second step)\nWanted: %v\nGot: %v\nIndex: %v", b, actual[i], i)
		}
	}
	s.mp3.Player.Play()
	s.w.Update(1)
	p.Step()
	exp = []byte{247, 255, 250, 255, 247, 255, 250, 255, 246, 255, 250, 255, 246, 255, 249}
	actual = p.Bytes()
	for i, b := range exp {
		if b != actual[i] {
			t.Logf("Audio byte values were incorrect. (Third step)\nWanted: %v\nGot: %v\nIndex: %v", b, actual[i], i)
		}
	}
	s.wav.Player.Play()
	s.w.Update(1)
	p.Step()
	exp = []byte{70, 25, 55, 25, 181, 25, 164, 25, 46, 33, 26, 33, 198, 33, 178, 33, 25, 33, 5, 33}
	actual = p.Bytes()
	for i, b := range exp {
		if b != actual[i] {
			t.Logf("Audio byte values were incorrect. (Fourth step)\nWanted: %v\nGot: %v\nIndex: %v", b, actual[i], i)
		}
	}
	s.ogg.Player.Pause()
	s.w.Update(1)
	p.Step()
	exp = []byte{90, 43, 90, 43, 92, 43, 92, 43, 40, 9, 40, 9, 37, 9, 37, 9, 177, 213, 177, 213, 175}
	actual = p.Bytes()
	for i, b := range exp {
		if b != actual[i] {
			t.Logf("Audio byte values were incorrect. (Fifth step)\nWanted: %v\nGot: %v\nIndex: %v", b, actual[i], i)
		}
	}
	p.Step()
}

func TestAudioLoaderLoadOgg(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Error while loading. Error: %v", err)
	}
	_, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Error while getting LoadedPlayer for ogg. Error: %v", err)
		return
	}
}

func TestAudioLoaderLoadWav(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("sfx_coin_double2.wav"); err != nil {
		t.Errorf("Error while loading. Error: %v", err)
	}
	_, err := LoadedPlayer("sfx_coin_double2.wav")
	if err != nil {
		t.Errorf("Error while getting LoadedPlayer for wav. Error: %v", err)
		return
	}
}

func TestAudioLoaderLoadMP3(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("TripleShot.mp3"); err != nil {
		t.Errorf("Error while loading. Error: %v", err)
	}
	_, err := LoadedPlayer("TripleShot.mp3")
	if err != nil {
		t.Errorf("Error while getting LoadedPlayer. Error: %v", err)
		return
	}
}

type testReader struct {
	readError bool
}

func (t testReader) Read(p []byte) (int, error) {
	if t.readError {
		return 0, errors.New("Read Error")
	}
	return 0, io.EOF
}

func TestAudioLoaderLoadReadallError(t *testing.T) {
	tr := testReader{readError: true}
	if err := engo.Files.LoadReaderData("test.mp3", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being read.")
	}
	if err := engo.Files.LoadReaderData("test.wav", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being read.")
	}
	if err := engo.Files.LoadReaderData("test.ogg", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being read.")
	}
}

func TestAudioLoaderLoadDecodeError(t *testing.T) {
	tr := testReader{}
	if err := engo.Files.LoadReaderData("test.mp3", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being loaded.")
	}
	if err := engo.Files.LoadReaderData("test.wav", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being loaded.")
	}
	if err := engo.Files.LoadReaderData("test.ogg", tr); err == nil {
		t.Error("Malformed io.Reader did not throw the error while being loaded.")
	}
}

func TestAudioLoaderUnload(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("sfx_coin_double2.wav"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	if _, err := LoadedPlayer("sfx_coin_double2.wav"); err != nil {
		t.Errorf("Coud not get player from loaded file. Error was: %v\n", err)
	}
	if err := engo.Files.Unload("sfx_coin_double2.wav"); err != nil {
		t.Errorf("Could not unload file. Error was: %v\n", err)
	}
	if _, err := LoadedPlayer("sfx_coin_double2.wav"); err == nil {
		t.Error("Loaded a previously unloaded player.")
	} else {
		if !strings.HasPrefix(err.Error(), "resource not loaded by `FileLoader`:") {
			t.Errorf("Unexpected error while loading previously unloaded player. Error was: %v\n", err)
		}
	}
}

func TestAudioPlayerURL(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	p, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Could not get player. Error was: %v\n", err)
	}
	if p.URL() != "1.ogg" {
		t.Errorf("Wrong URL reported from player. Wanted: %v\nGot: %v\n", "1.ogg", p.URL())
	}
}

func TestAudioPlayerClose(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	p, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Could not get player. Error was: %v\n", err)
	}
	if err = p.Close(); err != nil {
		t.Errorf("Could not close player. Error was: %v\n", err)
	}
	if err = p.Close(); err == nil {
		t.Errorf("Did not get an error while closing already closed player.")
	} else {
		if !strings.HasPrefix(err.Error(), "audio: the player is already closed") {
			t.Errorf("Wrong error when closing already closed player. Error was: %v\n", err)
		}
	}
}

func TestAudioPlayerPlay(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	p, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Could not get player. Error was: %v\n", err)
	}
	if p.IsPlaying() {
		t.Error("Player was playing before play was called.")
	}
	p.Play()
	if !p.IsPlaying() {
		t.Error("Player was not playing after play was called.")
	}
	p.Pause()
	if p.IsPlaying() {
		t.Error("Player was playing after pause was called.")
	}
	p.Seek(time.Second / 5)
	if p.IsPlaying() {
		t.Error("Player was playing after pause when seek was called, but play was not.")
	}
	p.Play()
	if !p.IsPlaying() {
		t.Error("Player was not playing after play was called, after pause and seek.")
	}
}

func TestAudioPlayerSeek(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	p, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Could not get player. Error was: %v\n", err)
	}
	if p.Current() != time.Second*0 {
		t.Error("Newly created player's duration was not zero.")
	}
	p.Seek(time.Second / 5)
	if p.Current() != time.Second/5 {
		t.Error("Seek didn't seek to one fifth of a second.")
	}
	p.Seek(time.Second / 2)
	if p.Current() != time.Second/2 {
		t.Error("Didn't seek from one fifth of a second to one half of a second")
	}
	p.Rewind()
	if p.Current() != time.Second*0 {
		t.Error("Rewind didn't set to zero time.")
	}
}

func TestAudioPlayerVolume(t *testing.T) {
	engo.Files.SetRoot("testdata")
	if err := engo.Files.Load("1.ogg"); err != nil {
		t.Errorf("Could not load file. Error was: %v\n", err)
	}
	p, err := LoadedPlayer("1.ogg")
	if err != nil {
		t.Errorf("Could not get player. Error was: %v\n", err)
	}
	buf := bytes.NewBuffer([]byte{})
	log.SetOutput(buf)
	if p.GetVolume() != 1 {
		t.Error("Initial volume was not 1")
	}
	p.SetVolume(0.5)
	if p.GetVolume() != 0.5 {
		t.Error("Volume was not 0.5 after being set to it")
	}
	p.SetVolume(-1)
	if p.GetVolume() != 0.5 {
		t.Error("Volume was not retained after trying to set to an invalid value")
	}
	if !strings.HasSuffix(buf.String(), "Volume can only be set between zero and one. Volume was not set.\n") {
		t.Errorf("Logged value was not what was expected. Got: %v\n", buf.String())
	}
	buf.Reset()
	p.SetVolume(1)
	if p.GetVolume() != 1 {
		t.Error("Volume was not 1 after being set to it")
	}
	p.SetVolume(10)
	if p.GetVolume() != 1 {
		t.Error("Volume was not retained after trying to set it to an invalid value")
	}
	if !strings.HasSuffix(buf.String(), "Volume can only be set between zero and one. Volume was not set.\n") {
		t.Errorf("Logged value was not what was expected. Got: %v\n", buf.String())
	}
}

func TestAudioMasterVolume(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	log.SetOutput(buf)
	if GetMasterVolume() != 1 {
		t.Error("Initial volume was not 1")
	}
	SetMasterVolume(0.5)
	if GetMasterVolume() != 0.5 {
		t.Error("Volume was not 0.5 after being set to it")
	}
	SetMasterVolume(-1)
	if GetMasterVolume() != 0.5 {
		t.Error("Volume was not retained after trying to set to an invalid value")
	}
	if !strings.HasSuffix(buf.String(), "Master Volume can only be set between zero and one. Volume was not set.\n") {
		t.Errorf("Logged value was not what was expected. Got: %v\n", buf.String())
	}
	buf.Reset()
	SetMasterVolume(1)
	if GetMasterVolume() != 1 {
		t.Error("Volume was not 1 after being set to it")
	}
	SetMasterVolume(10)
	if GetMasterVolume() != 1 {
		t.Error("Volume was not retained after trying to set it to an invalid value")
	}
	if !strings.HasSuffix(buf.String(), "Master Volume can only be set between zero and one. Volume was not set.\n") {
		t.Errorf("Logged value was not what was expected. Got: %v\n", buf.String())
	}
}
