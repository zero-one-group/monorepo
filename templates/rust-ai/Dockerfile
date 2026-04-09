# Multi-stage Rust build for the {{ package_name | kebab_case }} service.
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

# Cache {{ package_name | kebab_case }} dependencies in a layer separate from source so
# rebuilds skip the long initial compile when only application code
# changes.
COPY apps/{{ package_name | kebab_case }}/Cargo.toml apps/{{ package_name | kebab_case }}/Cargo.toml
RUN mkdir -p apps/{{ package_name | kebab_case }}/src \
    && echo "fn main() {}" > apps/{{ package_name | kebab_case }}/src/main.rs \
    && echo "" > apps/{{ package_name | kebab_case }}/src/lib.rs \
    && cargo build --release -p {{ package_name | kebab_case }} 2>/dev/null || true

# Real source.
COPY apps/{{ package_name | kebab_case }}/src apps/{{ package_name | kebab_case }}/src
COPY apps/{{ package_name | kebab_case }}/migrations apps/{{ package_name | kebab_case }}/migrations
RUN touch apps/{{ package_name | kebab_case }}/src/main.rs apps/{{ package_name | kebab_case }}/src/lib.rs \
    && cargo build --release -p {{ package_name | kebab_case }}

# ---------- Stage 2: runtime ----------
FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        libssl3 \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1000 fastapi \
    && useradd --system --uid 1000 --gid fastapi --create-home --shell /usr/sbin/nologin fastapi

WORKDIR /opt/{{ package_name | kebab_case }}

COPY --from=builder /workspace/target/release/{{ package_name | kebab_case }} /usr/local/bin/{{ package_name | kebab_case }}
COPY --from=builder /workspace/apps/{{ package_name | kebab_case }}/migrations /opt/{{ package_name | kebab_case }}/migrations

USER fastapi

ARG HOST=0.0.0.0
ARG PORT={{ port_number }}
ENV HOST=${HOST}
ENV PORT=${PORT}

EXPOSE ${PORT}

CMD ["/usr/local/bin/{{ package_name | kebab_case }}"]
