---
title: Project Structure
slug: "project-structure"
---

The monorepo follows a standard structure designed to keep code organized and maintainable.

## Directory Layout

Here is an overview of the top-level directories:

- **`/apps`**: Contains the source code for all applications (web, mobile, backend services). Each project here is a separate deployable unit.
- **`/packages`**: Contains shared libraries and packages that are used across multiple applications. This promotes code reuse and consistency.
- **`/templates`**: Stores the source files for the project templates used by `moon generate`.
- **`/docsite`**: Contains the source code and content for this documentation site.
- **`/.moon`**: Holds the configuration for the moonrepo build system and toolchain.
- **`/builder`**: Contains scripts and utilities used for scaffolding and building templates.
- **`/docker`**: Docker related configuration and files.

## Key Files

- **`package.json`**: The root configuration for Node.js dependencies and scripts.
- **`pnpm-workspace.yaml`**: Defines the workspace topology for pnpm.
- **`moon.yml`**: (In project roots) configuration file for moonrepo tasks and project settings.
