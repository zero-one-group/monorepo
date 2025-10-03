# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG DISTROLESS_TAG=nonroot
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=22
ARG GO_VERSION=1.25

# -----------------------------------------------------------------------------
# Base image for building the Golang app.
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} golang:${GO_VERSION}-trixie AS base_go
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0 COREPACK_INTEGRITY_KEYS=0
ENV FNM_PATH=/usr/bin PNPM_HOME=/pnpm HOME=/root
ENV PATH="$HOME/.local/share/fnm:$PNPM_HOME:$PATH"
ENV PATH="/root/.moon/bin:$PATH"
ARG NODE_VERSION=22.18

# Install system dependencies
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean && apt-get update -y \
    && apt-get -yqq --no-install-recommends install build-essential curl \
       inotify-tools pkg-config libssl-dev git unzip ca-certificates jq \
    && update-ca-certificates

# Install Node.js and PNPM
# We need to use bash as the default shell to load fnm function.
SHELL ["/usr/bin/bash", "-c"]
RUN echo -e "Installing Node.js $NODE_VERSION" && touch $HOME/.bashrc \
    && curl -fsSL https://fnm.vercel.app/install | bash -s -- --skip-shell --install-dir $FNM_PATH \
    && echo 'eval "$(fnm env --use-on-cd --shell bash)"' >> $HOME/.bashrc && . $HOME/.bashrc \
    && fnm install ${NODE_VERSION} --corepack-enabled \
    && fnm default ${NODE_VERSION} --corepack-enabled \
    && corepack prepare pnpm@latest-10 --activate

# Install required SDK
RUN --mount=type=cache,id=go-sdk,target=/root/.cache/go-sdk \
    go install github.com/swaggo/swag/cmd/swag@latest

# Add tini for signal handling and zombie reaping
RUN set -eux; \
    TINI_DOWNLOAD_URL="https://github.com/krallin/tini/releases/download/v0.19.0" \
    ARCH="$(dpkg --print-architecture)"; \
    case "${ARCH}" in \
      amd64|x86_64) TINI_BIN_URL="${TINI_DOWNLOAD_URL}/tini" ;; \
      arm64|aarch64) TINI_BIN_URL="${TINI_DOWNLOAD_URL}/tini-arm64" ;; \
      *) echo "unsupported architecture: ${ARCH}"; exit 1 ;; \
    esac; \
    curl -fsSL "${TINI_BIN_URL}" -o /usr/bin/tini; \
    chmod +x /usr/bin/tini

# Install Moonrepo CLI
RUN curl -fsSL https://moonrepo.dev/install/moon.sh | bash

WORKDIR /srv

# -----------------------------------------------------------------------------
# Scaffold the specific project.
# -----------------------------------------------------------------------------
FROM base_go AS skeleton
COPY --link . .

# Bash will load fnm function via .bashrc
# since we're starting non-interactive shell,
# we wil need to tell bash to load .bashrc manually
SHELL ["/usr/bin/bash", "-c"]
ENV BASH_ENV=$HOME/.bashrc

RUN moon docker scaffold {{ package_name | kebab_case }}

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base_go AS builder
ENV CGO_ENABLED=1

# Enable caching for Go modules
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

# Copy workspace skeleton
COPY --from=skeleton /srv/.moon/docker/workspace .

# Install toolchain and dependencies
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache moon setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .

# Build the application and prune the workspace
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache moon run {{ package_name | kebab_case }}:tidy
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache moon run {{ package_name | kebab_case }}:build
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache moon docker prune

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/cc-debian12:${DISTROLESS_TAG} AS runner

# ----- Read application environment variables --------------------------------

ARG APP_ENVIRONMENT
ARG APP_HOST
ARG APP_PORT
ARG CORS_ALLOW_ORIGINS
ARG DATABASE_URL
ARG ENABLE_INSTRUMENTATION
ARG ENABLE_SWAGGER
ARG JWT_SECRET
ARG JWT_TOKEN_EXPIRY_MINUTES
ARG LOG_LEVEL
ARG OTEL_EXPORTER_OTLP_ENDPOINT
ARG SERVICE_NAME
ARG TRACING_SAMPLE_RATE

# ----- Read application environment variables --------------------------------

# Copy the build output files and necessary utilities from previous stage.
COPY --from=builder --chown=nonroot:nonroot /srv/apps/{{ package_name | kebab_case }}/build/{{ package_name | kebab_case }} /srv/{{ package_name | kebab_case }}
COPY --from=builder --chmod=775 /usr/bin/tini /usr/bin/tini

# Define the host and port to listen on.
ARG HOST=0.0.0.0 PORT={{ port_number }}
ENV HOST=$HOST PORT=$PORT
ENV TINI_SUBREAPER=true

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $PORT/tcp

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/srv/{{ package_name | kebab_case }}"]
