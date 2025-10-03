# Zero One Group Monorepo

A unified foundation for building modern backend and frontend applications with TypeScript and Go.
Powered by [moonrepo](https://moonrepo.dev/docs/install), this monorepo offers curated templates,
shared UI libraries, and essential developer tools to accelerate development and foster best practices
across teams.

**Motivation:**
We built this monorepo to accelerate development, promote best practices, and simplify collaboration
across teams by providing ready-to-use templates, consistent tooling, and a scalable structure for all
Zero One Group projects.

## Quick Start

1. **Clone the repository:**
   `pnpm dlx tiged zero-one-group/monorepo my-monorepo-project`
2. **Install dependencies:**
   `cd moon-project && pnpm install`
3. **Generate the application from template:**
   `moon generate TEMPLATE_NAME`
4. **Configure your environment:**
   copy `.env.example` to `.env` and adjust as needed.
5. **Start development:**
   run all by running command `moon :dev` or `moon '#app:dev'`

## Templates

This monorepo includes a wide range of templates to help you start new projects quickly and consistently.
Templates are available for backend (Go, FastAPI), frontend (Astro, Next.js, React), and infrastructure
(Strapi, shared UI libraries).

Each template is designed to follow best practices and comes with pre-configured tooling, recommended
folder structures, and example code to get you up and running fast. You can generate new apps from these
templates using provided commands, and customize them to fit your needs.

For the full list of available templates, and usage instructions, please refer to the
[documentation](https://oss.zero-one-group.com/monorepo/available-templates/).

## Learn More

For complete guides and advanced usage, visit: 👉 [https://oss.zero-one-group.com/monorepo](https://oss.zero-one-group.com/monorepo)

## Contributions

Contributions are welcome! Please open a pull request or ticket for questions and improvements.

Read the full guidelines at: 👉 [https://oss.zero-one-group.com/monorepo/contribution-guidelines](https://oss.zero-one-group.com/monorepo/contribution-guidelines)
