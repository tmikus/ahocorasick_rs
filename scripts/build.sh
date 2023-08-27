#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SRC_DIR="$DIR/.."

docker build -t ahocorasick-rs-builder ..

docker run -it --rm -v "$SRC_DIR":/app/src ahocorasick-rs-builder

(cd "$SRC_DIR" && go build -ldflags="-r $SRC_DIR/lib" .)