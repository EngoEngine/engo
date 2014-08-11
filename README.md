# [ENGi v0.4.0](http://ajhager.com/engi)

A multi-platform 2D game library for Go.

## Status

Expect bugs and major API changes. Just a proof of concept at the moment.

## OSX

The darwin backend depends on [glfw3](http://github.com/go-gl/glfw). You can install it using homebrew:

```bash
brew tap homebrew/versions; brew install glfw3
```

## Linux

On Ubuntu you will need libx11 and libxcb installed. I will add an exact apt-get command soon.

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
