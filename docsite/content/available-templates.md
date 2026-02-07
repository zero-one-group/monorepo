---
title: Available Templates
slug: "available-templates"
---

The monorepo comes with a variety of templates to kickstart your development. These templates are pre-configured with the necessary tooling and best practices.

## Frontend

### Web Applications
- **`astro`**: Astro website template.
- **`nextjs`**: Next.js application template.
- **`react-app`**: React Single Page Application (SPA) template.
- **`react-ssr`**: React Server-Side Rendering (SSR) template.
- **`tanstack-start`**: TanStack Start template.

### Mobile
- **`expo`**: Expo (React Native) mobile application template.

## Backend

### Go
- **`go-clean`**: Go application following Clean Architecture.
- **`go-modular`**: Modular Go application structure.

### Python
- **`fastapi-ai`**: FastAPI application optimized for AI/ML workloads.

### Headless CMS
- **`strapi`**: Strapi executable for content management.

### Other
- **`phoenix`**: Elixir Phoenix application template.

## Infrastructure & DevOps

- **`ansible`**: Ansible playbooks and configuration.
- **`gitlab-cicd`**: GitLab CI/CD pipeline configurations.
- **`load-balancer`**: Load balancer configuration.
- **`monitoring`**: Monitoring stack (Prometheus, Grafana, etc.).
- **`postgresql`**: PostgreSQL database configuration.
- **`squidproxy`**: Squid proxy configuration.
- **`swarm`**: Docker Swarm configuration.
- **`terragrunt`**: Terragrunt infrastructure as code.

## Libraries

- **`shared-ui`**: Shared UI component library.

## Usage

To use any of these templates, run:

```bash
moon generate <template-name>
```

For example: `moon generate nextjs`
