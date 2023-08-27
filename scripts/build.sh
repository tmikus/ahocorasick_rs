#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SRC_DIR="$DIR/.."

docker build -t ahocorasick-rs-builder -f "$DIR/Dockerfile" "$DIR"

docker run \
  -it \
  --rm \
  -v /etc/passwd:/etc/passwd:ro \
  -v /etc/group:/etc/group:ro \
  -v "$SRC_DIR":/app/src \
  -u $(id -u):$(id -g) \
  ahocorasick-rs-builder

(cd "$SRC_DIR" && go build -ldflags="-r $SRC_DIR/lib" .)