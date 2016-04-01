# Engo
[![Join the chat at https://gitter.im/EngoEngine/engo](https://badges.gitter.im/EngoEngine/engo.svg)](https://gitter.im/EngoEngine/engo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A cross-platform game engine written in Go following an interpretation of the Entity Component System paradigm. Engo is
 currently compilable for Mac OSX, Linux and Windows. With the release of Go 1.4, sporting Android and the inception of 
 iOS compatibility, mobile will [soon](https://github.com/engoengine/engo/issues/63) be added as a release target. Web 
 support  ([gopherjs](https://github.com/gopherjs/gopherjs)) is also [planned](https://github.com/engoengine/engo/issues/71). 

Currently documentation is pretty scarce, this is because we have not *completely* finalized the API and are about to 
go through a "prettification" process in order to increase elegance and usability. For a basic up-to-date example of 
most features, look at the demos.

## Getting in touch / Contributing
Currently we are active on IRC / Freenode at the `#engo` channel. You can also create an issue to start a discussion. 

## Getting Started

1. First, you have to install some dependencies if you're running on Debian/Ubuntu: 
`sudo apt-get install libopenal-dev libglu1-mesa-dev freeglut3-dev mesa-common-dev xorg-dev libgl1-mesa-dev`
2. Then, you can go get it:
`go get -u github.com/engoengine/engo`
3. Now, you have two choices:
  1. Read the [Wiki: Getting Started](https://github.com/engoengine/engo/wiki/Getting-Started), for an explanation on the basics;
  2. Check out some demos in our [demos folder](https://github.com/engoengine/engo/tree/master/demos). 
4. Finally, if you run into problems, if you've encountered a bug, or want to request a feature, feel free to shoot 
us a DM or [create an issue](https://github.com/engoengine/engo/issues/new). 
