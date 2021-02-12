//!build

package assets

import (
	// go-bindata required for code generation
	_ "github.com/go-bindata/go-bindata"
)

//go:generate go-bindata -pkg=assets ./...
//go:generate gofmt -s -w .
