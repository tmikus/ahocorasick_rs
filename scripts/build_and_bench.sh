#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

"$DIR/build.sh" || exit 1
(cd "$DIR/.." && go test -v .) || exit 1
"$DIR/bench.sh" || exit 1
