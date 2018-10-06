package common

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"engo.io/engo"
)

var testTMXtmpl = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" tiledversion="1.1.5" orientation="{{ .Orientation }}" renderorder="{{ .RenderOrder }}" width="3" height="3" tilewidth="16" tileheight="16" infinite="0" nextobjectid="1">
 <tileset firstgid="1" name="test" tilewidth="16" tileheight="16" spacing="1" tilecount="468" columns="26">
  {{ if .Grid }}
  <grid orientation="isometric" width="1" height="1"/>
  {{ end }}
	{{ if .Tiles }}
	<tile id="0">
   <properties>
    <property name="walkable" type="bool" value="false"/>
   </properties>
   <image width="132" height="99" source="test.png{{ .BadExtensions }}"/>
  </tile>
  <tile id="1">
   <image width="132" height="99" source="test.png{{ .BadExtensions }}"/>
  </tile>
	{{ else }}
  <image source="test.png{{ .BadExtensions }}" width="457" height="305"/>
	{{ end }}
 </tileset>
 <layer name="Tile Layer 1" width="3" height="3">
 	{{ if .ChunkData }}
	<data encoding="base64" compression="zlib">
	 <chunk x="-32" y="-16" width="16" height="16">
	 eJxjYBgFQwUoMzIwqKBhUoADUL0jGh5K+mOB6uPQ8CigDAAA/zEGbg==
	</chunk>
	 <chunk x="-16" y="-16" width="16" height="16">
	 eJxjYBgFAwlUGDGxKiPx+h0ZMbHTENIfx4iJ40nQPwooAwCK+AfC
	</chunk>
	</data>
	{{ else }}
  <data encoding="base64" compression="zlib">
   eJx7zcDA8AaI3wLxdyBOYWRgkAJiZyB2AWJBIAYAfvME1w==
  </data>
	{{ end }}
 </layer>
 <layer name="Tile Layer 2">
	<data encoding="base64" compression="zlib">
	 eJx7zcDA8AaI3wLxdyBOYWRgkAJiZyB2AWJBIAYAfvME1w==
	</data>
 </layer>
 <objectgroup name="Object Layer 1">
	<object id="1" name="Rectangle" x="10" y="13" width="25" height="23"/>
	<object id="2" x="2" y="2" width="25" height="22">
		{{ if .ObjectImageTest }}
		<image source="objimgtest.png{{ .BadObjectImageExtension }}"/>
		{{ else }}
		<image source="test.png{{ .BadObjectImageExtension }}"/>
		{{ end }}
	</object>
	<object id="12" x="45.875" y="-12.875">
	 <polygon points="0,0 0.5,24.75 30.375,19.875 13.75,-15 -21.375,-11.125 -19.75,14.75"/>
	</object>
	<object id="13" x="3.875" y="51.375" width="47.25" height="20.5">
	 <ellipse/>
	</object>
	<object id="14" x="-12.75" y="45.25">
	 <polyline points="0,0 -17.125,-8.625 3.75,-13.75 -8.375,-26.125 1.25,-28.5 -2.75,-31"/>
	</object>
	<object id="15" x="6.05469" y="71.8516" width="89.5469" height="18.8438">
	 <properties>
		<property name="beep" type="bool" value="false"/>
		<property name="boop" type="bool" value="true"/>
	 </properties>
	 <text wrap="1">Hello World</text>
	</object>
 </objectgroup>
 <group name="Group 1" offsetx="2" offsety="2">
	<objectgroup name="Object Layer 2"/>
	<imagelayer name="Image Layer 2" offsetx="5" offsety="5"/>
 </group>
 <imagelayer name="Image Layer 1">
	<image source="test.png{{ .BadImageExtension }}"/>
 </imagelayer>
</map>
`

var badTMX = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" tiledversion="1.1.5" orientation="orthogonal" renderorder="right-down" width="3" height="3" tilewidth="16" tileheight="16" infinite="0" nextobjectid="1">
`

type tmxData struct {
	Orientation, RenderOrder                                  string
	BadExtensions, BadImageExtension, BadObjectImageExtension string
	Grid, Tiles, InvalidImageTile, ObjectImageTest, ChunkData bool
}

type tmxTestScene struct{}

func (*tmxTestScene) Preload() {}

func (*tmxTestScene) Setup(engo.Updater) {}

func (*tmxTestScene) Type() string { return "testScene" }

