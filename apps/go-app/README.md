# Golang Project Template
Basic Go application for backend

### Prerequisites

- Golang v1.24
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- Postgres v17.x.x
- Docker and Docker Compose (for containerized development)

### Quick Start

#### Create New Application

1. Generate a new application from a template.

```bash
moon generate template-golang
```

2. Set up environment variables:
```bash
cp .env.example .env
```
Adjust the variable to the desired value.

3. Install dependencies:
```bash
moon run tidy
```

#### Running The Application

1. Running on development mode
```bash
moon run dev
```

2. Build the application
```bash
moon run build
```

3. Running on production mode
```bash
moon run start
```
