# Go Application

This template already pre-configured basic auth, take a look at `apps/{{ package_name | kebab_case }}/database/seeders/user_factory.go`.

## Common Tasks
```sh
moon {{ package_name | kebab_case }}:dev                  # Run development
moon {{ package_name | kebab_case }}:run                  # Execute `go run`
moon {{ package_name | kebab_case }}:build                # Build the application
moon {{ package_name | kebab_case }}:start                # Start application from build
moon {{ package_name | kebab_case }}:test                 # Run testing
moon {{ package_name | kebab_case }}:coverage             # Run test coverage
moon {{ package_name | kebab_case }}:format               # Run code formatting
moon {{ package_name | kebab_case }}:tidy                 # Install dependencies
moon {{ package_name | kebab_case }}:docker-build         # Build docker image
moon {{ package_name | kebab_case }}:docker-run           # Run docker image
moon {{ package_name | kebab_case }}:docker-shell         # Execute docker shell
moon {{ package_name | kebab_case }}:dump                 # Dump the database
moon {{ package_name | kebab_case }}:install-mockery      # Install mockery
moon {{ package_name | kebab_case }}:generate-swagger     # Generate Swagger OpenAPI docs
moon {{ package_name | kebab_case }}:generate-mock        # Generate Mock
```

## Migration Tasks
```sh
# Initiate or reset migrations and seed
moon {{ package_name | kebab_case }}:run -- migrate:reset --up --seed --force

# Common migration commands
moon {{ package_name | kebab_case }}:run -- migrate:up
moon {{ package_name | kebab_case }}:run -- migrate:status
moon {{ package_name | kebab_case }}:run -- migrate:version
moon {{ package_name | kebab_case }}:run -- migrate:create [MIGRATION_NAME]
moon {{ package_name | kebab_case }}:run -- migrate:down
moon {{ package_name | kebab_case }}:run -- migrate:reset
moon {{ package_name | kebab_case }}:run -- migrate:seed
```

## Generate Sample Configuration
This command will generate `.env.example` from `apps/{{ package_name | kebab_case }}/internal/config/default.go`:

```sh
moon {{ package_name | kebab_case }}:run -- generate:config
```

## Scaffold New Module
```sh
mkdir -p apps/{{ package_name | kebab_case }}/modules/dummy
mkdir -p apps/{{ package_name | kebab_case }}/modules/dummy/handler
mkdir -p apps/{{ package_name | kebab_case }}/modules/dummy/models
mkdir -p apps/{{ package_name | kebab_case }}/modules/dummy/repository
mkdir -p apps/{{ package_name | kebab_case }}/modules/dummy/services
echo 'package handler' > apps/{{ package_name | kebab_case }}/modules/dummy/handler/handler.go
echo 'package models' > apps/{{ package_name | kebab_case }}/modules/dummy/models/model.go
echo 'package models' > apps/{{ package_name | kebab_case }}/modules/dummy/models/schema.go
echo 'package repository' > apps/{{ package_name | kebab_case }}/modules/dummy/repository/repository.go
echo 'package services' > apps/{{ package_name | kebab_case }}/modules/dummy/services/services.go
echo 'package dummy' > apps/{{ package_name | kebab_case }}/modules/dummy/module.go
```

> Replace `dummy` with the module name
