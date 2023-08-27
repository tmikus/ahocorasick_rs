#!/usr/bin/env bash

# Get script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SRC_DIR="$DIR/.."

(cd "$SRC_DIR" && go test -bench=. -ldflags="-r $SRC_DIR/lib")
