# [ENGi v0.6.0](http://ajhager.com/engi)

A multi-platform 2D game library for Go.

## Status

*SUPER ALPHA* Expect bugs and major API changes. Just a proof of concept at the moment.

		* Support for custom mappings on fonts
		* Better asset management
		* Add support for built-in colors and images
		* Clean native interfaces
		* More consistent and powerful input
		* Support multiple windows / games
		* Better windowed / borderless windowed / fullscreen support

## Desktop

The desktop backend depends on glfw3, but  includes the source code and links it statically.

## Web

The web backend depends on [gopherjs](http://github.com/neelance/gopherjs).

## Android

The android backend is in the works, following the daily updates to the go.mobile repo.

## Install

```bash
go get -u github.com/ajhager/engi
```

## Documentation

[godoc.org](http://godoc.org/github.com/ajhager/engi)
