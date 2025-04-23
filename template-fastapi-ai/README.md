# FastAPI AI Project Template

Short brief description about the project.

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
moon {{ package_name | kebab_case }}:dev
```
## Useful Links
- [FastAPI Dependency Injection](https://fastapi.tiangolo.com/tutorial/dependencies/)
    - tl;dr: FastAPI Dependency Injection: Your code tells FastAPI what it needs, and FastAPI automatically provides those requirements when needed. This eliminates repetitive code and makes your application cleaner and easier to maintain.
- [FastAPI Deployment](https://fastapi.tiangolo.com/deployment/)
