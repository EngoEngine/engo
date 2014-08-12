# [ENGi v0.4.0](http://ajhager.com/engi)

A multi-platform 2D game library for Go.

## Status

*SUPER ALPHA* Expect bugs and major API changes. Just a proof of concept at the moment.

## OSX

The darwin backend depends on [glfw3](http://github.com/go-gl/glfw). You can install it using homebrew:

```bash
brew tap homebrew/versions; brew install glfw3
```

## Linux

Ubuntu 14.04:

```bash
sudo apt-get install build-essential git mesa-common-dev libx11-dev libx11-xcb-dev libxcb-icccm4-dev libxcb-image0-dev libxcb-randr0-dev libxcb-render-util0-dev libxcb-xkb-dev libfreetype6-dev libbz2-dev
```

Arch:

```bash
pacman -Sy base-devel git mesa libx11 libxcb xcb-util-wm xcb-util-image libxrandr xcb-util-renderutil libxkbcommon-x11 freetype2 bzip2
```

## Windows

As long as you have mingw installed correctly for CGO, everything should work out of the box.

## Web

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs). Check out the [SERVi](http://github.com/ajhager/engi/tree/master/srvi) utility for testing your games in the browser.

## Install

```bash
go get -u github.com/ajhager/engi
```

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/engi)
