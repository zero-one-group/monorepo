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
RUN apt-get update && apt-get -yqq --no-install-recommends install curl tini jq ca-certificates git
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

ARG APP_BASE_URL
ARG CORS_CREDENTIALS
ARG CORS_MAX_AGE
ARG CORS_ORIGINS
ARG ENABLE_API_DOCS
ARG JWT_ALGORITHM
ARG JWT_SECRET_KEY
ARG RATE_LIMIT_BURST_SIZE
ARG RATE_LIMIT_ENABLED
ARG RATE_LIMIT_REQUESTS
ARG DATABASE_URL
ARG PG_MAX_POOL_SIZE
ARG PG_MAX_RETRIES
ARG SMTP_HOST
ARG SMTP_PASSWORD
ARG SMTP_PORT
ARG SMTP_SECURE
ARG SMTP_SENDER_EMAIL
ARG SMTP_SENDER_NAME
ARG SMTP_USERNAME
ARG PUBLIC_ASSETS_URL
ARG S3_ACCESS_KEY
ARG S3_BUCKET_NAME
ARG S3_ENDPOINT
ARG S3_FORCE_PATH_STYLE
ARG S3_REGION
ARG S3_SECRET_KEY
ARG S3_USE_SSL
ARG LOG_FORMAT
ARG LOG_LEVEL
ARG LOG_NO_COLOR
ARG OTEL_ENABLE_TELEMETRY
ARG OTEL_EXPORTER_OTLP_ENDPOINT
ARG OTEL_EXPORTER_OTLP_HEADERS
ARG OTEL_EXPORTER_OTLP_PROTOCOL
ARG OTEL_INSECURE_MODE
ARG OTEL_SERVICE_NAME
ARG OTEL_TRACING_SAMPLE_RATE

# ----- Read application environment variables --------------------------------

# Copy the build output files from the installer stage.
COPY --from=builder --chown=nonroot:nonroot /srv/apps/{{ package_name | kebab_case }}/build/{{ package_name | kebab_case }} /srv/{{ package_name | kebab_case }}

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
ARG APP_MODE=production
ARG SERVER_HOST=0.0.0.0
ARG SERVER_PORT={{ port_number }}
ENV APP_MODE=$APP_MODE
ENV SERVER_HOST=$SERVER_HOST
ENV SERVER_PORT=$SERVER_PORT
ENV TINI_SUBREAPER=true

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $SERVER_PORT/tcp

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/srv/{{ package_name | kebab_case }}", "serve"]
