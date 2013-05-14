eng 2D game library for go
===

+NOTHING ABOUT THIS LIBRARY IS STABLE+

eng depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw).
* Ubuntu: apt-get install libglfw-dev
* OSX: brew install glfw
* Windows: download the glfw binaries, then drop the GL directory into C:\MinGW\include and the files for your arch under libmingw into C:\MinGW\lib. You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
`go get github.com/ajhager/eng`

TODO
----

* Port over 2d shaders and effect composer
* Port over particle system
* Create Cache, a static Batch
* Shape renderer
* Investigate audio solutions
* Add texture packer support
* Tiled map support
* Documentation and website
