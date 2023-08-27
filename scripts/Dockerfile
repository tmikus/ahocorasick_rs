FROM rust:1.71.1

WORKDIR /app

# installing rust toolchains
# Linux
RUN rustup target add aarch64-unknown-linux-gnu
RUN rustup target add x86_64-unknown-linux-gnu
# Apple
RUN rustup target add aarch64-apple-darwin
RUN rustup target add x86_64-apple-darwin
# Windows
RUN rustup target add x86_64-pc-windows-gnu

COPY build_internal.sh .
RUN chmod +x build_internal.sh

ENTRYPOINT ["/app/build_internal.sh"]