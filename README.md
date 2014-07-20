# [ENGi v0.3.0](http://ajhager.com/engi)

A multi-platform 2D game library for Go.

## Status

Expect bugs and major API changes. Just a proof of concept at the moment.

## Install

`go get -u github.com/ajhager/engi`

## Desktop

The desktop backend depends on [glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install glfw3
* OSX: brew tap homebrew/versions; brew install glfw3
* Windows: download the [glfw3](http://www.glfw.org/docs/latest/) binaries, then drop the GL directory into C:\MinGW\include (64bit: C:\MinGW\mingw-w64-x86_6\include) and the library files into C:\MinGW\lib (64bit: C:\MinGW\mingw-w64-x86_6\lib). You will then need to install glfw.dll system wide or have it in the directory with your game.

## Web

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs). Check out the [SERVi](http://github.com/ajhager/engi/srvi) utility for trying out your games in the browser.

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/engi)
