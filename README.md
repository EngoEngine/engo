# [ENGi v0.6.0](http://ajhager.com/engi)

A multi-platform 2D game library for Go.

## Status

*SUPER ALPHA* Expect bugs and major API changes. Just a proof of concept at the moment.

		* Clean native interfaces
		* Better asset management
		* More consistent and powerful input
		* Support multiple windows / games
		* Better windowed / borderless windowed / fullscreen support
		* Add support for built-in colors and images

## Desktop

The desktop backend depends on glfw3, but includes the source code and links it statically. If you are having linker errors on Windows, I suggest using [TDM-GCC](http://tdm-gcc.tdragon.net/download) instead of MinGW as your cgo compiler.

## Web

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs). ```gopherjs build``` is very much like ```go build```, then you can embed the resulting javascript file into your html document.

During development you can use [SRVi](https://github.com/ajhager/srvi) to automatically rebuild and serve your project every time you refresh. Quickly try out new ideas without even needing to setup a new index.html every time. 

## Android

The android backend is in the works, following the daily updates to the go.mobile repo.

## Install

```bash
go get -u github.com/ajhager/engi
```

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/engi)
