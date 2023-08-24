#!/usr/bin/env bash

go test -bench=. -ldflags="-r $(pwd)/lib"
