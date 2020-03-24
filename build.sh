#!/bin/bash
OS=$1
ARCH=$2

echo "go build $ARCH"
GOOS=$OS GOARCH=$ARCH go build -o fever-pass

mkdir build
cd ..
tar -Jcvf fever-pass/build/fever-pass-$OS-$ARCH.tar.xz fever-pass/fever-pass fever-pass/static fever-pass/templates fever-pass/doc
cd fever-pass 
go clean