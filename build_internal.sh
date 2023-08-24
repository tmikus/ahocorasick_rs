#!/usr/bin/env bash

TARGETS=(
  "aarch64-unknown-linux-gnu"
  "x86_64-unknown-linux-gnu"
  "aarch64-apple-darwin"
  "x86_64-apple-darwin"
  "x86_64-pc-windows-gnu"
)

cd src/lib/ahocorasick_rs
for TARGET in "${TARGETS[@]}"; do
  echo "Building ahocorasick_rs for $TARGET..."
  cargo build --release --target=$TARGET
done

# Go back to lib directory
cd ..
rm -rf darwin linux windows
mkdir -p darwin linux windows

cp ahocorasick_rs/target/aarch64-unknown-linux-gnu/release/libahocorasick_rs.a linux/libahocorasick_rs_arm64.a
cp ahocorasick_rs/target/x86_64-unknown-linux-gnu/release/libahocorasick_rs.a linux/libahocorasick_rs_amd64.a
cp ahocorasick_rs/target/aarch64-apple-darwin/release/libahocorasick_rs.a darwin/libahocorasick_rs_arm64.a
cp ahocorasick_rs/target/x86_64-apple-darwin/release/libahocorasick_rs.a darwin/libahocorasick_rs_amd64.a
cp ahocorasick_rs/target/x86_64-pc-windows-gnu/release/libahocorasick_rs.a windows/libahocorasick_rs_amd64.a
