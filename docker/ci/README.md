# `01group/moon-ci` — CI image for moon + pnpm monorepos

Public image on Docker Hub with the toolchain this monorepo pins: proto → moon, Go, Node,
pnpm, plus `swag`, `gotestsum`, `mockery`. Pipelines skip the 3–4 minute apt/proto/download
dance per job. Nothing project-specific is baked in.

| Tag                              | Meaning                                                            |
| -------------------------------- | ------------------------------------------------------------------ |
| `moon2.5.3-go1.27-node24-pnpm11` | immutable — what a project's CI pins                               |
| `latest`                         | newest build from `main`                                           |

The immutable tag is **derived from the pins** (`.moon/workspace.yml` `versionConstraint`,
`.moon/toolchains.yml` `go` / `node` / `pnpm`) by `.github/workflows/ci-image.yaml`, so the tag
can never disagree with what the repo actually uses.

## Rebuild

Automatic: every push to `main` that touches `docker/ci/**`, `.moon/toolchains.yml` or
`.moon/workspace.yml` builds `linux/amd64` and pushes the derived tag + `latest`. Requires the
`DOCKERHUB_USERNAME` / `DOCKERHUB_TOKEN` organisation secrets with push rights to `01group`.

Manual:

```sh
docker login                        # as 01group
docker buildx build --platform linux/amd64 -f docker/ci/Dockerfile \
  --build-arg GO_VERSION=1.27.0 --build-arg MOON_VERSION=2.5.3 --build-arg PNPM_VERSION=11.23.0 \
  -t 01group/moon-ci:moon2.5.3-go1.27-node24-pnpm11 -t 01group/moon-ci:latest \
  --push docker/ci
```

Verify a tag before pinning it:

```sh
docker run --rm -it 01group/moon-ci:<tag> bash -lc 'moon --version; go version; node --version; pnpm --version; swag --version'
```

## Using it

`templates/gitlab-cicd/dependencies/base.yml` sets `.setup.image` to the immutable tag. When you
bump a toolchain pin, bump that tag in the same change; the `.gitlab-ci.yml` of every generated
project inherits it on the next `moon generate`.
