# Embed Demo

## What does it do?

It demonstrates how one can bundle assets into their binary to reduce distribution complexity using the go embed package

## What are important aspects of the code?

These lines are key in this demo:

- `import _ "github.com/go-bindata/go-bindata"` to track go-bindata in go.mod
- `//go:generate go-bindata -pkg=assets ./...` to setup the go generate tool
- `data, err := assets.Asset(file)` to retrieve a bundled asset
- `err = engo.Files.LoadReaderData(file, bytes.NewReader(data))` to load it into engo.Files

## Load each file in a list

```go
func (s *GameScene) Preload() {
	files := []string{
		"tilemap.tmx",
		"sprites.png,
	}
	for _, file := range files {
		data, err := assets.Asset(file)
		if err != nil {
			log.Fatalf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Unable to load asset with URL: %v\n", file)
		}
	}
}
```

## Setup a bundler go file to use go generate

assets/assets.go

```go
//!build

package assets

import (
	// go-bindata required for code generation
	_ "github.com/go-bindata/go-bindata"
)

//go:generate go-bindata -pkg=assets ./...
//go:generate gofmt -s -w .
```
