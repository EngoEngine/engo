# Engo
[![Join the chat at https://gitter.im/EngoEngine/engo](https://badges.gitter.im/EngoEngine/engo.svg)](https://gitter.im/EngoEngine/engo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A cross-platform game engine written in Go following an interpretation of the Entity Component System paradigm. Engo is
currently compilable for Mac OSX, Linux and Windows. With the release of Go 1.4, sporting Android and the inception of
iOS compatibility, mobile will [soon](https://github.com/EngoEngine/engo/issues/63) be added as a release target. Web
support  ([gopherjs](https://github.com/gopherjs/gopherjs)) is also [planned](https://github.com/EngoEngine/engo/issues/71).

Currently documentation is pretty scarce, this is because we have not *completely* finalized the API and are about to
go through a "prettification" process in order to increase elegance and usability. For a basic up-to-date example of
most features, look at the demos.

## Getting in touch / Contributing

We have a [gitter](https://gitter.im/EngoEngine/engo) chat for people to join who want to further discuss `engo`. We are happy to discuss bugs, feature requests and would love to hear about the projects you are building!

## Getting Started

1. First, you have to install some dependencies if you're running on Debian/Ubuntu:
`sudo apt-get install libopenal-dev libglu1-mesa-dev freeglut3-dev mesa-common-dev xorg-dev libgl1-mesa-dev`
2. Then, you can go get it:
`go get -u engo.io/engo`
3. Now, you have two choices:
  1. Read the [Wiki: Getting Started](https://github.com/EngoEngine/engo/wiki/Getting-Started), for an explanation on the basics;
  2. Check out some demos in our [demos folder](https://github.com/EngoEngine/engo/tree/master/demos).
4. Finally, if you run into problems, if you've encountered a bug, or want to request a feature, feel free to shoot
us a DM or [create an issue](https://github.com/EngoEngine/engo/issues/new).

## History

Engo, originally known as `Engi` was written by [ajhager](https://github.com/ajhager) as a general purpose Go game engine. With a desire to build it into an "ECS" game engine, it was forked to `github.com/paked/engi`. After passing through several iterations, it was decided that the project would be rebranded and rereleased as Engo on its own GitHub organisation.

## Credits

Thank you to everyone who has worked on, or with `Engo`. Non of this would be possible without you, and your help has been truly amazing.

- [ajhager](https://github.com/ajhager): Building the original `engi`, which engo was based off of
- [paked](https://github.com/paked): Adding ECS element, project maintenance and management
- [EtienneBruines](https://github.com/EtienneBruines): Rewriting the OpenGL code, maintenance and helping redesign the API
- [Everyone else who has submitted PRs over the years, to any iteration of the project](https://github.com/EngoEngine/engo/graphs/contributors)
