# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG MOON_VERSION=2.5.3
ARG NODE_ENV=production
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
COPY . .
RUN MOON_LOG=error moon docker scaffold {{ package_name | kebab_case }}

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder

# Copy workspace skeleton
COPY --from=skeleton /srv/.moon/docker/configs .

# Install toolchain and dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup

# Copy source files
COPY --from=skeleton /srv/.moon/docker/sources .

# Build the application and prune the workspace (optional).
# RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run {{ package_name | kebab_case }}:build
# RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker prune

# -----------------------------------------------------------------------------
# Rebuild Strapi CMS with production dependencies.
# -----------------------------------------------------------------------------
FROM base AS installer

# Only copy the require files from builder stage
COPY --from=builder /srv/apps/{{ package_name | kebab_case }} /srv

# Fix sharp dependencies then build the Strapi application.
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --ignore-scripts
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm rebuild --verbose sharp
RUN --mount=type=cache,id=pnpm,target=/pnpm/store NODE_ENV=$NODE_ENV pnpm build

# -----------------------------------------------------------------------------
# Cleanup the builder stage and create data directory.
# -----------------------------------------------------------------------------
FROM base AS pruner

# Copy the required files from the installer stage
COPY --from=installer /srv/package.json /srv/package.json
COPY --from=installer /srv/pnpm-lock.yaml /srv/pnpm-lock.yaml
COPY --from=installer /srv/node_modules /srv/node_modules
COPY --from=installer /srv/favicon.png /srv/favicon.png
COPY --from=installer /srv/server.cjs /srv/server.cjs
COPY --from=installer /srv/public /srv/public
COPY --from=installer /srv/dist /srv/dist

# Copy migrations and the generated files from Strapi plugins
COPY --from=installer /srv/src/extensions /srv/src/extensions
COPY --from=installer /srv/database /srv/database

# Cleanup uncessary files and install production dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod \
    --frozen-lockfile --ignore-scripts && pnpm approve-builds && pnpm \
    prune --prod --ignore-scripts && pnpm dlx clean-modules clean --yes \
    "!**/sharp/**" "**/codecov*"

# Set permissions for server.cjs and uploads directories
RUN mkdir -p /srv/public/uploads && chmod -R 0775 /srv/public
RUN mkdir -p /srv/_data && chmod -R 0775 /srv/_data
RUN chmod +x /srv/server.cjs

# -----------------------------------------------------------------------------
# Production image, copy build output files and run the application (runner).
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/nodejs${NODE_VERSION}-debian13:${DISTROLESS_TAG} AS runner

# Read application environment variables
ARG APP_KEYS
ARG API_TOKEN_SALT
ARG ADMIN_JWT_SECRET
ARG TRANSFER_TOKEN_SALT
ARG JWT_SECRET
ARG DATABASE_URL
ARG DATABASE_SSL
ARG SMTP_HOST
ARG SMTP_PORT
ARG SMTP_USERNAME
ARG SMTP_PASSWORD
ARG AWS_ACCESS_KEY_ID
ARG AWS_ACCESS_SECRET
ARG AWS_REGION
ARG AWS_BUCKET
ARG S3_PATH_PREFIX

# Copy build output files and required binaries from previous stage
COPY --from=pruner --chown=nonroot:nonroot /srv /srv
COPY --from=base /usr/bin/tini /usr/bin/tini

# Define the host and port to listen on.
ARG NODE_ENV=production HOST=0.0.0.0 PORT={{ port_number }}
ENV NODE_ENV=$NODE_ENV HOST=$HOST PORT=$PORT
ENV STRAPI_TELEMETRY_DISABLED=false
ENV TINI_SUBREAPER=true

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/nodejs/bin/node", "server.cjs"]
