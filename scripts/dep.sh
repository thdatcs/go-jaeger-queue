#!/bin/bash
GOPATH=${PWD}

for pkg in $@
do
    echo "dep $pkg"
    cd $GOPATH/src/$pkg && $GOPATH/bin/dep ensure -v && cd $GOPATH
done