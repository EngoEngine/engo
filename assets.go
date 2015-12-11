// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"

	"github.com/golang/freetype/truetype"
	"github.com/paked/webgl"
)

type Resource struct {
	kind string
	name string
	url  string
}

type Loader struct {
	resources []Resource
	images    map[string]*Texture
	jsons     map[string]string
	levels    map[string]*Level
	sounds    map[string]string
	fonts     map[string]*truetype.Font
}

func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		images:    make(map[string]*Texture),
		jsons:     make(map[string]string),
		levels:    make(map[string]*Level),
		sounds:    make(map[string]string),
		fonts:     make(map[string]*truetype.Font),
	}
}

func NewResource(url string) Resource {
	kind := path.Ext(url)
	//name := strings.TrimSuffix(path.Base(url), kind)
	name := path.Base(url)
	return Resource{name: name, url: url, kind: kind[1:]}
}

func (l *Loader) AddFromDir(url string, recurse bool) {
	files, err := ioutil.ReadDir(url)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		furl := url + "/" + f.Name()
		if !f.IsDir() {
			Files.Add(furl)
		} else if recurse {
			Files.AddFromDir(furl, recurse)
		}
	}
}

func (l *Loader) Add(urls ...string) {
	for _, u := range urls {
		r := NewResource(u)
		l.resources = append(l.resources, r)
		log.Println(r)
	}
}

func (l *Loader) Image(name string) *Texture {
	return l.images[name]
}

func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

func (l *Loader) Level(name string) *Level {
	return l.levels[name]
}

func (l *Loader) Sound(name string) ReadSeekCloser {
	f, err := os.Open(l.sounds[name])
	if err != nil {
		return nil
	}
	return f
}

func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		switch r.kind {
		case "png":
			data, err := loadImage(r)
			if err == nil {
				l.images[r.name] = NewTexture(data)
			}
		case "jpg":
			data, err := loadImage(r)
			if err == nil {
				l.images[r.name] = NewTexture(data)
			}
		case "json":
			data, err := loadJSON(r)
			if err == nil {
				l.jsons[r.name] = data
			}
		case "tmx":
			data, err := createLevelFromTmx(r)
			if err == nil {
				l.levels[r.name] = data
			}
		case "wav":
			l.sounds[r.name] = r.url
		case "ttf":
			f, err := loadFont(r)
			if err == nil {
				l.fonts[r.name] = f
			}
		}
	}
	onFinish()
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

func LoadShader(vertSrc, fragSrc string) *webgl.Program {
	vertShader := Gl.CreateShader(Gl.VERTEX_SHADER)
	Gl.ShaderSource(vertShader, vertSrc)
	Gl.CompileShader(vertShader)
	defer Gl.DeleteShader(vertShader)

	fragShader := Gl.CreateShader(Gl.FRAGMENT_SHADER)
	Gl.ShaderSource(fragShader, fragSrc)
	Gl.CompileShader(fragShader)
	defer Gl.DeleteShader(fragShader)

	program := Gl.CreateProgram()
	Gl.AttachShader(program, vertShader)
	Gl.AttachShader(program, fragShader)
	Gl.LinkProgram(program)

	return program
}

type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
	width, height float32
}

func (r *Region) Render(b *Batch, render *RenderComponent, space *SpaceComponent) {
	b.Draw(r, space.Position.X, space.Position.Y, 0, 0, render.Scale.X, render.Scale.Y, 0, render.Color, render.Transparency)
}

func NewRegion(texture *Texture, x, y, w, h int) *Region {
	invTexWidth := 1.0 / float32(texture.Width())
	invTexHeight := 1.0 / float32(texture.Height())

	u := float32(x) * invTexWidth
	v := float32(y) * invTexHeight
	u2 := float32(x+w) * invTexWidth
	v2 := float32(y+h) * invTexHeight
	width := float32(math.Abs(float64(w)))
	height := float32(math.Abs(float64(h)))

	return &Region{texture, u, v, u2, v2, width, height}
}

func (r *Region) Width() float32 {
	return float32(r.width)
}

func (r *Region) Height() float32 {
	return float32(r.height)
}

func (r *Region) Texture() *webgl.Texture {
	return r.texture.id
}

func (r *Region) View() (float32, float32, float32, float32) {
	return r.u, r.v, r.u2, r.v2
}

type Texture struct {
	id     *webgl.Texture
	width  int
	height int
}

func NewTexture(img Image) *Texture {
	var id *webgl.Texture
	if !headless {
		id = Gl.CreateTexture()

		Gl.BindTexture(Gl.TEXTURE_2D, id)

		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_S, Gl.CLAMP_TO_EDGE)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_WRAP_T, Gl.CLAMP_TO_EDGE)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MIN_FILTER, Gl.LINEAR)
		Gl.TexParameteri(Gl.TEXTURE_2D, Gl.TEXTURE_MAG_FILTER, Gl.NEAREST)

		if img.Data() == nil {
			panic("Texture image data is nil.")
		}

		Gl.TexImage2D(Gl.TEXTURE_2D, 0, Gl.RGBA, Gl.RGBA, Gl.UNSIGNED_BYTE, img.Data())
	}

	return &Texture{id, img.Width(), img.Height()}
}

func (t *Texture) Render(b *Batch, render *RenderComponent, space *SpaceComponent) {
	b.Draw(t,
		space.Position.X, space.Position.Y,
		0, 0,
		render.Scale.X, render.Scale.Y,
		0,
		render.Color, render.Transparency)

}

// Width returns the width of the texture.
func (t *Texture) Width() float32 {
	return float32(t.width)
}

// Height returns the height of the texture.
func (t *Texture) Height() float32 {
	return float32(t.height)
}

func (t *Texture) Texture() *webgl.Texture {
	return t.id
}

func (r *Texture) View() (float32, float32, float32, float32) {
	return 0.0, 0.0, 1.0, 1.0
}

type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float32
	Color    color.Color
	Alpha    float32
	Region   *Region
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Anchor:   &Point{0, 0},
		Rotation: 0,
		Color:    color.White,
		Alpha:    1,
		Region:   region,
	}
}

var batchVert = ` 
attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;
uniform vec3 center;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
  gl_Position = vec4(in_Position.x /  uf_Projection.x - center.x,
				 	 in_Position.y / -uf_Projection.y + center.y,
										 0, center.z);
}`

var batchFrag = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`

var hudVert = `
attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;
uniform vec3 center;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
  gl_Position = vec4(in_Position.x / uf_Projection.x - 1.0,
  					 in_Position.y / -uf_Projection.y + 1.0, 0, 1.0);
}`

var hudFrag = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`
