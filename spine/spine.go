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

type Image struct {
	*eng.Region
	Width      int
	Height     int
	attachment *spine.Attachment
}

func NewImage(path string) *Image {
	texture := eng.NewTexture(path)
	width := texture.Width()
	height := texture.Height()
	region := eng.NewRegion(texture, 0, 0, width, height)
	return &Image{region, width, height, nil}
}

func NewImageRA(path string, a *spine.Attachment) *RegionAttachment {
	texture := eng.NewTexture(path)
	width := texture.Width()
	height := texture.Height()
	region := eng.NewRegion(texture, 0, 0, width, height)
	return NewRegionAttachment(region, a)
	//	return &Image{region, width, height, nil}
}

type Skeleton struct {
	*spine.Skeleton
	base   string
	images map[*spine.Attachment]*RegionAttachment
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
			/*
				x := slot.Bone.WorldX + attachment.X*slot.Bone.M00 + attachment.Y*slot.Bone.M01
				y := slot.Bone.WorldY + attachment.X*slot.Bone.M10 + attachment.Y*slot.Bone.M11
				rotation := slot.Bone.WorldRotation + attachment.Rotation
				xScale := slot.Bone.WorldScaleX + attachment.ScaleX - 1
				yScale := slot.Bone.WorldScaleY + attachment.ScaleY - 1
				if s.FlipX {
					xScale = -xScale
					rotation = -rotation
				}
				if s.FlipY {
					yScale = -yScale
					rotation = -rotation
				}

				color.R = slot.R
				color.G = slot.G
				color.B = slot.B
				color.A = slot.A
			*/
			image.Update(slot)
			//			batch.Draw(image.region, s.X+x, s.Y-y, attachment.OriginX,
			//			attachment.OriginY,
			//			xScale*attachment.WidthRatio,
			//			yScale*attachment.HeightRatio, rotation,
			//			color)
			batch.DrawVerts(image.region, image.Vertices(), nil)
		}
	}
	//}
}

func NewSkeleton(base, file string) *Skeleton {
	path := base + "/" + file

	if data, ok := skeletons[path]; ok {
		log.Println(path)
		return &Skeleton{spine.NewSkeleton(data), base, nil}
	}

	skeletonData := spine.Load(path)
	skeletons[path] = skeletonData
	return &Skeleton{spine.NewSkeleton(skeletonData), base, nil}
}
