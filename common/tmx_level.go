package common

import (
	"io"
	"path"
	"strconv"
	"strings"

	"engo.io/engo"
	"github.com/Noofbiz/tmx"
)

// createLevelFromTmx unmarshalls and unpacks tmx data into a Level
func createLevelFromTmx(r io.Reader, tmxURL string) (*Level, error) {
	tmx.TMXURL = tmxURL
	tmxLevel, err := tmx.Parse(r)
	if err != nil {
		return nil, err
	}
	level := &Level{}
	level.Orientation = orth
	level.resourceMap = make(map[uint32]Texture)
	level.pointMap = make(map[mapPoint]*Tile)

	// get a map of the gids to textures from the tilesets
	for _, ts := range tmxLevel.Tilesets {
		for _, g := range ts.Grid {
			level.Orientation = g.Orientation
		}
		for _, t := range ts.Tiles {
			for _, i := range t.Image {
				if i.Source != "" {
					tex, err := LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
					if err != nil {
						if strings.HasPrefix(err.Error(), "resource not loaded") {
							err = engo.Files.Load(path.Join(path.Dir(tmxURL), i.Source))
							if err != nil {
								return nil, err
							}
							tex, err = LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
						} else {
							return nil, err
						}
					}
					level.resourceMap[ts.FirstGID+t.ID] = *tex
				}
			}
		}
		for _, i := range ts.Image {
			if i.Source != "" {
				_, err := LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
				if err != nil {
					if strings.HasPrefix(err.Error(), "resource not loaded") {
						err = engo.Files.Load(path.Join(path.Dir(tmxURL), i.Source))
						if err != nil {
							return nil, err
						}
						_, err = LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
					} else {
						return nil, err
					}
				}
				ss := NewSpritesheetWithBorderFromFile(path.Join(path.Dir(tmxURL), i.Source), ts.TileWidth, ts.TileHeight, ts.Spacing, ts.Spacing)
				for i, tex := range ss.Cells() {
					level.resourceMap[ts.FirstGID+uint32(i)] = tex
				}
			}
		}
	}

	level.Orientation = tmxLevel.Orientation
	level.RenderOrder = tmxLevel.RenderOrder
	level.TileWidth = tmxLevel.TileWidth
	level.width = tmxLevel.Width
	level.height = tmxLevel.Height
	level.TileHeight = tmxLevel.TileHeight
	level.NextObjectID = tmxLevel.NextObjectID
	level.Properties = getProperties(tmxLevel.Properties)

	// tile layers
	for _, l := range tmxLevel.Layers {
		tl := &TileLayer{}
		tl.Name = l.Name
		tl.X = float32(l.X)
		tl.OffSetX = float32(l.OffsetX)
		tl.Y = float32(l.Y)
		tl.OffSetY = float32(l.OffsetY)
		tl.Opacity = float32(l.Opacity)
		tl.Visible = l.Visible == 1
		if l.Width != 0 {
			tl.Width = l.Width
		} else {
			tl.Width = tmxLevel.Width
		}
		if l.Height != 0 {
			tl.Height = l.Height
		} else {
			tl.Height = tmxLevel.Height
		}
		tl.Properties = getProperties(l.Properties)
		tl.Tiles = level.unpackTiles(0, 0, tl.Height, tl.Width, l.Data)
		level.TileLayers = append(level.TileLayers, tl)
	}

	//image layers
	for _, l := range tmxLevel.ImageLayers {
		il := &ImageLayer{}
		il.Name = l.Name
		il.Opacity = float32(l.Opacity)
		il.Visible = l.Visible == 1
		il.OffSetX = float32(l.OffsetX)
		il.OffSetY = float32(l.OffsetY)
		il.Properties = getProperties(l.Properties)
		il.Images, err = level.imageTiles(tmxURL, l.Images, il.OffSetX, il.OffSetY)
		if err != nil {
			return nil, err
		}
		level.ImageLayers = append(level.ImageLayers, il)
	}

	// Objects
	for _, o := range tmxLevel.ObjectGroups {
		ol := &ObjectLayer{}
		ol.Color = o.Color
		ol.Name = o.Name
		ol.DrawOrder = o.DrawOrder
		ol.OffSetX = float32(o.OffsetX)
		ol.OffSetY = float32(o.OffsetY)
		ol.Opacity = float32(o.Opacity)
		ol.Visible = o.Visible == 1
		ol.Properties = getProperties(o.Properties)
		for _, tmxobj := range o.Objects {
			object := Object{}
			object.ID = tmxobj.ID
			object.Name = tmxobj.Name
			object.Type = tmxobj.Type
			object.X = float32(tmxobj.X)
			object.Y = float32(tmxobj.Y)
			object.Width = float32(tmxobj.Width)
			object.Height = float32(tmxobj.Height)
			object.Properties = getProperties(tmxobj.Properties)
			object.Tiles = append(object.Tiles, level.tileFromGID(tmxobj.GID, engo.Point{
				X: object.X,
				Y: object.Y,
			}))
			tiles, err := level.imageTiles(tmxURL, tmxobj.Images, object.X, object.Y)
			if err != nil {
				return nil, err
			}
			object.Tiles = append(object.Tiles, tiles...)
			for _, l := range tmxobj.Polygons {
				line := TMXLine{}
				line.Lines = pointStringToLines(l.Points, tmxobj.X, tmxobj.Y)
				line.Type = "Polygon"
				object.Lines = append(object.Lines, line)
			}
			for _, l := range tmxobj.Polylines {
				line := TMXLine{}
				line.Lines = pointStringToLines(l.Points, tmxobj.X, tmxobj.Y)
				line.Type = "Polyline"
				object.Lines = append(object.Lines, line)
			}
			for range tmxobj.Ellipses {
				object.Ellipses = append(object.Ellipses, TMXCircle{
					X:      object.X,
					Y:      object.Y,
					Width:  object.Width,
					Height: object.Height,
				})
			}
			for _, t := range tmxobj.Text {
				text := TMXText{}
				text.Bold = t.Bold == 1
				text.Color = t.Color
				text.FontFamily = t.FontFamily
				text.Halign = t.Halign
				text.Italic = t.Italic == 1
				text.Kerning = t.Kerning == 1
				text.Size = float32(t.PixelSize)
				text.Strikeout = t.Strikeout == 1
				text.Underline = t.Underline == 1
				text.Valign = t.Valign
				text.WordWrap = t.Wrap == 1
				text.CharData = t.CharData
				object.Text = append(object.Text, text)
			}
		}
		level.ObjectLayers = append(level.ObjectLayers, ol)
	}

	return level, nil
}

