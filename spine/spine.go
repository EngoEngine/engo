package spine

import (
	"github.com/ajhager/eng"
	"github.com/ajhager/spine"
	"log"
)

var skeletons map[string]*spine.SkeletonData

func init() {
	skeletons = make(map[string]*spine.SkeletonData)
}

type Animation struct {
	*spine.Animation
}

func NewImageRA(path string, a *spine.Attachment) *RegionAttachment {
	texture := eng.NewTexture(path)
	texture.SetFilter(eng.FilterLinear, eng.FilterLinear)
	width := texture.Width()
	height := texture.Height()
	region := eng.NewRegion(texture, 0, 0, width, height)
	return NewRegionAttachment(region, a)
}

type Skeleton struct {
	*spine.Skeleton
	base   string
	images map[*spine.Attachment]*RegionAttachment
}

func (s *Skeleton) Apply(a *Animation, time float32, loop bool) {
	a.Apply(s.Skeleton, time, loop)
}

func (s *Skeleton) Mix(a *Animation, time float32, loop bool, alpha float32) {
	a.Mix(s.Skeleton, time, loop, alpha)
}

var color *eng.Color

func init() {
	color = eng.NewColor(1, 1, 1, 1)
}

func (s *Skeleton) Draw(batch *eng.Batch) {
	if s.images == nil {
		s.images = make(map[*spine.Attachment]*RegionAttachment)
	}
	images := s.images

	for _, slot := range s.DrawOrder {
		attachment := slot.Attachment
		image, ok := images[attachment]
		//		if attachment == nil && ok {
		//delete(images, attachment)
		//		} else {
		if ok && image.a != attachment {
			ok = false
		}
		if !ok {
			image = NewImageRA(s.base+"/"+attachment.Name+".png", attachment)
			if image != nil {
				//				image.a = attachment
				imageWidth := float32(image.region.Width())
				imageHeight := float32(image.region.Height())
				attachment.WidthRatio = attachment.Width / imageWidth
				attachment.HeightRatio = attachment.Height / imageHeight
				attachment.OriginX = imageWidth / 2.0
				attachment.OriginY = imageHeight / 2.0
			} else {
				// blah
			}
			s.images[attachment] = image
		}
		if image != nil {
			image.Update(slot)
			batch.DrawVerts(image.region, image.Vertices(), nil)
		}
	}
	//}
}

func (s *Skeleton) Animation(name string) *Animation {
	return &Animation{s.FindAnimation(name)}
}

func NewSkeleton(base, file string) *Skeleton {
	path := base + "/" + file

	if data, ok := skeletons[path]; ok {
		log.Println(path)
		return &Skeleton{spine.NewSkeleton(data), base, nil}
	}

	skeletonData := spine.Load(path)
	skeletons[path] = skeletonData
	s := &Skeleton{spine.NewSkeleton(skeletonData), base, nil}
	s.FlipY = true
	return s
}
