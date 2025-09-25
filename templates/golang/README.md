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

#### Running Migration

1. Create new migration file
```bash
moon run migration-create -- {migration_name}
```

2. Migration up
```bash
moon run migration-up
```

3. Migration down
```bash
moon run migration-down
```

4. Migration reset
```bash
moon run migration-reset
```

4. Check Migration version
```bash
moon run migration-version
```

#### Running Seeders

1. Run seeders for all tables
```bash
moon run seed -- all
```

2. Run seeder for certain table
```bash
moon run seed -- {table_name}
```
#### Running Tests

##### 1. Install mockery (v3.5.1)

We use [mockery](https://github.com/vektra/mockery) to generate interface mocks. Make sure you have exactly v3.5.1:

Option A – via moon command: `moon {{ package_name | kebab_case }}:install-mockery`
Option B – via GitHub binary

Verify you have the right version:
```bash
mockery --version
# ⇒ mockery version 3.5.1
```

##### 2. Run the test suite

```bash
moon run {{ package_name | kebab_case }}:test
```

NOTE: Everytime test run, it will automatically generate mock

##### 3. Generate documentation

```bash
moon run {{ package_name | kebab_case }}:generate-swagger
```

## Production

### Instrumentation
Tracing is enabled exclusively in the production environment. Set `APP_ENVIRONMENT` to `production` to activate tracing. Alternatively, you may customize the tracing rules in `apps/{{ package_name | kebab_case }}/config/tracer.go`.

For instructions on customizing span tracing, please refer to the example located at:
- `apps/{{ package_name | kebab_case }}/internal/rest/user.go`
    - From rest layer all the way down to repository layer

