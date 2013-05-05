package eng

import (
	gl "github.com/chsc/gogl/gl21"
	"math"
)

type Region struct {
	texture       *Texture
	u, v          gl.Float
	u2, v2        gl.Float
	width, height int
}

func NewRegion(texture *Texture, x, y, w, h int) *Region {
	invTexWidth := 1.0 / float32(texture.Width())
	invTexHeight := 1.0 / float32(texture.Height())

	u := float32(x) * invTexWidth
	v := float32(y+h) * invTexHeight
	v2 := float32(y) * invTexHeight
	u2 := float32(x+w) * invTexWidth
	width := int(math.Abs(float64(w)))
	height := int(math.Abs(float64(h)))

	return &Region{texture, gl.Float(u), gl.Float(v), gl.Float(u2), gl.Float(v2), width, height}
}
