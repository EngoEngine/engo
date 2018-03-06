#!/usr/bin/env bash

go get golang.org/x/tools/cmd/cover
go get github.com/mdempsky/unconvert
go get github.com/mattn/goveralls

if [ "$TEST_TYPE" == "js_test" ] || [ "$TEST_TYPE" == "js_build" ]
then
    go get github.com/gopherjs/gopherjs
    gopherjs get "honnef.co/go/js/dom"
    gopherjs get "honnef.co/go/js/xhr"
elif [ "$TEST_TYPE" == "android_test" ] || [ "$TEST_TYPE" == "android_build" ]
then
    go get golang.org/x/mobile/cmd/gomobile
    gomobile init
    git clone https://github.com/Noofbiz/android-ndk.git $HOME/android-ndk-root
    printf "$HOME/android-ndk-root" > $GOPATH/pkg/gomobile/android_ndk_root
elif [ "$TEST_TYPE" == "linux_test" ] || [ "$TEST_TYPE" == "linux_build" ]
then
    echo "nothing more to add for linux"
else
    echo "environment variable TEST_TYPE was not set"
fi
