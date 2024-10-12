#!/bin/sh

git config --global --add safe.directory $PWD

set -x

export TARGETPLATFORM=linux/amd64
xx-apk add perl
xx-go --wrap
make build
mkdir -p .bin/linux/$GOARCH
mv mc .bin/linux/$GOARCH/mc

export TARGETPLATFORM=linux/arm64
xx-apk add perl
xx-go --wrap
make build
mkdir -p .bin/linux/$GOARCH
mv mc .bin/linux/$GOARCH/mc
