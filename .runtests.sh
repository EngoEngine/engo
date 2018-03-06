#!/usr/bin/env bash

echo "Using GOPATH=$GOPATH"

echo "Getting engo.io/engo using 'go get'"
go get -t -v ./... || exit 1

# These can fail without us minding it
blacklist="engo.io/engo/demos/demoutils"

if [ "$TEST_TYPE" == "linux_test" ]
then
    echo "Testing engo.io/engo using coveralls"
    $HOME/gopath/bin/goveralls -service=travis-ci

    echo "Testing and benchmarking engo.io/engo"
    go test -v -bench=. ./... || exit 1

    echo "Checking for unnecessary conversions using unconvert"
    unconvert -v engo.io/engo
elif [ "$TEST_TYPE" == "linux_build" ]
then
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

        # Creating the output directory, attempting to build and exit 1 if it failed
        mkdir -p "$outdir/linux/"
        go build -o "$outdir/linux/${dir}" -tags demo ${dir} || exit 1
    done
elif [ "$TEST_TYPE" == "js_test" ]
then
    # TODO: Fix the build so this actually passes
    # echo "Testing engo.io/engo using 'gopherjs test'"
    # gopherjs test
    echo "Skipping tests for engo.io/engo using 'gopherjs test' (won't pass)"
elif [ "$TEST_TYPE" == "js_build" ]
then
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

        # Creating the output directory, attempting to build and exit 1 if it failed
        mkdir -p "$outdir/gopherjs/"
        gopherjs build -o "$outdir/gopherjs/${dir}" -tags demo ${dir} || exit 1
    done
elif [ "$TEST_TYPE" == "android_test" ]
then
    echo "Skipping tests for engo.io/engo using 'gomobile' (no tools exist yet)"
elif [ "$TEST_TYPE" == "android_build" ]
then
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

        # Creating the output directory, attempting to build and exit 1 if it failed
        mkdir -p `dirname "$outdir/android/${dir}.apk"`
        gomobile build -o "$outdir/android/${dir}.apk" -target android -tags demo ${dir} || exit 1
    done
else
    echo "environment variable TEST_TYPE was not set"
fi

# Test the TrafficManager as well
# TODO