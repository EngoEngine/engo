#!/bin/bash

test_demo () {
	if [ "$1" = "demoutils" ]; then
		echo Skipping $1
		return
	fi
	cd $1
	echo Verifying $1
	go build -tags demo -o $1
	cd ..
}

go test -tags appveyor -v ./...

cd demos
for D in *; do [ -d "${D}" ] && test_demo "${D}"; done

