#!/usr/bin/env bash

docker build -t ahocorasick-rs-builder .

docker run -it --rm -v "$(pwd)":/app/src ahocorasick-rs-builder

go build -ldflags="-r $(pwd)/lib" search.go