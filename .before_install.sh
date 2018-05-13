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
    git clone https://github.com/golang/mobile.git $GOPATH/src/golang.org/x/mobile
    cd $GOPATH/src/golang.org/x/mobile/cmd/gomobile
    git reset --hard 598bfe4b20d39a660581f014b68e60c5ad425336
    go install
    cd ~
    gomobile init
    git clone https://github.com/Noofbiz/android-ndk.git $HOME/android-ndk-root
    printf "$HOME/android-ndk-root" > $GOPATH/pkg/gomobile/android_ndk_root
elif [ "$TEST_TYPE" == "linux_test" ] || [ "$TEST_TYPE" == "linux_build" ]
then
    echo "nothing more to add for linux"
elif [ "$TEST_TYPE" == "traffic_manager" ]
then
    go get github.com/EngoEngine/TrafficManager
else
    echo "environment variable TEST_TYPE was not set"
fi
