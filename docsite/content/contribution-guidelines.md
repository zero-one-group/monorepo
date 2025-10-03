---
title: Contribution Guidelines
slug: "contribution-guidelines"
---

We welcome contributions to the Zero One Group Monorepo! This document outlines the guidelines for
contributing to ensure a smooth and collaborative development process.

## Table of Contents

1.  [Getting Started](#getting-started)
    *   [Prerequisites](#prerequisites)
    *   [Local Development Setup](#local-development-setup)
2.  [Project Structure and Naming Conventions](#project-structure-and-naming-conventions)
    *   [App Naming](#app-naming)
    *   [Development Folder (`apps/`)](#development-folder-apps)
    *   [Template Folders](#template-folders)
3.  [Development Workflow](#development-workflow)
    *   [Creating a New Application](#creating-a-new-application)
    *   [Moving Changes to Templates](#moving-changes-to-templates)
4.  [Helpful Moon Commands](#helpful-moon-commands)
5.  [Code Style and Standards](#code-style-and-standards)
6.  [Pull Request Guidelines](#pull-request-guidelines)
7.  [Seeking Help](#seeking-help)
8.  [Related Helpful Links](#related-helpful-links)

---

## 1. Getting Started

To get your development environment set up, please follow these steps.

### Prerequisites

Before you begin, ensure you have the following installed:

*   **Node.js** (LTS version recommended)
*   **pnpm**
*   **moon** (install globally as recommended by moonrepo):
    ```bash
    pnpm install -g @moonrepo/cli
    ```
*   **Go** (if you're working on Go applications)
*   **uv** (if you're working on Python applications)

### Local Development Setup

1.  **Clone the monorepo:**

2.  **Install dependencies:**
    ```bash
    pnpm install
    ```

3.  **Configure `.env` files:**
    Create a `.env` file or duplicate `.env.example` inside the specific application or package you are working on.

4.  **Start Local Development Server (Optional):**
    ```bash
    pnpm compose:up    # Start local development server
    pnpm compose:down  # Stop local development server
    pnpm compose:cleanup # Remove all local development server data
    ```
    To start a development server for a specific application you want to contribute to (e.g., `go-clean`):
    ```bash
    moon go-clean:dev
    ```

---

## 2. Project Structure and Naming Conventions

This monorepo follows specific conventions to maintain consistency and ease of management.

### App Naming

All applications within the monorepo should follow the `template-{app-name}` format.
**Example:** `template-go-clean`, `template-react-app`

### Development Folder (`apps/`)

The `apps/` directory is designated for active development. When you are developing a new application or making significant changes to an existing one, you should work within a sub-folder inside `apps/`. This allows for iteration and experimentation without affecting the stable templates.

### Template Folders

After development in `apps/` is complete and the application is stable and ready for wider use, the changes should be moved to the corresponding template folder (e.g., `templates/go-clean` if it's a Go application). This ensures that new projects can be easily bootstrapped from stable, well-defined templates.

---

## 3. Development Workflow

### Creating a New Template

To establish a new template, manually initialize the desired project within the `apps/` directory. Once the project is considered stable and ready as a template, migrate the code to a dedicated template folder following the established naming conventions (e.g., `template-{app-name}`).

Helpful Resource:
* [moonrepo Toolchain Configuration](https://moonrepo.dev/docs/config/toolchain)

### Contributing to an Existing Template

Select an existing template or create an issue to address the required updates. Create a branch and perform the necessary development work. After completing the updates, migrate the code to existing template folder following the appropriate naming conventions (e.g., `template-{app-name}`, `template-go-clean`).

Subsequently, submit a pull request in accordance with the [Pull Request Guidelines](#pull-request-guidelines).

**Note:** Ensure that all package names are updated accordingly. For example, after verifying the template, change the package name from `"go-clean/internal/rest/middleware"` to `"{{ package_name }}/internal/rest/middleware"` within the template code.

---

## 4. Helpful Moon Commands

**Adding helpful commands to `{template-folder}/moon.yml`:**

For project-specific commands that are useful for contributors, consider adding them directly to the `moon.yml` file within the respective template folder (e.g., `templates/go-clean/moon.yml`). This makes them easily discoverable and runnable by others using commands like `moon go-clean:dev`.

**Example for `templates/go-clean/moon.yml`:**

```yaml
# templates/go-clean/moon.yml
language: 'go'
type: 'application' # or 'library'
tasks:
  dev:
    command: 'go run ./cmd/server' # Example: command to start the Go app in dev mode
    local: true
    options:
      persistent: true
  build:
    command: 'go build -o dist/go-clean ./cmd/server' # Example: command to build the Go app
    inputs:
      - 'src/**/*.go'
      - 'go.mod'
      - 'go.sum'
    outputs:
      - 'dist'
  test:
    command: 'go test ./...' # Example: command to run Go tests
    inputs:
      - 'src/**/*.go'
      - 'go.mod'
      - 'go.sum'
```

*   Detailed link for moon project configuration: [https://moonrepo.dev/docs/config/project](https://moonrepo.dev/docs/config/project)

---

## 5. Code Style and Standards

We use **Biome** for code formatting and linting. Please ensure your code adheres to the configured standards by running the appropriate moon commands (e.g., `moon run :format`, `moon run :lint`) before submitting your changes.

*   Biome Documentation: [https://biomejs.dev/](https://biomejs.dev/)

---

## 6. Pull Request Guidelines

To ensure your contributions can be reviewed and merged efficiently, please follow these guidelines:

*   **Create an Issue:** We encourage you to first create an issue describing the feature or bug you plan to work on. You can optionally tag the relevant maintainer(s) in the issue.
*   **Branch Naming:** Use descriptive branch names following a convention (e.g., `feat/new-feature-name`, `fix/bug-description`, `docs/update-readme`).
*   **Commit Messages:** Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for clear and concise commit messages. This helps with automated changelog generation and understanding the history.
*   **Tests:** Ensure your changes are covered by appropriate tests (unit, integration, E2E). New features should have tests, and bug fixes should include a test that reproduces the bug.
*   **Documentation:** Update relevant documentation for any new features, significant changes, or breaking changes.
*   **Review:** All pull requests require at least one approval from a designated reviewer before merging.

---

## 7. Seeking Help

If you have any questions, encounter issues, or need clarification on any aspect of contributing, please:

*   Open an issue on the GitHub repository: [https://github.com/zero-one-group/monorepo/issues](https://github.com/zero-one-group/monorepo/issues)
*   Reach out on our internal communication channels (e.g. Slack).
*   For specific template-related questions, you can reach out to our maintainers:
    *   **General Monorepo:** `@riipandi`, `@rubiagatra`
    *   **Go Template:** `@ameliarahman`, `@mirfmaster`
    *   **Python Template:** `@mirfmaster`
    *   **Infra Template:** `@prihuda`

---

## 8. Related Helpful Links

*   **Zero One Group Monorepo Main README:** [https://github.com/zero-one-group/monorepo?tab=readme-ov-file](https://github.com/zero-one-group/monorepo?tab=readme-ov-file)
*   **moonrepo Documentation:** [https://moonrepo.dev/docs/](https://moonrepo.dev/docs/)
*   **moonrepo Tasks Documentation:** [https://moonrepo.dev/docs/run-task](https://moonrepo.dev/docs/run-task)
