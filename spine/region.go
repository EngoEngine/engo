package spine

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/spine"
	"math"
)

type RegionAttachment struct {
	region *eng.Region
	verts  [8]float32
	offset [8]float32
	a      *spine.Attachment
}

func NewRegionAttachment(region *eng.Region, a *spine.Attachment) *RegionAttachment {
	ra := new(RegionAttachment)
	ra.region = region
	ra.a = a
	ra.updateOffset()
	return ra
}

func (r *RegionAttachment) Region() *eng.Region {
	return r.region
}

func (r *RegionAttachment) Vertices() [8]float32 {
	return r.verts
}

func (r *RegionAttachment) updateOffset() {
	width := r.a.Width
	height := r.a.Height
	localX2 := width / 2
	localY2 := height / 2
	localX := -localX2
	localY := -localY2
	scaleX := r.a.ScaleX
	scaleY := -r.a.ScaleY
	localX *= scaleX
	localY *= scaleY
	localX2 *= scaleX
	localY2 *= scaleY
	rotation := r.a.Rotation
	rads := float64(rotation) * math.Pi / 180
	cos := float32(math.Cos(rads))
	sin := float32(math.Sin(rads))
	x := r.a.X
	y := r.a.Y
	localXCos := localX*cos + x
	localXSin := localX * sin
	localYCos := localY*cos + y
	localYSin := localY * sin
	localX2Cos := localX2*cos + x
	localX2Sin := localX2 * sin
	localY2Cos := localY2*cos + y
	localY2Sin := localY2 * sin
	r.offset[0] = localXCos - localYSin
	r.offset[1] = localYCos + localXSin
	r.offset[2] = localXCos - localY2Sin
	r.offset[3] = localY2Cos + localXSin
	r.offset[4] = localX2Cos - localY2Sin
	r.offset[5] = localY2Cos + localX2Sin
	r.offset[6] = localX2Cos - localYSin
	r.offset[7] = localYCos + localX2Sin
}

func (r *RegionAttachment) Update(slot *spine.Slot) {
	bone := slot.Bone
	s := slot.Skeleton()
	x := s.X + bone.WorldX
	y := s.Y + bone.WorldY
	m00 := bone.M00
	m01 := bone.M01
	m10 := bone.M10
	m11 := bone.M11
	r.verts[0] = r.offset[0]*m00 + r.offset[1]*m01 + x
	r.verts[1] = r.offset[0]*m10 + r.offset[1]*m11 + y
	r.verts[2] = r.offset[2]*m00 + r.offset[3]*m01 + x
	r.verts[3] = r.offset[2]*m10 + r.offset[3]*m11 + y
	r.verts[4] = r.offset[4]*m00 + r.offset[5]*m01 + x
	r.verts[5] = r.offset[4]*m10 + r.offset[5]*m11 + y
	r.verts[6] = r.offset[6]*m00 + r.offset[7]*m01 + x
	r.verts[7] = r.offset[6]*m10 + r.offset[7]*m11 + y
}
