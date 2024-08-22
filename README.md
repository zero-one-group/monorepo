# Zero One Group Monorepo

This repository contains framework projects utilizing [moonrepo][moonrepo] and the technology
stack commonly used within Zero One Group, with TypeScript and Go as the main languages.

To get started, several templates for backend and frontend applications are available.
The frontend stack includes React with the Vite bundler, as well as Next.js. A shared
UI Library, shadcn, is also provided.

Additionally, the following tools and libraries are included:

- **Tailwind CSS** for utility-first CSS framework.
- **Radix UI** for accessible, unstyled UI components.
- **Vitest** for fast unit testing.
- **Playwright** (optional) for end-to-end testing.
- **Biome** for code formatting and linting.

## 🏁 Quick Start

To begin, we suggest installing moon globally, read the documentation [here](https://moonrepo.dev/docs/install).
Then, follow these steps (_don't forget to replace `my-new-project` with your project name_):

1. Clone this repository: `npx tiged zero-one-group/monorepo my-new-project`
2. Initialize git repository: `cd my-new-project && git init`
3. Install the required toolchain: `moon setup`
4. Install the necessary dependencies: `pnpm install`
5. Add template submodule: `git submodule add https://github.com/zero-one-group/templates`
6. Initialize template submodule: `pnpm submodule:init`
7. Update submodule (optional): `pnpm submodule:update`
8. Create `.env` file or duplicate the `.env.example` file, then configure required variables.

Find and replace the `myorg` namespace and `example.com` string with your own organization
or project-specific namespace. This is necessary to ensure that all configurations,
dependencies, and references are correctly aligned with your project's unique identifier.
This includes updating any configuration files, package names, and other references
throughout the codebase where `myorg` is used.

Optinally, you'll need to install [Static Web Server][static-web-server] to preview
generated websites or Single-Page Applications (SPAs).

### Golang application

Currently, Go is not supported as an official moonrepo toolchain. You need to manually
install and configure it for your project. Please read the [Go installation docs][go-docs].

For a list of supported toolchains, visit [moonrepo documentation][moon-toolchain].

### Creating application from template

To get started, you can use the following command to generate a new application from a template.

```sh
moon generate TEMPLATE_NAME
```

Example, creating React application:

```sh
moon generate moon-vite-react-tailwind
```

Explore the [`templates`](./templates/) directory to see all available templates.
Each template is prefixed with `moon-` to indicate its purpose and usage. The original
templates repository can be found at [`zero-one-group/templates`][zog-templates].

### Moon commands

After setting up your project, use the following commands for common tasks:

| Command                    | Description                              |
|----------------------------|------------------------------------------|
| `moon :dev`                | Start developing the project             |
| `moon :build`              | Build all projects                       |
| `moon :test`               | Run tests in all projects                |
| `moon :lint`               | Lint code in all projects                |
| `moon :format`             | Format code in all projects              |
| `moon <project_id>:<task>` | Run specific task by project             |
| `moon check <project_id>`  | Run check for individual project         |
| `moon check --all`         | Run check for all tasks                  |
| `moon run '#tag:task'`     | Run a task in all projects with a tag    |
| `moon project-graph`       | Display an interactive graph of projects |

Type `moon help` for more information. Refer to the [moon tasks documentation](https://moonrepo.dev/docs/run-task) for more details.

## 📦 Managing Dependencies

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
pnpm --filter website add -D vitest
```

### Updating dependencies

To update all projects dependencies, you can use the following command:

```sh
pnpm run update-deps
```

### Cleanup Projects

Sometimes it is necessary to clean up dependencies and build artifacts from the project.
To do this, you can use the following command:

```sh
pnpm run cleanup
```

Cleanup workspace (optional):

```sh
pnpm run cleanup-workspace
```

After all, you can reinstall the dependencies and build the project.

## ✅ Tasks to Complete

After creating a new project from this template repository, ensure you update the documentation, including:

1. **Project Overview:** Briefly describe the project's purpose and features.
2. **Installation Instructions:** Update steps to reflect any project-specific changes.
3. **Usage Guide:** Provide instructions on how to use the project, including commands and configurations.
4. **API Documentation:** If applicable, update API endpoints, request/response formats, and examples.
5. **Contributing Guidelines:** Reflect any new processes or requirements for contributions.
6. **License Information:** Ensure the license is accurate for the new project.

Keeping documentation current helps others understand, use, and contribute to the project.

<!-- link reference definition -->
[moonrepo]: https://moonrepo.dev/
[zog-templates]: https://github.com/zero-one-group/templates
[moon-toolchain]: https://moonrepo.dev/docs/concepts/toolchain
[go-docs]: https://go.dev/doc/install
[static-web-server]: https://static-web-server.net/download-and-install/
