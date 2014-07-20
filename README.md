# [ENG v0.2.0](http://ajhager.com/engi)

A 2D game library for go. Expect bugs and major API changes. Just a proof of concept at the moment.

## Desktop Build

The desktop backend depends on [glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install glfw3
* OSX: brew tap homebrew/versions; brew install glfw3
* Windows: download the [glfw3](http://www.glfw.org/docs/latest/) binaries, then drop the GL directory into C:\MinGW\include (64bit: C:\MinGW\mingw-w64-x86_6\include) and the library files into C:\MinGW\lib (64bit: C:\MinGW\mingw-w64-x86_6\lib). You will then need to install glfw.dll system wide or have it in the directory with your game.

## Web Build

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs). eng comes with a utility for quickly testing out your games in the browser.

`go get github.com/ajhager/engi/srvi`

Run `srvi` in the same directory as your game, with your static files in a directory named 'data'. Access http://localhost:8080/ if your game file is at ./main.go. Any other file name can be accessed at http://localhost:8080/name, where 'name' would be name.go.

You can supply a custom flags to srvi:

`
Usage of srv:
	-host="127.0.0.1": The host at which to serve your games
	-port=8080: The port at which to serve your games
	-static="data": The relative path to your assets
`

## Install

`go get github.com/ajhager/engi`

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/engi)