func TestTMXFiletypeLoad(t *testing.T) {
	// Start an instance of engo
	engo.Run(engo.RunOptions{
		NoRun:        true,
		HeadlessMode: true,
	}, &tmxTestScene{})

	// Create an image in memory
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 457, 305))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	// Load the image
	err = engo.Files.LoadReaderData("test.png", imgbuf)
	if err != nil {
		t.Errorf("Unable to load test png. Error was: %v", err)
	}

	// Load the template
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	// Load tmx
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}
}

func TestTMXTileNotLoadedFileNotExist(t *testing.T) {
	// Ensure the file is not loaded
	engo.Files.Unload("test.png")

	// Load the template
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-up",
		ChunkData:   true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	// Load tmx
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Errorf("Unable to load tmx file from unloaded png. Error was: %v", err)
	}
}

func TestTMXTileNotLoadedTempFile(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 132, 99))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	engo.Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "test.png")
	if err = ioutil.WriteFile(tmpfn, imgbuf.Bytes(), 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		ChunkData:   true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("test.png")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load test image from file while loading tmx. Error: %v", err)
	}
}

func TestTMXTileNotLoadedTempFileBadExtensions(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 132, 99))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	engo.Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "test.test")
	if err = ioutil.WriteFile(tmpfn, imgbuf.Bytes(), 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:   "orthogonal",
		RenderOrder:   "right-down",
		BadExtensions: ".test",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("test.test")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Errorf("Able to load test image with bad extension from file while loading tmx.")
	}
}

func TestTMXBadFile(t *testing.T) {
	err := engo.Files.LoadReaderData("bad.tmx", bytes.NewBufferString(badTMX))
	if err == nil {
		t.Error("able to load bad tmx file without an error")
	}
}

func TestTMXGrid(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 457, 305))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	err = engo.Files.LoadReaderData("test.png", imgbuf)
	if err != nil {
		t.Errorf("Unable to load test png. Error was: %v", err)
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "isometric",
		RenderOrder: "left-down",
		Grid:        true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}

	resource, err := engo.Files.Resource("test.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(TMXResource)
	levelData := tmxResource.Level
	if levelData.Orientation != "isometric" {
		t.Errorf("orientation is not isometric. Was: %v", levelData.Orientation)
	}
}

func TestTMXTileImages(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 132, 99))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "left-up",
		Tiles:       true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}
}

func TestTMXTileImagesNotLoadedFileNotExist(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Tiles:       true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("test.png")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Errorf("Able to load tmx file even though assets aren't found.")
	}
}

func TestTMXTileImagesNotLoadedTempFile(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 132, 99))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	engo.Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "test.png")
	if err = ioutil.WriteFile(tmpfn, imgbuf.Bytes(), 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Tiles:       true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("test.png")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load test image from file while loading tmx. Error: %v", err)
	}
}

func TestTMXTileImageWrongFileType(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:   "orthogonal",
		RenderOrder:   "right-down",
		Tiles:         true,
		BadExtensions: ".test",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Error("Able to load tmx when it contains bad image files")
	}
}

func TestTMXBadImageExtension(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:       "orthogonal",
		RenderOrder:       "right-down",
		Tiles:             true,
		BadImageExtension: ".test",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Error("Able to load tmx with bad image layer extension")
	}
}

func TestTMXBadObjectImageExtension(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:             "orthogonal",
		RenderOrder:             "left-down",
		BadObjectImageExtension: ".test",
		ChunkData:               true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("test.png.test")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Error("Able to load tmx with bad object image layer extension")
	}
}

func TestObjectImageNotExistTempFile(t *testing.T) {
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 55, 55))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	dir, err := ioutil.TempDir(".", "testing")
	if err != nil {
		t.Errorf("failed to create temp directory for testing, error: %v", err)
	}
	defer os.RemoveAll(dir)

	engo.Files.SetRoot(dir)

	tmpfn := filepath.Join(dir, "objimgtest.png")
	if err = ioutil.WriteFile(tmpfn, imgbuf.Bytes(), 0666); err != nil {
		t.Errorf("failed to create temp file for testing, file: %v, error: %v", tmpfn, err)
	}

	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:     "orthogonal",
		RenderOrder:     "right-up",
		ObjectImageTest: true,
		ChunkData:       true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("objimgtest.png")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load test image from file while loading object image. Error: %v", err)
	}
}

func TestTMXBadObjectImageNotExist(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation:     "orthogonal",
		RenderOrder:     "left-up",
		ObjectImageTest: true,
		ChunkData:       true,
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	engo.Files.Unload("objimgtest.png")
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err == nil {
		t.Error("Able to load tmx with bad object image layer extension")
	}
}

