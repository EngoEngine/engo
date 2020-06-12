//+build !vulkan

package common

import "github.com/EngoEngine/engo"

func (t Texture) close() {
	engo.Gl.DeleteTexture(t.id)
}
