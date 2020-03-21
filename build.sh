#!/bin/bash
ARCH=$1

echo "go build $ARCH"
GOARCH=$ARCH go build -o fever-pass

mkdir build
cd ..
tar -Jcvf fever-pass/build/fever-pass-$ARCH.tar.xz fever-pass/fever-pass fever-pass/static fever-pass/templates fever-pass/doc
cd fever-pass 
go clean