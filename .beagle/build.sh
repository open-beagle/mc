#!/bin/sh

set -x

export GOARCH=amd64
make build
mkdir -p .bin/linux/$GOARCH
mv mc .bin/linux/$GOARCH/mc

export GOARCH=arm64
make build
mkdir -p .bin/linux/$GOARCH
mv mc .bin/linux/$GOARCH/mc

export GOARCH=ppc64le
make build
mkdir -p .bin/linux/$GOARCH
mv mc .bin/linux/$GOARCH/mc
