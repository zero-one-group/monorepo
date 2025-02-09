# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=20

FROM nginx:1.27-alpine-slim as base_nginx
RUN apk update && apk add --no-cache tini

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
RUN moon docker scaffold react-app

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
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run react-app:build

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} nginx:1.27-alpine-slim as runner

# ----- Read application environment variables --------------------------------

ARG VITE_API_URL VITE_BASE_URL

# ----- Read application environment variables --------------------------------

# Copy the build output files from the installer stage.
COPY --from=builder --chown=nonroot:nonroot /srv/apps/react-app/build/client /srv

# Create the user and group that will run the application.
RUN addgroup --system --gid 1001 nonroot && adduser --system --uid 1001 nonroot
COPY --from=skeleton /srv/docker/nginx-static.conf /etc/nginx/nginx.conf
COPY --from=base_nginx /sbin/tini /sbin/tini

# Define the host and port to listen on.
ARG NODE_ENV=production NGINX_HOST=0.0.0.0 NGINX_PORT=3000
ENV NGINX_HOST=$NGINX_HOST NGINX_PORT=$NGINX_PORT
ENV NODE_ENV=$NODE_ENV TINI_SUBREAPER=true
STOPSIGNAL SIGQUIT

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $NGINX_PORT

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/sbin/nginx"]