func pointStringToLines(str string, xOff, yOff float64) []*engo.Line {
	pts := strings.Split(str, " ")
	floatPts := make([][]float64, len(pts))
	for i, x := range pts {
		pt := strings.Split(x, ",")
		floatPts[i] = make([]float64, 2)
		floatPts[i][0], _ = strconv.ParseFloat(pt[0], 64)
		floatPts[i][1], _ = strconv.ParseFloat(pt[1], 64)
	}

	lines := make([]*engo.Line, len(floatPts)-1)

	// Now to globalize line coordinates
	for i := 0; i < len(floatPts)-1; i++ {
		x1 := float32(floatPts[i][0] + xOff)
		y1 := float32(floatPts[i][1] + yOff)
		x2 := float32(floatPts[i+1][0] + xOff)
		y2 := float32(floatPts[i+1][1] + yOff)

		p1 := engo.Point{X: x1, Y: y1}
		p2 := engo.Point{X: x2, Y: y2}
		newLine := &engo.Line{P1: p1, P2: p2}

		lines[i] = newLine
	}

	return lines
}

func (l *Level) unpackTiles(x, y, w, h int, d []tmx.Data) []*Tile {
	var ret []*Tile
	const (
		rd = "right-down"
		ru = "right-up"
		ld = "left-down"
		lu = "left-up"
	)

	switch l.RenderOrder {
	case ru:
		y = h - 1
	case ld:
		x = w - 1
	case lu:
		x = w - 1
		y = h - 1
	}

	for _, data := range d {
		for _, t := range data.Tiles {
			tile := l.tileFromGID(t.GID, l.screenPoint(engo.Point{
				X: float32(x),
				Y: float32(y),
			}))
			ret = append(ret, tile)
			l.pointMap[mapPoint{X: x, Y: y}] = tile
			switch l.RenderOrder {
			case rd:
				x++
				if x >= w {
					x = 0
					y++
				}
			case ru:
				x++
				if x >= w {
					x = 0
					y--
				}
			case ld:
				x--
				if x < 0 {
					x = w
					y++
				}
			case lu:
				x--
				if x < 0 {
					x = w
					y--
				}
			}
		}
		for _, c := range data.Chunks {
			x = c.X
			y = c.Y
			switch l.RenderOrder {
			case ru:
				y += c.Height - 1
			case ld:
				x += c.Width - 1
			case lu:
				x += c.Width - 1
				y += c.Height - 1
			}
			for _, t := range c.Tiles {
				tile := l.tileFromGID(t.GID, l.screenPoint(engo.Point{
					X: float32(x),
					Y: float32(y),
				}))
				ret = append(ret, tile)
				l.pointMap[mapPoint{X: x, Y: y}] = tile
				switch l.RenderOrder {
				case rd:
					x++
					if x >= c.X+c.Width {
						x = c.X
						y++
					}
				case ru:
					x++
					if x >= c.X+c.Width {
						x = c.X
						y--
					}
				case ld:
					x--
					if x < c.X {
						x = c.X + c.Width - 1
						y++
					}
				case lu:
					x--
					if x < c.Width {
						x = c.X + c.Width - 1
						y--
					}
				}
			}
		}
	}
	return ret
}

func (l *Level) imageTiles(tmxURL string, imgs []tmx.Image, x, y float32) ([]*Tile, error) {
	ret := make([]*Tile, 0)
	for _, i := range imgs {
		if i.Source != "" {
			tex, err := LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
			if err != nil {
				if strings.HasPrefix(err.Error(), "resource not loaded") {
					err = engo.Files.Load(path.Join(path.Dir(tmxURL), i.Source))
					if err != nil {
						return nil, err
					}
					tex, err = LoadedSprite(path.Join(path.Dir(tmxURL), i.Source))
				} else {
					return nil, err
				}
			}
			tile := &Tile{
				Image: tex,
				Point: engo.Point{
					X: x,
					Y: y,
				},
			}
			ret = append(ret, tile)
		}
		// TODO: handle image data (not supported by Tiled, but some Java versions do have it)
	}
	return ret, nil
}

func (l *Level) tileFromGID(gid uint32, pt engo.Point) *Tile {
	ret := &Tile{}
	tex := l.resourceMap[gid]
	ret.Image = &tex
	ret.Point = pt
	return ret
}

func getProperties(props []tmx.Property) []Property {
	ret := make([]Property, 0)
	for _, p := range props {
		ret = append(ret, Property{
			Name:  p.Name,
			Type:  p.Type,
			Value: p.Value,
		})
	}
	return ret
}
