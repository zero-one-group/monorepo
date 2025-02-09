# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=20

FROM busybox:1.37-glibc as glibc
FROM joseluisq/static-web-server:2 as httpd

# -----------------------------------------------------------------------------
# Base image with pnpm package manager.
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} node:${NODE_VERSION}-bookworm-slim AS base
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0 COREPACK_INTEGRITY_KEYS=0
ENV LEFTHOOK=0 CI=true PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=true
ENV MOON_TOOLCHAIN_FORCE_GLOBALS=1 MOON_INSTALL_DIR="/usr/bin"
ENV PNPM_HOME="/pnpm" PATH="$PNPM_HOME:$MOON_INSTALL_DIR:$PATH"
RUN corepack enable && corepack prepare pnpm@latest-10 --activate
WORKDIR /srv

# Install system dependencies and moon cli
RUN apt-get update && apt-get -yqq --no-install-recommends install curl tini jq ca-certificates
RUN apt -yqq purge && update-ca-certificates && apt -yqq autoremove && apt -yqq clean
RUN curl -fsSL https://moonrepo.dev/install/moon.sh | bash

# -----------------------------------------------------------------------------
# Scaffold the specific project.
# -----------------------------------------------------------------------------
FROM base AS skeleton
COPY --link . .
RUN moon docker scaffold {{ package_name | kebab_case }}

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder

# Copy workspace skeleton
COPY --from=skeleton /root/.proto /root/.proto
COPY --from=skeleton /srv/.moon/docker/workspace .

# Install toolchain and dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .

# Build the application and prune the workspace (optional).
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run {{ package_name | kebab_case }}:build

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/base-nossl-debian12 as runner

# ----- Read application environment variables --------------------------------

ARG VITE_API_URL VITE_BASE_URL

# ----- Read application environment variables --------------------------------

# Copy the build output files from the installer stage.
COPY --from=builder --chown=nonroot:nonroot /srv/apps/{{ package_name | kebab_case }}/build/client /srv

# Copy the required binaries and configuration file.
COPY --from=skeleton /srv/docker/sws.toml /etc/sws.toml
COPY --from=httpd /static-web-server /usr/bin/sws
COPY --from=base /usr/bin/tini /usr/bin/tini

# Copy some necessary system utilities from previous stage (~7MB).
# To enhance security, consider avoiding the copying of sysutils.
# COPY --from=glibc /usr/bin/env /usr/bin/env
# COPY --from=glibc /bin/clear /bin/clear
# COPY --from=glibc /bin/mkdir /bin/mkdir
# COPY --from=glibc /bin/which /bin/which
# COPY --from=glibc /bin/cat /bin/cat
# COPY --from=glibc /bin/ls /bin/ls
# COPY --from=glibc /bin/sh /bin/sh

# Define the log level, host and port to listen on.
ARG SERVER_LOG_LEVEL=info SERVER_HOST=0.0.0.0 SERVER_PORT={{ port_number }}
ENV SERVER_CONFIG_FILE=/etc/sws.toml SERVER_LOG_LEVEL=$SERVER_LOG_LEVEL
ENV SERVER_HOST=$SERVER_HOST SERVER_PORT=$SERVER_PORT TINI_SUBREAPER=true
STOPSIGNAL SIGQUIT

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $SERVER_PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/usr/bin/sws"]
