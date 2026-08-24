# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG MOON_VERSION=2.5.3
ARG NODE_VERSION=24
ARG DISTROLESS_TAG=nonroot

# -----------------------------------------------------------------------------
# Base image with pnpm package manager.
# -----------------------------------------------------------------------------
FROM node:${NODE_VERSION}-trixie-slim AS base
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0 COREPACK_INTEGRITY_KEYS=0 DO_NOT_TRACK=1
ENV NODE_VERSION=$NODE_VERSION PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=true PNPM_HOME="/pnpm"
ENV CI=true MOON_TOOLCHAIN_FORCE_GLOBALS=1 MOON_INSTALL_DIR="/usr/bin"
ENV PATH="$PNPM_HOME:$MOON_INSTALL_DIR:$PATH"

# Install and enable PNPM via corepack
RUN corepack enable && corepack prepare pnpm@latest-10 --activate

# Install system dependencies (optional tools for building and debugging)
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean && apt-get update -y \
    && apt-get -yqq --no-install-recommends install build-essential curl \
       inotify-tools pkg-config libssl-dev git unzip ca-certificates jq \
    && update-ca-certificates

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

# Install Moon CLI
RUN bash -c 'bash <(curl -fsSL https://moonrepo.dev/install/moon.sh) "$MOON_VERSION"'

WORKDIR /srv

# -----------------------------------------------------------------------------
# Scaffold the specific project.
# -----------------------------------------------------------------------------
FROM base AS skeleton
COPY --link . .
RUN MOON_LOG=error moon docker scaffold nextjs-app

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder
ARG NODE_ENV=production
ENV NODE_ENV=$NODE_ENV NEXT_TELEMETRY_DISABLED=1

# Copy workspace skeleton
COPY --from=skeleton /srv/.moon/docker/configs .

# Install toolchain and dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .
COPY --from=skeleton /srv/.npmrc .

# Build the application and prune the workspace (optional).
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run nextjs-app:build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker prune

# -----------------------------------------------------------------------------
# Cleanup the builder stage and create data directory.
# -----------------------------------------------------------------------------
FROM base AS pruner

# Copy only necessary files from builder stage
COPY --from=builder /srv/apps/nextjs-app/.next/standalone /srv
COPY --from=builder /srv/apps/nextjs-app/.next/static /srv/apps/nextjs-app/.next/static
COPY --from=builder /srv/apps/nextjs-app/public /srv/apps/nextjs-app/public

# Copy Next.js configuration files.
COPY --from=skeleton /srv/apps/nextjs-app/next.config.mjs /srv/apps/nextjs-app/next.config.mjs

# Set permissions for the public directory and the server.js file.
RUN chmod -R 0775 /srv/apps/nextjs-app/public && chmod +x /srv/apps/nextjs-app/server.js

# Cleanup uncessary files and install production dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm dlx \
    clean-modules clean --directory "apps/node_modules" --yes \
    "!**/next/**" "!**/sharp/**" "**/codecov*" "!**/@repo/**"

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/nodejs${NODE_VERSION}-debian13:${DISTROLESS_TAG} AS runner

# Read application environment variables
ARG DATABASE_URL
ARG SMTP_HOST
ARG SMTP_PORT
ARG SMTP_USERNAME
ARG SMTP_PASSWORD
ARG SMTP_EMAIL_FROM

# Copy build output files and required binaries from previous stage
COPY --from=pruner --chown=nonroot:nonroot /srv /srv
COPY --from=base /usr/bin/tini /usr/bin/tini

# Define the host and port to listen on.
ARG NODE_ENV=production HOSTNAME=0.0.0.0 PORT=3200
ENV NODE_ENV=$NODE_ENV HOSTNAME=$HOSTNAME PORT=$PORT
ENV TINI_SUBREAPER=true NEXT_TELEMETRY_DISABLED=1

WORKDIR /srv/apps/nextjs-app
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/nodejs/bin/node", "server.js"]
