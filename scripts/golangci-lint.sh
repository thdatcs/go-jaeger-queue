#!/bin/bash
GOPATH=${PWD}

cd src
for pkg in $@
do
    echo "golangci-lint $pkg"
    sources=$(go list $pkg/...)
    ${GOPATH}/bin/golangci-lint run ${sources}
done
cd ..