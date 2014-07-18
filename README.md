# eng v0.2.0

A 2D game library for go. Expect bugs and major API changes. Just a proof of concept at the moment.

## Desktop Build

The desktop backend depends on [glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install glfw3
* OSX: brew tap homebrew/versions; brew install glfw3
* Windows: download the [glfw3](http://www.glfw.org/docs/latest/) binaries, then drop the GL directory into C:\MinGW\include (64bit: C:\MinGW\mingw-w64-x86_6\include) and the library files into C:\MinGW\lib (64bit: C:\MinGW\mingw-w64-x86_6\lib). You will then need to install glfw.dll system wide or have it in the directory with your game.

## Web Build

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs). eng comes with a utility for quickly testing out your games in the browser.

`go get github.com/ajhager/eng/bld`

Run `bld` in the same directory as your game, with your static files in a directory named 'data'. Access http://localhost:8080/main if your game file is at ./main.go.

## Install

`go get github.com/ajhager/eng`

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/eng)
