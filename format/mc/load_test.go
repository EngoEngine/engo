package mc_test

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/format/mc"
)

var dataMC = `{
  "mc": {
    "run": {
      "frameRate": 24,
      "events": [
        {
          "name": "@action",
          "frame": 1
        }
      ],
      "labels": [
        {
          "name": "action",
          "frame": 1,
          "end": 1
        }
      ],
      "frames": [
        {
          "y": -178,
          "res": "action_1",
          "duration": 1,
          "x": -147
        }
      ]
    }
  },
  "file": "test.png",
  "res": {
    "action_1": {
      "x": 1,
      "y": 1,
      "w": 256,
      "h": 216
    }
  }
}`

func TestLoad(t *testing.T) {
	dir, err := createAssets()
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}
	defer os.RemoveAll(dir)

	// Start an instance of engo
	engo.Run(engo.RunOptions{
		NoRun:        true,
		HeadlessMode: true,
		AssetsRoot:   dir,
	}, &testScene{})

	err = engo.Files.Load("test.mc.json")
	if err != nil {
		t.Errorf("Unable to load MC file for testing. Error was: %v", err)
	}

	mcr, err := mc.LoadResource("test.mc.json")
	if err != nil || mcr == nil {
		t.Errorf("Unable to load MCResource. Error was: %v", err)
	}

	err = engo.Files.Unload("test.mc.json")
	if err != nil {
		t.Errorf("Unable to unload MCResource. Error was: %v", err)
	}

	mcr, err = mc.LoadResource("test.mc.json")
	if err == nil || mcr != nil {
		t.Errorf("Cat not unload MCResource")
	}
}

type testScene struct{}

func (*testScene) Preload() {}

func (*testScene) Setup(engo.Updater) {}

func (*testScene) Type() string { return "testScene" }

func createAssets() (string, error) {
	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory for testing, error: %v", err)
	}

	// Create an image
	buf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	err = png.Encode(buf, img)
	if err != nil {
		return "", fmt.Errorf("unable to encode png from image: %v", err)
	}

	tmpfn := filepath.Join(dir, "test.png")
	if err = ioutil.WriteFile(tmpfn, buf.Bytes(), 0666); err != nil {
		return "", fmt.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	// Create an MC file
	buf = bytes.NewBuffer([]byte(dataMC))
	tmpfn = filepath.Join(dir, "test.mc.json")
	if err = ioutil.WriteFile(tmpfn, buf.Bytes(), 0666); err != nil {
		return "", fmt.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	return dir, nil
}
