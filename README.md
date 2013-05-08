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

* Wrap up framebuffers
* Port over 2d shaders and effect composer
* Port over particle system
* Create Cache, a static Batch
* Add true bmfont support
* Add Orthographic camera support
* Add simple 2d scenegraph that wraps a camera and batch
