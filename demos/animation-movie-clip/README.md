# Animation form `MovieClip` format Demo

## What does it do?

It demonstrates how one can use animations from `MovieClip` format.
  
For doing so, it exported an animation in `MovieClip` format form Editor eg DragonBones (probably not only this editor supports this format):
`File > Export` and for param:
 
 * `Data config > Type` select: `Egret MC`
 * `Texture Config > Image Type` select: `Texture Atlas` and in `Settings` select: Region Padding X: 50, Y:50

> **note** support only `scale textures` and not support `multiple textures` yet.

After export is complete, two files will be created: `name_mc.json` and `name_tex.png`.
Rename `name_mc.json` in `name.mc.json` and place it assets folder.


## What are important aspects of the code?

These lines are key in this demo:

```go
func (*DefaultScene) Preload() {
	engo.Files.Load("sheep.mc.json")
}

func (scene *DefaultScene) Setup(u engo.Updater) {
    ...
	mcr, _ := mc.LoadResource("sheep.mc.json")
	hero := NewHeroEntity(engo.Point{200, 50}, mcr)
    ...
}
```

## What can I do in this demo?

* You can random changes animation (press `Enter`)
