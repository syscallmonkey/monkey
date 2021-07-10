#!/bin/bash

export VERSION=${VERSION:-`git rev-parse HEAD`}
export BUILD=${BUILD:-`date`}

go build \
    -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=$BUILD'" \
    -o ./bin/${BIN} \
    ./cmd/${BIN}