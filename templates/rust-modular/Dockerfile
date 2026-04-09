# Multi-stage Rust build for the {{ package_name | kebab_case }} service (Phase D).
#
# Stage 1: build the release binary inside a full Rust image.
# Stage 2: copy the binary into a minimal debian-slim runtime.
#
# Built from the monorepo root so the workspace Cargo.lock is visible
# to the build context. moon's `docker-build` task sets
# `runFromWorkspaceRoot: true`.

# ---------- Stage 1: builder ----------
FROM rust:1.90-slim-bookworm AS builder

WORKDIR /workspace

RUN apt-get update && apt-get install -y --no-install-recommends \
        pkg-config \
        libssl-dev \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock rust-toolchain.toml ./

# Workspace members not built in this image need stub Cargo.toml files
# so the resolver is happy without pulling them into the build.
COPY apps/templates-cli/Cargo.toml apps/templates-cli/Cargo.toml
RUN mkdir -p apps/templates-cli/src \
    && echo "fn main() {}" > apps/templates-cli/src/main.rs \
    && echo "" > apps/templates-cli/src/lib.rs

COPY apps/rust-ai/Cargo.toml apps/rust-ai/Cargo.toml
RUN mkdir -p apps/rust-ai/src \
    && echo "fn main() {}" > apps/rust-ai/src/main.rs \
    && echo "" > apps/rust-ai/src/lib.rs

COPY apps/rust-clean/Cargo.toml apps/rust-clean/Cargo.toml
RUN mkdir -p apps/rust-clean/src \
    && echo "fn main() {}" > apps/rust-clean/src/main.rs \
    && echo "" > apps/rust-clean/src/lib.rs

# Cache {{ package_name | kebab_case }} dependencies in a layer separate from source.
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
    && groupadd --system --gid 1000 gomodular \
    && useradd --system --uid 1000 --gid gomodular --create-home --shell /usr/sbin/nologin gomodular

WORKDIR /opt/{{ package_name | kebab_case }}

COPY --from=builder /workspace/target/release/{{ package_name | kebab_case }} /usr/local/bin/{{ package_name | kebab_case }}
COPY --from=builder /workspace/apps/{{ package_name | kebab_case }}/migrations /opt/{{ package_name | kebab_case }}/migrations

USER gomodular

ARG SERVER_HOST=0.0.0.0
ARG SERVER_PORT={{ port_number }}
ENV SERVER_HOST=${SERVER_HOST}
ENV SERVER_PORT=${SERVER_PORT}

EXPOSE ${SERVER_PORT}

ENTRYPOINT ["/usr/local/bin/{{ package_name | kebab_case }}"]
CMD ["serve"]
