---
title: Create New Project
slug: "create-new-project"
---

Getting started with a new project in the monorepo is straightforward. We provide CLI tools to help you initialize the workspace and generate applications from pre-defined templates.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js**: Version 22.19 or later (managed via moon).
- **pnpm**: Version 10.18.0 or later.
- **Docker**: For running containerized services.
- **Moonrepo**: The build system and task runner.

## 1. Initialize a New Monorepo

If you are setting up a brand new monorepo instance, use the `moci` tool to initialize the workspace structure.

```bash
pnpm dlx moci init <project-name>
```

This command will set up the necessary configuration files and directory structure for a moonrepo-powered workspace.

## 2. Generate an Application

Once your monorepo is set up, you can generate new applications using the `moon generate` command. This command uses the templates defined in the workspace to scaffold a new project with all the necessary boilerplate.

```bash
moon generate <template-name>
```

Replace `<template-name>` with the name of the template you wish to use (e.g., `react-app`, `go-clean`, `nextjs`).

> [!TIP]
> You can view a list of all available templates in the [Available Templates](/available-templates) section.

### Example: creating a React App

To create a new React application, run:

```bash
moon generate react-app
```

Follow the interactive prompts to configure your new application, such as setting the project name and description.

## 3. Configure Your Application

After generation, you may need to perform some initial configuration:

1.  **Environment Variables**: Copy `.env.example` to `.env` in your new project's directory and update the values as needed.
2.  **Dependencies**: Run `pnpm install` or `moon setup` to install any new dependencies required by the generated application.

## 4. Start Development

With your application generated and configured, you are ready to start coding!

Use `moon` to run the development server:

```bash
moon run <project-name>:dev
```

Or, to run all applications in the workspace (if applicable):

```bash
moon :dev
```
