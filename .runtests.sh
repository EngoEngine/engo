#!/usr/bin/env bash

# Install dependencies
echo "Installing gopherjs ..."
go get github.com/gopherjs/gopherjs

# TODO: These few lines are required until https://github.com/gopherjs/gopherjs/issues/455 is fixed.
echo "Installing engo_js dependencies ..."
gopherjs get "honnef.co/go/js/dom"
gopherjs get "honnef.co/go/js/xhr"

echo "Installing gomobile ..."
go get golang.org/x/mobile/cmd/gomobile

echo "Initializing gomobile ..."
gomobile init

echo "Using GOPATH=$GOPATH"

# These can fail without us minding it
blacklist="engo.io/engo/demos/demoutils,engo.io/engo/demos/tilemap"

for dir in `pwd`/demos/*/
do
    # Formatting the directory to be usable by Go
    dir=${dir%*/}
    dir=${dir#$GOPATH/src/}

    # Ignore the directory if it's in the blacklist
    if [[ $blacklist == *"${dir}"* ]]
    then
        echo "Skipping ${dir}"
        continue
    fi

    # Some debug output and output directory initialization
    echo "Verifying ${dir} ..."
    outdir="/tmp/go-builds"

    # Per build method, creating the output directory, attempting to build/test and exit 1 if it failed

    mkdir -p "$outdir/native/"
    go build -o "$outdir/native/${dir}" ${dir} || exit 1
    go test ${dir} || exit 1

    mkdir -p "$outdir/gopherjs/"
    echo "Gopherjs get is running now"
    gopherjs get ${dir} || exit 1
    echo "Gopherjs get is done now"
    gopherjs build -o "$outdir/gopherjs/${dir}" ${dir} || exit 1

    mkdir -p `dirname "$outdir/android/${dir}.apk"`
    gomobile build -o "$outdir/android/${dir}.apk" -target android ${dir} || exit 1

done

# Test the TrafficManager as well
# TODO