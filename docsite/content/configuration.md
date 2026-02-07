---
title: Configuration
slug: "configuration"
---

Configuration in the monorepo is managed primarily through **moonrepo** and environment variables.

## Moonrepo Configuration

Moonrepo uses a set of YAML files to configure the workspace and individual projects.

### Workspace Configuration (`.moon/`)

Located in the root `/.moon` directory, these files control the global behavior:

- **`workspace.yml`**: Defines the projects, generator templates, and VCS settings.
- **`toolchain.yml`**: Configures the languages and tools available in the workspace (Node.js, Go, Python, etc.).
- **`tasks.yml`**: Defines global tasks that can be inherited by all projects.

### Project Configuration (`moon.yml`)

Each project in `/apps` or `/packages` has a `moon.yml` file. This file defines:

- **`language`**: The primary language of the project.
- **`type`**: Whether it is an `application` or a `library`.
- **`tasks`**: Project-specific scripts (e.g., `dev`, `build`, `test`).
- **`dependsOn`**: Internal dependencies on other projects in the monorepo.

## Environment Variables

We use `.env` files to manage environment-specific configurations.

- **`.env.example`**: Included in each project template. This file lists all required environment variables with placeholder values.
- **`.env`**: You should create this file in your project root (copied from `.env.example`) and populate it with your local secrets and configuration. **This file is git-ignored.**

## Package Scripts

While moonrepo is the primary task runner, we also use `package.json` scripts for standard Node.js workflows. Moon is configured to infer tasks from these scripts automatically where appropriate.
