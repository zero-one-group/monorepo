# Multi-stage Rust build for the rust-ai service.
#
# Stage 1: build the release binary inside a full Rust image.
# Stage 2: copy the binary into a minimal debian-slim runtime that has
#          just enough libc for sqlx + rustls + tokio.
#
# Built from the monorepo root so the workspace Cargo.lock is visible
# to the build context. The moon `docker-build` task already sets
# `runFromWorkspaceRoot: true`.

# ---------- Stage 1: builder ----------
FROM rust:1.90-slim-bookworm AS builder

WORKDIR /workspace

# Build deps (most are cheap; pkg-config + libssl-dev cover transitive
# crates that need them, even though sqlx uses rustls).
RUN apt-get update && apt-get install -y --no-install-recommends \
        pkg-config \
        libssl-dev \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Workspace manifest + lockfile + toolchain pin.
COPY Cargo.toml Cargo.lock rust-toolchain.toml ./

# templates-cli is a workspace member; provide a stub so the resolver
# is happy without pulling it into the build.
COPY apps/templates-cli/Cargo.toml apps/templates-cli/Cargo.toml
RUN mkdir -p apps/templates-cli/src \
    && echo "fn main() {}" > apps/templates-cli/src/main.rs \
    && echo "" > apps/templates-cli/src/lib.rs

# Cache rust-ai dependencies in a layer separate from source so
# rebuilds skip the long initial compile when only application code
# changes.
COPY apps/rust-ai/Cargo.toml apps/rust-ai/Cargo.toml
RUN mkdir -p apps/rust-ai/src \
    && echo "fn main() {}" > apps/rust-ai/src/main.rs \
    && echo "" > apps/rust-ai/src/lib.rs \
    && cargo build --release -p rust-ai 2>/dev/null || true

# Real source.
COPY apps/rust-ai/src apps/rust-ai/src
COPY apps/rust-ai/migrations apps/rust-ai/migrations
RUN touch apps/rust-ai/src/main.rs apps/rust-ai/src/lib.rs \
    && cargo build --release -p rust-ai

# ---------- Stage 2: runtime ----------
FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        libssl3 \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1000 fastapi \
    && useradd --system --uid 1000 --gid fastapi --create-home --shell /usr/sbin/nologin fastapi

WORKDIR /opt/rust-ai

COPY --from=builder /workspace/target/release/rust-ai /usr/local/bin/rust-ai
COPY --from=builder /workspace/apps/rust-ai/migrations /opt/rust-ai/migrations

USER fastapi

ARG HOST=0.0.0.0
ARG PORT=8080
ENV HOST=${HOST}
ENV PORT=${PORT}

EXPOSE ${PORT}

CMD ["/usr/local/bin/rust-ai"]
