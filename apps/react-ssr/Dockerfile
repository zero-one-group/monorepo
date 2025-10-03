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
RUN moon docker scaffold react-ssr

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder

# Copy workspace skeleton
# COPY --from=skeleton /root/.proto /root/.proto
COPY --from=skeleton /srv/.moon/docker/workspace .

# Install toolchain and dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .
COPY --from=skeleton /srv/.npmrc .

# Build the application and prune the workspace (optional).
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run react-ssr:build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker prune

# -----------------------------------------------------------------------------
# Cleanup the builder stage and create data directory.
# TODO: reduce image size, cleanup unnecessary dependencies from other packages
# -----------------------------------------------------------------------------
FROM base AS pruner

# Copy global dependencies from builder stage
COPY --from=builder /srv/node_modules /srv/node_modules
COPY --from=builder /srv/package.json /srv/package.json
COPY --from=builder /srv/pnpm-lock.yaml /srv/pnpm-lock.yaml

# Only copy the require files from builder stage
COPY --from=builder /srv/apps/react-ssr/node_modules /srv/apps/react-ssr/node_modules
COPY --from=builder /srv/apps/react-ssr/package.json /srv/apps/react-ssr/package.json
COPY --from=builder /srv/apps/react-ssr/server.mjs /srv/apps/react-ssr/server.mjs
COPY --from=builder /srv/apps/react-ssr/build /srv/apps/react-ssr/build

# Cleanup uncessary files and set executable permissions for the server.
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod \
    --frozen-lockfile --ignore-scripts && pnpm prune --prod \
    --ignore-scripts && pnpm dlx clean-modules clean --yes \
    "!**/sharp/**" "!**/@repo/**"

RUN chmod +x /srv/apps/react-ssr/server.mjs

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/nodejs${NODE_VERSION}-debian12

# ----- Read application environment variables --------------------------------

ARG VITE_API_URL VITE_BASE_URL

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
ARG NODE_ENV=production HOST=0.0.0.0 PORT=3100
ENV NODE_ENV=$NODE_ENV HOST=$HOST PORT=$PORT
ENV TINI_SUBREAPER=true

WORKDIR /srv/apps/react-ssr
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/nodejs/bin/node", "server.mjs"]
