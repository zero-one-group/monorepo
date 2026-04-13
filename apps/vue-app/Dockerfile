# syntax=docker/dockerfile:1.7

ARG MOON_VERSION=2.1.3
ARG NODE_VERSION=22

FROM joseluisq/static-web-server:2 AS httpd

# -----------------------------------------------------------------------------
# Base image with pnpm package manager.
# -----------------------------------------------------------------------------
FROM node:${NODE_VERSION}-trixie-slim AS base
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0 COREPACK_INTEGRITY_KEYS=0 DO_NOT_TRACK=1
ENV NODE_VERSION=$NODE_VERSION PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=true PNPM_HOME="/pnpm"
ENV CI=true MOON_TOOLCHAIN_FORCE_GLOBALS=1 MOON_INSTALL_DIR="/usr/bin"
ENV PATH="$PNPM_HOME:$MOON_INSTALL_DIR:$PATH"

RUN corepack enable && corepack prepare pnpm@latest-10 --activate

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean && apt-get update -y \
    && apt-get -yqq --no-install-recommends install build-essential curl \
       inotify-tools pkg-config libssl-dev git unzip ca-certificates jq \
    && update-ca-certificates

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

RUN bash -c 'bash <(curl -fsSL https://moonrepo.dev/install/moon.sh) "$MOON_VERSION"'

WORKDIR /srv

# -----------------------------------------------------------------------------
# Scaffold the specific project.
# -----------------------------------------------------------------------------
FROM base AS skeleton
COPY --link . .
RUN MOON_LOG=error moon docker scaffold vue-app

# -----------------------------------------------------------------------------
# Install dependencies and build the application.
# -----------------------------------------------------------------------------
FROM base AS builder

COPY --from=skeleton /srv/.moon/docker/configs .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon docker setup
COPY --from=skeleton /srv/.moon/docker/sources .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store moon run vue-app:build

# -----------------------------------------------------------------------------
# Production image (runner).
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/base-nossl-debian12 AS runner

ARG VITE_API_URL VITE_BASE_URL

COPY --from=builder --chown=nonroot:nonroot /srv/apps/vue-app/dist /srv
COPY --from=skeleton /srv/docker/sws-spa.toml /etc/sws.toml
COPY --from=httpd /static-web-server /usr/bin/sws
COPY --from=base /usr/bin/tini /usr/bin/tini

ARG SERVER_LOG_LEVEL=info SERVER_HOST=0.0.0.0 SERVER_PORT=3000
ENV SERVER_CONFIG_FILE=/etc/sws.toml SERVER_LOG_LEVEL=$SERVER_LOG_LEVEL
ENV SERVER_HOST=$SERVER_HOST SERVER_PORT=$SERVER_PORT TINI_SUBREAPER=true
STOPSIGNAL SIGQUIT

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $SERVER_PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/usr/bin/sws"]
