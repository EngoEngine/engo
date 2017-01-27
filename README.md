# Engo
[![GoDoc](https://godoc.org/engo.io/engo?status.svg)](https://godoc.org/engo.io/engo)
[![Join the chat at https://gitter.im/EngoEngine/engo](https://badges.gitter.im/EngoEngine/engo.svg)](https://gitter.im/EngoEngine/engo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) ![License](https://img.shields.io/badge/License-MIT-blue.svg)
[![Build Status](https://travis-ci.org/EngoEngine/engo.svg?branch=master)](https://travis-ci.org/EngoEngine/engo)
[![Build status](https://ci.appveyor.com/api/projects/status/019qc8hncmhnje83?svg=true)](https://ci.appveyor.com/project/otraore/engo)
[![Go Report Card](https://goreportcard.com/badge/engo.io/engo)](https://goreportcard.com/report/engo.io/engo)
[![Coverage Status](https://coveralls.io/repos/github/EngoEngine/engo/badge.svg?branch=master)](https://coveralls.io/github/EngoEngine/engo?branch=master)

A cross-platform game engine written in Go following an interpretation of the Entity Component System paradigm. Engo is
currently compilable for Mac OSX, Linux and Windows. With the release of Go 1.4, supporting Android and the inception of
iOS compatibility, mobile has been be added as a release target. Web support
([gopherjs](https://github.com/gopherjs/gopherjs)) is also available.

This table shows the current state of what's working and what's not. We are aiming to turn all **No** columns into **Yes** columns.

| Aspect/Target | Linux | Windows | OSX | WebGL | Android | iOS* |
| ------------- | ----- | ------- | --- | ----- | ------- | ---  |
| Compilation   | Yes   | Yes     | Yes | Yes   | Yes     | ?    |
| AudioSystem   | Yes   | No      | Yes | No    | No      | ?    |
| RenderSystem  | Yes   | Yes     | Yes | Yes   | Yes     | ?    |
| XXXXXSystem   | Yes   | Yes     | Yes | Yes   | Yes     | ?    |

*iOS will not build at all (due to `gomobile` limitations), if it would, most of these would be supported the same way Android is.

We are working towards a `1.0` [release](https://github.com/EngoEngine/engo/milestones/1.0), and making good progress.
We welcome any contributions.

Currently documentation is pretty scarce, this is because we have not *completely* finalized the API and are about to
go through a "prettification" process in order to increase elegance and usability. For a basic up-to-date example of
most features, look at the demos.

## Getting in touch / Contributing

We have a [gitter](https://gitter.im/EngoEngine/engo) chat for people to join who want to further discuss `engo`. We are happy to discuss bugs, feature requests and would love to hear about the projects you are building!

## Getting Started

### Theory: `common` vs `engo`

There are currently two major important packages within this repository: `engo.io/engo` and `engo.io/engo/common`.

The top level `engo` package contains the functionality of creating windows, starting the game, creating an OpenGL
context and handling input. It is designed to be used with Systems designed as per `engo.io/ecs` specifications.
The `common` package contains our ECS implementations of common game development Systems like a  `RenderSystem` or
`CameraSystem`.

### Practice: Getting it to Run

1. First, you have to install some dependencies:
  1. If you're running on Debian/Ubuntu:
    `sudo apt-get install libopenal-dev libglu1-mesa-dev freeglut3-dev mesa-common-dev xorg-dev libgl1-mesa-dev`
  2. If you're running on Windows you'll need a gcc compiler that the go tool can use and have `gcc.exe` in your PATH environmental variable. We recommend [Mingw](http://mingw-w64.org/doku.php/start) since it has been tested.
  3. If you're on OSX, you should be OK. [Open an issue if you are not](https://github.com/EngoEngine/engo/issues/new)
2. Then, you can go get it:
`go get -u engo.io/engo`
3. Now, you have two choices:
  1. Visit [our website](https://engo.io/), which hosts a full-blown tutorial series on how to create your own game, and on top of that, has some conceptual explanations;
  2. Check out some demos in our [demos folder](https://github.com/EngoEngine/engo/tree/master/demos).
4. Finally, if you run into problems, if you've encountered a bug, or want to request a feature, feel free to shoot
us a DM or [create an issue](https://github.com/EngoEngine/engo/issues/new).

## Breaking Changes
Engo is currently undergoing a lot of optimizations and constantly gets new features. However, this sometimes means things break. In order to make transitioning easier for you,
we have a list of those changes, with the most recent being at the top. If you run into any problems, please contact us at [gitter](https://gitter.im/EngoEngine/engo).

* `engo.PreloadedSpriteSingle` is now `engo.LoadedSprite`
* `engo.Files.Load` and `engo.Files.LoadMany` have been merged into one function `engo.Files.Load` which does the same thing as `engo.Files.LoadMany` allowing an indefinite ammount of parameters to be passed in.
* `engo` has been split in `engo` (which contains stuff about creating windows, starting the game, creating an OpenGL context, input handling, etc.) - and `common` (which contains a lot of common `System` implementations for common tasks (`RenderSystem`, `CameraSystem`, `AudioSystem`, etc.)
* `engo.Width()` and `engo.Height()` have been changed to `engo.GameWidth()` and `engo.GameHeight()` respectively.
* `RenderComponent.Scale` is now no longer a method, but a variable you can change / access directly.
* `engo.NewRenderComponent` was removed. You can now define the values you want directly by using `engo.RenderComponent{}`. Note that the `Drawable` is still required.
* `ecs.Entity` changed to `ecs.BasicEntity`, `world.AddEntity` is gone - **a lot** has changed here. The entire issue is described [here](https://github.com/EngoEngine/ecs/issues/13), while [this comment](https://github.com/EngoEngine/ecs/issues/13#issuecomment-210887914) in particular, should help you migrate your code.
* Renamed `engo.io/webgl` to `engo.io/gl`, because the package handles more than only *web*gl.
* `github.com/EngoEngine/engo` -> `engo.io/engo` - Our packages `engo`, `ecs` and `gl` should now be imported using the `engo.io` path.
* `engi.XXX` -> `engo.XXX` - We renamed our package `engi` to `engo`.

## History

Engo, originally known as `Engi` was written by [ajhager](https://github.com/ajhager) as a general purpose Go game engine. With a desire to build it into an "ECS" game engine, it was forked to `github.com/paked/engi`. After passing through several iterations, it was decided that the project would be rebranded and rereleased as Engo on its own GitHub organisation.

## Credits

Thank you to everyone who has worked on, or with `Engo`. None of this would be possible without you, and your help has been truly amazing.

- [ajhager](https://github.com/ajhager): Building the original `engi`, which engo was based off of
- [paked](https://github.com/paked): Adding ECS element, project maintenance and management
- [Newbrict](https://github.com/Newbrict): Font rendering, TMX support
- [EtienneBruines](https://github.com/EtienneBruines): Rewriting the OpenGL code, maintenance and helping redesign the API
- [otraore](https://github.com/otraore): Adding in GopherJS support, maintenance
- [Everyone else who has submitted PRs or issues over the years, to any iteration of the project](https://github.com/EngoEngine/engo/graphs/contributors)

These are 3rd party projects that have made `engo` possible.
- The original [engi](https://github.com/ajhager/engi) game engine which engo was based off of ([BSD license](https://github.com/ajhager/engi/blob/master/LICENSE))
