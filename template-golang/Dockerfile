# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=22
ARG GO_VERSION=1.24

FROM busybox:1.37-glibc as glibc

# -----------------------------------------------------------------------------
# Base image for building the Golang app.
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} golang:${GO_VERSION}-bookworm AS base_go
ENV MOON_TOOLCHAIN_FORCE_GLOBALS=1 MOON_INSTALL_DIR="/usr/bin"
ENV PATH="$MOON_INSTALL_DIR:$PATH" LEFTHOOK=0 CI=true
WORKDIR /srv

# Install system dependencies and moon cli
RUN apt-get update && apt-get -yqq --no-install-recommends install curl tini jq ca-certificates
RUN apt -yqq purge && update-ca-certificates && apt -yqq autoremove && apt -yqq clean
RUN curl -fsSL https://moonrepo.dev/install/moon.sh | bash

# -----------------------------------------------------------------------------
# Scaffold the specific project.
# -----------------------------------------------------------------------------
FROM base_go AS skeleton
COPY --link . .
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
COPY --from=skeleton /root/.proto /root/.proto
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
FROM --platform=${PLATFORM} gcr.io/distroless/cc-debian12 AS runner

# ----- Read application environment variables --------------------------------

ARG  DATABASE_URL SMTP_HOST SMTP_PORT SMTP_USERNAME SMTP_PASSWORD SMTP_EMAIL_FROM

# ----- Read application environment variables --------------------------------

# Copy the build output files from the installer stage.
COPY --from=builder --chown=nonroot:nonroot /srv/apps/{{ package_name | kebab_case }}/build/release/{{ package_name | kebab_case }} /srv/{{ package_name | kebab_case }}

# Copy some necessary system utilities from previous stage (~7MB).
# To enhance security, consider avoiding the copying of sysutils.
COPY --from=base_go /usr/bin/tini /usr/bin/tini
COPY --from=glibc /usr/bin/env /usr/bin/env
COPY --from=glibc /bin/clear /bin/clear
COPY --from=glibc /bin/mkdir /bin/mkdir
COPY --from=glibc /bin/which /bin/which
COPY --from=glibc /bin/cat /bin/cat
COPY --from=glibc /bin/ls /bin/ls
COPY --from=glibc /bin/sh /bin/sh

# Define the host and port to listen on.
ARG SERVER_ENV=production HOST=0.0.0.0 PORT={{ port_number }}
ENV SERVER_ENV=$SERVER_ENV TINI_SUBREAPER=true
ENV HOST=$HOST PORT=$PORT

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/srv/{{ package_name | kebab_case }}"]
