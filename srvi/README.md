# SERVi

SERVi is a utility for quickly testing out your games in the browser. It supports hosting a server, displaying errors, and rebuilding your game each time you refesh the page.

## Install

`go get -u github.com/ajhager/engi/srvi`

## Usage

Run `srvi` in the same directory as your game, with your static files in a directory named 'data'. Access http://localhost:8080/ if your game file is at ./main.go. Any other file name can be accessed at http://localhost:8080/name, where 'name' would be name.go.

`
Usage of srvi:
	-host="127.0.0.1": The host at which to serve your games
	-port=8080: The port at which to serve your games
	-static="data": The relative path to your assets
`
