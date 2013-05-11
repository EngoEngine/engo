package eng

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type offset struct {
	xoffset  float32
	yoffset  float32
	xadvance float32
}

type Font struct {
	texture *Texture
	regions []*Region
	offsets []*offset
	mapping map[rune]int
}

func NewFont(fnt string, img string) *Font {
	file, err := os.Open(fnt)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	defer file.Close()

	texture := NewTexture(img)
	texture.SetFilter(FilterLinear, FilterLinear)

	font := new(Font)
	font.regions = make([]*Region, 0)
	font.offsets = make([]*offset, 0)
	font.mapping = make(map[rune]int)
	font.texture = texture

	index := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		groups := strings.Split(line, " ")
		switch groups[0] {
		default:
		case "char":
			char := make(map[string]int64)
			for i := 1; i <= 10; i++ {
				k, v := split(groups[i])
				num, _ := strconv.ParseInt(v, 10, 64)
				char[k] = num
			}

			font.mapping[rune(char["id"])] = index
			os := &offset{float32(char["xoffset"]), float32(char["yoffset"]), float32(char["xadvance"])}
			font.offsets = append(font.offsets, os)
			r := NewRegion(texture, int(char["x"]), int(char["y"]), int(char["width"]), int(char["height"]))
			font.regions = append(font.regions, r)
			index += 1
		}
	}

	return font
}

func split(s string) (string, string) {
	strs := strings.Split(s, "=")
	return strs[0], strs[1]
}

func (f *Font) mapRune(ch rune) (int, bool) {
	if f.mapping == nil {
		return int(ch), true
	}
	position, ok := f.mapping[ch]
	return position, ok
}

func (f *Font) Print(batch *Batch, text string, x, y float32, color *Color) {
	xx := x
	for _, v := range text {
		i, ok := f.mapRune(v)
		if ok {
			region := f.regions[i]
			offset := f.offsets[i]
			batch.Draw(region, xx+offset.xoffset, y+offset.yoffset, 0, 0, 1, 1, 0, color)
			xx += offset.xadvance
		}
	}
}

func (f *Font) Texture() *Texture {
	return f.texture
}
