#!/bin/bash

DIR="fever-pass"
DEST="build"

clean() {
    rm -rf $DIR
}

mk_work_dir() {
    mkdir -p $DIR
    cp -r templates $DIR/
    cp -r static $DIR/
    cp -r doc $DIR/
    cp -r testdata $DIR/
    cp LICENSE $DIR/
    cp README.md $DIR/
    mkdir -p $DEST
}

build() {
    OS=$1
    ARCH=$2

    echo "building $OS $ARCH"
    GOOS=$OS GOARCH=$ARCH go build -o $DIR/fever-pass

    echo "packaging $OS $ARCH"
    tar -Jcf $DEST/fever-pass_$OS-$ARCH.tar.xz $DIR
}

clean
mk_work_dir
build linux amd64
build linux 386
build linux arm
build windows amd64
build windows 386
build darwin amd64
build darwin 386
clean