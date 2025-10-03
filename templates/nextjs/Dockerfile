# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=22

FROM busybox:1.37-glibc AS glibc

# -----------------------------------------------------------------------------
# Base image with pnpm package manager.
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} node:${NODE_VERSION}-bookworm-slim AS base
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0 COREPACK_INTEGRITY_KEYS=0
ENV LEFTHOOK=0 CI=true PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=true
ENV MOON_TOOLCHAIN_FORCE_GLOBALS=1 MOON_INSTALL_DIR="/usr/bin" PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$MOON_INSTALL_DIR:$PATH"
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
ARG NODE_ENV=production
ENV NODE_ENV=$NODE_ENV NEXT_TELEMETRY_DISABLED=1

# Copy workspace skeleton
# COPY --from=skeleton /root/.proto /root/.proto
COPY --from=skeleton /srv/.moon/docker/workspace .

# Install toolchain and dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .
COPY --from=skeleton /srv/.npmrc .

# Build the application and prune the workspace (optional).
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run {{ package_name | kebab_case }}:build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker prune

# -----------------------------------------------------------------------------
# Cleanup the builder stage and create data directory.
# -----------------------------------------------------------------------------
FROM base AS pruner

# Copy only necessary files from builder stage
COPY --from=builder /srv/apps/{{ package_name | kebab_case }}/.next/standalone /srv
COPY --from=builder /srv/apps/{{ package_name | kebab_case }}/.next/static /srv/apps/{{ package_name | kebab_case }}/.next/static
COPY --from=builder /srv/apps/{{ package_name | kebab_case }}/public /srv/apps/{{ package_name | kebab_case }}/public

# Copy Next.js configuration files.
COPY --from=skeleton /srv/apps/{{ package_name | kebab_case }}/next.config.mjs /srv/apps/{{ package_name | kebab_case }}/next.config.mjs

# Set permissions for the public directory and the server.js file.
RUN chmod -R 0775 /srv/apps/{{ package_name | kebab_case }}/public && chmod +x /srv/apps/{{ package_name | kebab_case }}/server.js

# Cleanup uncessary files and install production dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm dlx \
    clean-modules clean --directory "apps/node_modules" --yes \
    "!**/next/**" "!**/sharp/**" "**/codecov*" "!**/@repo/**"

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/nodejs${NODE_VERSION}-debian12

# ----- Read application environment variables --------------------------------

ARG DATABASE_URL SMTP_HOST SMTP_PORT SMTP_USERNAME SMTP_PASSWORD SMTP_EMAIL_FROM

# ----- Read application environment variables --------------------------------

# Copy the build output files from the installer stage.
COPY --from=pruner --chown=nonroot:nonroot /srv /srv

# Copy some necessary system utilities from previous stage (~7MB).
# To enhance security, consider avoiding the copying of sysutils.
COPY --from=base /usr/bin/tini /usr/bin/tini
COPY --from=glibc /usr/bin/env /usr/bin/env
COPY --from=glibc /bin/clear /bin/clear
COPY --from=glibc /bin/mkdir /bin/mkdir
COPY --from=glibc /bin/which /bin/which
COPY --from=glibc /bin/cat /bin/cat
COPY --from=glibc /bin/ls /bin/ls
COPY --from=glibc /bin/sh /bin/sh

# Define the host and port to listen on.
ARG NODE_ENV=production HOSTNAME=0.0.0.0 PORT={{ port_number }}
ENV NODE_ENV=$NODE_ENV HOSTNAME=$HOSTNAME PORT=$PORT
ENV TINI_SUBREAPER=true NEXT_TELEMETRY_DISABLED=1

WORKDIR /srv/apps/{{ package_name | kebab_case }}
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/nodejs/bin/node", "server.js"]
