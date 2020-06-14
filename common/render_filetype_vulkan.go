//+build vulkan

package common

import "sync/atomic"

type TextureID *VkTextureID

type VkTextureID uint64

var texMap = make(map[*VkTextureID]Image)

var idInc uint64

func createTextureID(img Image) TextureID {
	id := VkTextureID(atomic.AddUint64(&idInc, 1))
	texMap[&id] = img
	return &id
}

// close texture in an api specific way
func (t Texture) close() {
	delete(texMap, t.Texture())
}
