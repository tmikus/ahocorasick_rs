#!/usr/bin/env bash

(cd lib/ahocorasick_rs && cargo build --release)
cp lib/ahocorasick_rs/target/release/libahocorasick_rs.so lib/
go run -ldflags="-r $(pwd)/lib" search.go
