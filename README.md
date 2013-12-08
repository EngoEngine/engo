eng 2D game library for go
===

eng depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install libglfw-dev libxrandr-dev
* OSX: brew tap versions -> brew install glfw2
* Windows: download the glfw2 binaries, then drop the GL directory into C:\MinGW\include and the files for your arch under libmingw into C:\MinGW\lib. You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
`go get github.com/ajhager/eng`

Documentation
-------------
[godoc.org](http://godoc.org/github.com/ajhager/eng)
