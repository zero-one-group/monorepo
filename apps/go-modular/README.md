# Go Application

This template already pre-configured basic auth, take a look at `apps/go-modular/database/seeders/user_factory.go`.

## Common Tasks
```sh
moon go-modular:dev                  # Run development
moon go-modular:run                  # Execute `go run`
moon go-modular:build                # Build the application
moon go-modular:start                # Start application from build
moon go-modular:test                 # Run testing
moon go-modular:coverage             # Run test coverage
moon go-modular:format               # Run code formatting
moon go-modular:tidy                 # Install dependencies
moon go-modular:docker-build         # Build docker image
moon go-modular:docker-run           # Run docker image
moon go-modular:docker-shell         # Execute docker shell
moon go-modular:dump                 # Dump the database
moon go-modular:install-mockery      # Install mockery
moon go-modular:generate-swagger     # Generate Swagger OpenAPI docs
moon go-modular:generate-mock        # Generate Mock
```

## Migration Tasks
```sh
# Initiate or reset migrations and seed
moon go-modular:run -- migrate:reset --up --seed --force

# Common migration commands
moon go-modular:run -- migrate:up
moon go-modular:run -- migrate:status
moon go-modular:run -- migrate:version
moon go-modular:run -- migrate:create [MIGRATION_NAME]
moon go-modular:run -- migrate:down
moon go-modular:run -- migrate:reset
moon go-modular:run -- migrate:seed
```

## Generate Sample Configuration
This command will generate `.env.example` from `apps/go-modular/internal/config/default.go`:

```sh
moon go-modular:run -- generate:config
```

## Scaffold New Module
```sh
mkdir -p apps/go-modular/modules/dummy
mkdir -p apps/go-modular/modules/dummy/handler
mkdir -p apps/go-modular/modules/dummy/models
mkdir -p apps/go-modular/modules/dummy/repository
mkdir -p apps/go-modular/modules/dummy/services
echo 'package handler' > apps/go-modular/modules/dummy/handler/handler.go
echo 'package models' > apps/go-modular/modules/dummy/models/model.go
echo 'package models' > apps/go-modular/modules/dummy/models/schema.go
echo 'package repository' > apps/go-modular/modules/dummy/repository/repository.go
echo 'package services' > apps/go-modular/modules/dummy/services/services.go
echo 'package dummy' > apps/go-modular/modules/dummy/module.go
```

> Replace `dummy` with the module name