func TestTMXAsset(t *testing.T) {
	// Create an image in memory
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 457, 305))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	// Load the image
	err = engo.Files.LoadReaderData("test.png", imgbuf)
	if err != nil {
		t.Errorf("Unable to load test png. Error was: %v", err)
	}

	// Load the template
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	// Load tmx
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}

	resource, err := engo.Files.Resource("test.tmx")
	if err != nil {
		t.Errorf("Unable to retrieve resource. Error was: %v", err)
	}
	tmxResource := resource.(TMXResource)
	url := tmxResource.URL()
	if url != "test.tmx" {
		t.Errorf("URL did not match expected\nWanted: %v\nGot: %v", "test.tmx", url)
	}

	// Unload it
	engo.Files.Unload("test.tmx")

	// Check that it's unloaded
	resource, err = engo.Files.Resource("test.tmx")
	if err == nil {
		t.Error("After unloading, resource was retireved.")
	}
}

func TestTMXLevel(t *testing.T) {
	// Create an image in memory
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 457, 305))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	// Load the image
	err = engo.Files.LoadReaderData("test.png", imgbuf)
	if err != nil {
		t.Errorf("Unable to load test png. Error was: %v", err)
	}

	// Load the template
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "orthogonal",
		RenderOrder: "right-down",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	// Load tmx
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}

	resource, err := engo.Files.Resource("test.tmx")
	if err != nil {
		t.Errorf("Unable to retrieve resource. Error was: %v", err)
	}
	tmxResource := resource.(TMXResource)

	bounds := tmxResource.Level.Bounds()
	exp := engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 48, Y: 48},
	}
	if bounds.Min.X != exp.Min.X || bounds.Min.Y != exp.Min.Y || bounds.Max.X != exp.Max.X || bounds.Max.Y != exp.Max.Y {
		t.Errorf("Bounds was not returned correctly\nWanted: %v\nGot: %v", exp, bounds)
	}

	if tmxResource.Level.Width() != 3 {
		t.Error("Level width was not returned correctly.")
	}

	if tmxResource.Level.Height() != 3 {
		t.Error("Level height was not returned correctly.")
	}

	tile := tmxResource.Level.GetTile(engo.Point{X: 20, Y: 20})
	expTile := engo.Point{X: 16, Y: 16}
	if tile.Point.X != expTile.X || tile.Point.Y != expTile.Y {
		t.Errorf("Tile was not returned correctly\nWanted: %v\nGot: %v", expTile, tile.Point)
	}

	if tile.Width() != 16 {
		t.Error("Tile width was not returned correctly")
	}

	if tile.Height() != 16 {
		t.Error("Tile height was not returned correctly")
	}
}

func TestTMXLevelIsometric(t *testing.T) {
	// Create an image in memory
	imgbuf := bytes.NewBuffer([]byte{})
	img := image.NewRGBA(image.Rect(0, 0, 457, 305))
	err := png.Encode(imgbuf, img)
	if err != nil {
		t.Errorf("Unable to encode png from image")
	}

	// Load the image
	err = engo.Files.LoadReaderData("test.png", imgbuf)
	if err != nil {
		t.Errorf("Unable to load test png. Error was: %v", err)
	}

	// Load the template
	buf := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("test").Parse(testTMXtmpl)
	if err != nil {
		t.Error("Error parsing tmx template")
	}
	err = tmpl.Execute(buf, tmxData{
		Orientation: "isometric",
		RenderOrder: "right-down",
	})
	if err != nil {
		t.Error("Error executing tmx template")
	}

	// Load tmx
	err = engo.Files.LoadReaderData("test.tmx", buf)
	if err != nil {
		t.Errorf("Unable to load tmx file for testing. Error was: %v", err)
	}

	resource, err := engo.Files.Resource("test.tmx")
	if err != nil {
		t.Errorf("Unable to retrieve resource. Error was: %v", err)
	}
	tmxResource := resource.(TMXResource)

	bounds := tmxResource.Level.Bounds()
	exp := engo.AABB{
		Min: engo.Point{X: -16, Y: 0},
		Max: engo.Point{X: 32, Y: 56},
	}
	if bounds.Min.X != exp.Min.X || bounds.Min.Y != exp.Min.Y || bounds.Max.X != exp.Max.X || bounds.Max.Y != exp.Max.Y {
		t.Errorf("Bounds was not returned correctly\nWanted: %v\nGot: %v", exp, bounds)
	}

	tile := tmxResource.Level.GetTile(engo.Point{X: 20, Y: 20})
	expTile := engo.Point{X: 16, Y: 16}
	if tile.Point.X != expTile.X || tile.Point.Y != expTile.Y {
		t.Errorf("Tile was not returned correctly\nWanted: %v\nGot: %v", expTile, tile.Point)
	}
}
