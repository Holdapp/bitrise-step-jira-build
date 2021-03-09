#!/bin/bash

VERSION="1.1.0"

docker build -t libgit2-build:${VERSION} .

mkdir -p products/
docker run -v $(pwd)/products:/products libgit2-build:${VERSION}
