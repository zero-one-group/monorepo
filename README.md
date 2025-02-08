# Zero One Group Monorepo

## Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
  - [Golang application](#golang-application)
  - [Local development server](#local-development-server)
  - [Creating application from template](#creating-application-from-template)
  - [Moon commands](#moon-commands)
- [E2E Testing](#e2e-testing)
- [Managing Dependencies](#managing-dependencies)
  - [Updating dependencies](#updating-dependencies)
  - [Cleanup projects](#cleanup-projects)
- [Tasks to Complete](#tasks-to-complete)
- [Contributions](#contributions)

## Overview

This repository contains framework projects utilizing [moonrepo][moonrepo] and the technology
stack commonly used within Zero One Group, with TypeScript and Go as the main languages.

To get started, several templates for backend and frontend applications are available.
The frontend stack includes React with the Vite bundler, as well as Next.js. A shared
UI Library, shadcn/ui, is also provided.

Additionally, the following tools and libraries are included:

- **Tailwind CSS** for utility-first CSS framework.
- **Radix UI** for accessible, unstyled UI components.
- **Vitest** for fast unit testing.
- **Playwright** for end-to-end testing.
- **Biome** for code formatting and linting.

## Quick Start

To begin, we suggest installing `moon` globally, read the documentation [here][moonrepo].
Then, follow these steps (_don't forget to replace `moon-project` with your project name_):

1. Clone this repository: `npx tiged zero-one-group/monorepo moon-project`
2. Install the necessary dependencies: `cd moon-project && pnpm install`
3. Create `.env` file or duplicate the `.env.example` file, then configure required variables.

Find and replace the `myorg` namespace and `example.com` string with your own organization
or project-specific namespace. This is necessary to ensure that all configurations,
dependencies, and references are correctly aligned with your project's unique identifier.
This includes updating any configuration files, package names, and other references
throughout the codebase where `myorg` is used.

### Golang application

Currently, Go is not supported as an official moonrepo toolchain. You need to manually
install and configure it for your project. Please read the [Go installation docs][go-docs].

For a list of supported toolchains, visit [moonrepo documentation][moon-toolchain].

### Local development server

This repository includes a local development server for testing and development purposes.
Currently, it supports PostgreSQL, Valkey (drop-in replacement for Redis), [mailpit][mailpit] (SMTP server),
and [pgweb][pgweb] (PostgreSQL web interface).

These commands are used for managing the local development server:

```sh
pnpm compose:up       # Start local development server
pnpm compose:down     # Stop local development server
pnpm compose:cleanup  # Remove all local development server data
```

### Creating application from template

To get started, you can use the following command to generate a new application from a template.

```sh
moon generate TEMPLATE_NAME
```

Example, creating React application:

```sh
moon generate template-react-app
```

View all available templates by looking at folders with `template-` prefix in this repository.

Current available templates are:

| Template Name        | Description                                       |
|----------------------|---------------------------------------------------|
| `template-golang`    | Basic Go application for backend                  |
| `template-react-app` | SPA React Router application with Tailwind CSS    |
| `template-react-ssr` | SSR React Router application with Tailwind CSS    |
| `template-shared-ui` | Collections of UI components based on `shadcn/ui` |

### Moon commands

After setting up your project, use the following commands for common tasks:

| Command                    | Description                              |
|----------------------------|------------------------------------------|
| `moon :dev`                | Start developing the project             |
| `moon :build`              | Build all projects                       |
| `moon :test`               | Run tests in all projects                |
| `moon :lint`               | Lint code in all projects                |
| `moon :format`             | Format code in all projects              |
| `moon <project_id>:env`    | Print system env and individual project  |
| `moon <project_id>:<task>` | Run specific task by project             |
| `moon check <project_id>`  | Run check for individual project         |
| `moon check --all`         | Run check for all tasks                  |
| `moon run '#tag:task'`     | Run a task in all projects with a tag    |
| `moon project-graph`       | Display an interactive graph of projects |

Type `moon help` for more information. Refer to the [moon tasks documentation](https://moonrepo.dev/docs/run-task) for more details.

## E2E Testing

This monorepo includes E2E tests for testing the application, powered by Playwright.
To run E2E tests, you need to install Playwright dependencies. You can do this by
running the following command:

**Install Playwright dependencies for all projects**

```sh
moon <project_id>:e2e-install
```

**Install Playwright dependencies for a specific project**

```sh
moon <project_id>:e2e-install
```

Run E2E tests for specific project in headless mode:

```sh
moon <project_id>:e2e
```

If you want to use Playwright UI mode, you can use the following command:

```sh
moon <project_id>:e2e-ui
```

To run E2E test for specific browser, you can use the following command:

```sh
moon <project_id>:e2e-chrome   # Run E2E test for Chrome browser
moon <project_id>:e2e-firefox  # Run E2E test for Firefox browser
moon <project_id>:e2e-mobile   # Run E2E test for Chrome Mobile browser
moon <project_id>:e2e-safari   # Run E2E test for Safari browser
```

## Managing Dependencies

To add a new dependency to a project, you can use the following command:

```sh
pnpm --filter <project_id> add <dependency>
```

Or, if you want to add development dependencies, you can use the following command:

```sh
pnpm --filter <project_id> add -D <dependency>
```

Example:

```sh
pnpm --filter react-app add -D vitest
```

### Updating dependencies

To update all projects dependencies, you can use the following command:

```sh
pnpm run update-deps
```

### Cleanup projects

Sometimes it is necessary to clean up dependencies and build artifacts from the project.
To do this, you can use the following command:

```sh
pnpm run cleanup
```

After all, you can reinstall the dependencies and build the project.

## Tasks to Complete

After creating a new project from this template repository, ensure you update the documentation, including:

1. **Project Overview:** Briefly describe the project's purpose and features.
2. **Installation Instructions:** Update steps to reflect any project-specific changes.
3. **Usage Guide:** Provide instructions on how to use the project, including commands and configurations.
4. **API Documentation:** If applicable, update API endpoints, request/response formats, and examples.
5. **Contributing Guidelines:** Reflect any new processes or requirements for contributions.
6. **License Information:** Ensure the license is accurate for the new project.

Keeping documentation current helps others understand, use, and contribute to the project.

## Contributions

Contributions are welcome! Please open a pull requests for your changes and tickets in case you would like to discuss something or have a question.

Read [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed documentation.

<!-- link reference definition -->
[go-docs]: https://go.dev/doc/install
[mailpit]: https://mailpit.axllent.org/
[moon-toolchain]: https://moonrepo.dev/docs/concepts/toolchain
[moonrepo]: https://moonrepo.dev/docs/install
[pgweb]: https://sosedoff.github.io/pgweb
