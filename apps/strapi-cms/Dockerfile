# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_ENV=production
ARG NODE_VERSION=20

FROM busybox:1.37-glibc as glibc

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
RUN moon docker scaffold strapi-cms

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder

# Copy workspace skeleton
COPY --from=skeleton /root/.proto /root/.proto
COPY --from=skeleton /srv/.moon/docker/workspace .

# Install toolchain and dependencies
RUN moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .

# Build the application and prune the workspace
# RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run strapi-cms:build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker prune

# -----------------------------------------------------------------------------
# Rebuild Strapi CMS with production dependencies.
# -----------------------------------------------------------------------------
FROM base AS installer

# Only copy the require files for Strapi from builder stage
COPY --from=builder /srv/apps/strapi-cms /srv
COPY --from=builder /srv/.npmrc /srv/.npmrc

# Fix sharp dependencies then build the Strapi application.
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --ignore-scripts
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm rebuild --verbose better-sqlite3 sharp
RUN --mount=type=cache,id=pnpm,target=/pnpm/store NODE_ENV=$NODE_ENV pnpm build

# -----------------------------------------------------------------------------
# Cleanup the builder stage and create data directory.
# -----------------------------------------------------------------------------
FROM base AS pruner

# Copy the required files from the installer stage
COPY --from=installer /srv/package.json /srv/package.json
COPY --from=installer /srv/pnpm-lock.yaml /srv/pnpm-lock.yaml
COPY --from=installer /srv/node_modules /srv/node_modules
COPY --from=installer /srv/server.cjs /srv/server.cjs
COPY --from=installer /srv/dist /srv/dist

# Copy migrations and the generated files from Strapi plugins
COPY --from=installer /srv/src/extensions /srv/src/extensions
COPY --from=installer /srv/database /srv/database

# Cleanup uncessary files and install production dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod \
    --frozen-lockfile --ignore-scripts && pnpm prune --prod \
    --ignore-scripts && pnpm dlx clean-modules clean --yes \
    "!**/sharp/**" "!**/better-sqlite3/**" "**/codecov*"

# Set permissions for server.cjs and uploads directories
RUN mkdir -p /srv/public/uploads && chmod -R 0775 /srv/public
RUN mkdir -p /srv/_data && chmod -R 0775 /srv/_data
RUN chmod +x /srv/server.cjs

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/nodejs${NODE_VERSION}-debian12

# ------ Read application environment variables --------------------------------

ARG APP_KEYS API_TOKEN_SALT ADMIN_JWT_SECRET TRANSFER_TOKEN_SALT JWT_SECRET \
    DATABASE_CLIENT DATABASE_URL DATABASE_FILENAME DATABASE_SSL SMTP_HOST \
    SMTP_PORT SMTP_USERNAME SMTP_PASSWORD

# ------ Read application environment variables --------------------------------

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
COPY --from=glibc /bin/head /bin/head
COPY --from=glibc /bin/hostname /bin/hostname

# Define the host and port to listen on.
ARG NODE_ENV=production HOST=0.0.0.0 PORT=1337
ENV NODE_ENV=$NODE_ENV HOST=$HOST PORT=$PORT
ENV TINI_SUBREAPER=true

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/nodejs/bin/node", "server.cjs"]
