//+build demo

package assets

import "embed"

//go:embed *.png *.tmx
var fs embed.FS

// ReadFile wrapper for embedded assets filesystem.
func ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(name)
}
