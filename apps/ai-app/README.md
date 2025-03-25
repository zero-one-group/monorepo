# FastAPI AI Project Template

## Project Structure

| Directory/File | Purpose |
|----------------|---------|
| `app/` | Main application package |
| `app/core/` | Core functionality including config, database, env management, and logging |
| `app/main.py` | Application entry point |
| `app/model/` | Data models including database models and schemas |
| `app/repository/` | Data access layer for different providers (postgres, ml models, etc.) |
| `app/router/` | API route definitions |
| `app/services/` | Business logic implementation |
| `infra/` | Infrastructure configuration |
| `infra/docker/` | Docker configuration for development and production |
| `tests/` | Test suite |
| `.env.example` | Example environment variables |
| `moon.yml` | Moonrepo configuration |
| `pyproject.toml` | Python package configuration |

## Check-In-Dance

### Prerequisites

- Python 3.10+
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- Docker and Docker Compose (optional, for containerized development)

### Installation

1. Create a new project using this template:

```bash
moon generate template-fastapi-ai
```

2. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:

```bash
moon ai-app:sync
```

### Running the Application

#### Development Mode

```bash
# Start the development server with hot reloading
moon ai-app:dev
```

#### Production Mode

```bash
moon ai-app:start
```

#### Using Custom App Name

If you've renamed your application, use:

```bash
moon {app-name}:dev
```

## Contributing

When contributing to this template, please follow the existing structure and naming conventions.
