#!/bin/bash
ARCH=$1

echo "go build $ARCH"
GOARCH=$ARCH go build -o fever-pass-$ARCH

cd ..
tar -Jcvf fever-pass/fever-pass-$ARCH.tar.xz fever-pass/fever-pass fever-pass/static fever-pass/templates
cd fever-pass 
go clean