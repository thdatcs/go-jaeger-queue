#!/bin/bash
GOPATH=${PWD}

cd src
for pkg in $@
do
    echo "golint $pkg"
    sources=$(go list $pkg/...)
    ${GOPATH}/bin/golint ${sources}
done
cd ..