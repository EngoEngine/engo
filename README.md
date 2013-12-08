eng 2D game library for go
===

eng depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install libglfw-dev libxrandr-dev
* OSX: brew tap versions -> brew install glfw2
* Windows: download the [glfw2](http://sourceforge.net/projects/glfw/files/glfw/2.7.7/) binaries, then drop the GL directory into C:\MinGW\include (64bit: C:\MinGW\mingw-w64-x86_6\include) and the library files into C:\MinGW\lib (64bit: C:\MinGW\mingw-w64-x86_6\lib). You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
`go get github.com/ajhager/eng`

Documentation
-------------
[godoc.org](http://godoc.org/github.com/ajhager/eng)
