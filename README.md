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

## Quick Start

To begin, we suggest installing moon globally, read the documentation [here](https://moonrepo.dev/docs/install).
Then, follow these steps (_don't forget to replace `my-new-project` with your project name_):

1. Clone this repository: `npx tiged zero-one-group/monorepo#moon my-new-project`
2. Initialize git repository: `cd my-new-project && git init`
3. Install the required toolchain: `moon setup`
4. Install the necessary dependencies: `pnpm install`
5. Initialize template submodule: `pnpm submodule:init`
6. Update ubmodule (optional): `pnpm submodule:update`

Find and replace the `myorg` namespace string with your own organization or project-specific
namespace. This is necessary to ensure that all configurations, dependencies, and references
are correctly aligned with your project's unique identifier. This includes updating any
configuration files, package names, and other references throughout the codebase where
`myorg` is used.

### Creating application from template

To get started, you can use the following command to generate a new application from a template.

```sh
moon generate TEMPLATE_NAME
```

Example, creating React application:

```sh
moon generate moon-vite-react-tailwind
```

Take a look at [`templates`](./templates/) directory for list all available templates.

### Moon commands

Once installed, run the following commands for common tasks:

| Command                 | Description                      |
|-------------------------|----------------------------------|
| `moon check --all`      | Run all tasks                    |
| `moon :dev`             | Start developing the project     |
| `moon :build`           | Build all projects               |
| `moon :lint`            | Lint code in all projects        |
| `moon :test`            | Run tests in all projects        |
| `moon :format`          | Format code in all projects      |
| `moon <project>:<task>` | Run specific task by project     |
| `moon check <project>`  | Run check for individual project |

Refer to the [moon tasks documentation](https://moonrepo.dev/docs/run-task) for more details.

[moonrepo]: https://moonrepo.dev/

## Tasks to Complete

After creating a new project from this template repository, ensure you update the documentation, including:

1. **Project Overview:** Briefly describe the project's purpose and features.
2. **Installation Instructions:** Update steps to reflect any project-specific changes.
3. **Usage Guide:** Provide instructions on how to use the project, including commands and configurations.
4. **API Documentation:** If applicable, update API endpoints, request/response formats, and examples.
5. **Contributing Guidelines:** Reflect any new processes or requirements for contributions.
6. **License Information:** Ensure the license is accurate for the new project.

Keeping documentation current helps others understand, use, and contribute to the project.
