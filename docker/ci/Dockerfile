# syntax=docker/dockerfile:1.7
# CI image for moon + pnpm monorepos (Go + Node). Public: docker.io/01group/moon-ci
#
# Contains exactly the toolchain this monorepo pins in .moon/toolchains.yml (and
# .moon/workspace.yml for moon), so a pipeline job skips apt + proto + tool downloads
# (~3–4 min/job). Nothing project-specific is baked in: sources, node_modules and the Go
# module cache come from the checkout and the CI cache.
#
# Size matters more than anything else here: runners cold-pull it once per tag. Node comes
# from the base image (not a second copy via proto) and the Go CLI tools are built in a
# throwaway stage so only their binaries land in the final image.
#
# Built and pushed by .github/workflows/ci-image.yaml on every push to main that touches
# docker/ci/**; the tag is derived from the pins so it cannot disagree with them:
#   moon<M>-go<G.major.minor>-node<N.major>-pnpm<P.major>   e.g. moon2.5.3-go1.27-node24-pnpm11
# Manual build (amd64):
#   docker buildx build --platform linux/amd64 -f docker/ci/Dockerfile \
#     -t 01group/moon-ci:moon2.5.3-go1.27-node24-pnpm11 -t 01group/moon-ci:latest --push docker/ci
# Bump a version → bump the pin, the workflow rebuilds and pushes the new tag; pin it in CI.

ARG NODE_VERSION=24
ARG GO_VERSION=1.27.0

# ---- stage 1: build Go CLI tools, keep only the binaries -------------------------------
FROM golang:${GO_VERSION}-trixie AS gotools
ARG SWAG_VERSION=v1.16.6
ARG GOTESTSUM_VERSION=v1.13.0
ARG MOCKERY_VERSION=v3.5.1
ENV GOFLAGS=-trimpath CGO_ENABLED=0 GOBIN=/out
RUN go install "github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}" \
 && go install "gotest.tools/gotestsum@${GOTESTSUM_VERSION}" \
 && go install "github.com/vektra/mockery/v3@${MOCKERY_VERSION}"

# ---- stage 2: the CI image ---------------------------------------------------------------
FROM node:${NODE_VERSION}-trixie-slim
ARG GO_VERSION
ARG PROTO_VERSION=0.61.1
ARG MOON_VERSION=2.5.3
ARG PNPM_VERSION=11.23.0

# proto installs moon and Go under PROTO_HOME (outside $HOME so any uid works; runners use root).
ENV PROTO_HOME=/opt/proto \
    GOPATH=/opt/go \
    GOTOOLCHAIN=local \
    GOFLAGS=-buildvcs=false \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0 \
    MOON_ALLOW_ROOT=1 \
    CI=true
ENV PATH="${PROTO_HOME}/shims:${PROTO_HOME}/bin:${GOPATH}/bin:/usr/local/bin:${PATH}"

RUN apt-get update -qq \
 && apt-get install -y -qq --no-install-recommends ca-certificates curl git jq unzip xz-utils openssh-client \
 && rm -rf /var/lib/apt/lists/*

# pnpm via corepack on the base image's Node — no second Node install.
RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate && pnpm --version

# proto → moon + Go only. `proto use` in a project still works: node/pnpm resolve to these.
RUN curl -fsSL https://moonrepo.dev/install/proto.sh | bash -s -- "${PROTO_VERSION}" --yes --no-profile \
 && proto install moon "${MOON_VERSION}" --pin=global \
 && proto install go "${GO_VERSION}" --pin=global \
 && rm -rf "${PROTO_HOME}/temp" "${PROTO_HOME}/tools/go/${GO_VERSION}/test" "${PROTO_HOME}/tools/go/${GO_VERSION}/doc" \
 && moon --version && go version

COPY --from=gotools /out/swag /out/gotestsum /out/mockery /usr/local/bin/
RUN swag --version && gotestsum --version && mockery version

# Login shells (`bash -l`, interactive `docker run`) source /etc/profile, which resets PATH.
RUN printf 'export PROTO_HOME=%s GOPATH=%s\nexport PATH="%s/shims:%s/bin:%s/bin:$PATH"\n' \
      "${PROTO_HOME}" "${GOPATH}" "${PROTO_HOME}" "${PROTO_HOME}" "${GOPATH}" > /etc/profile.d/moon-ci.sh

WORKDIR /workspace
CMD ["bash"]
