# Environment variables for ECS services

This directory holds the per-service, per-environment environment variables
injected into each ECS task definition. Actual values are git-ignored - never
commit secrets here.

## Layout

    env/
      dev/<service>.yml       # development values
      staging/<service>.yml   # staging values
      prod/<service>.yml      # production values

The per-env terragrunt.hcl (e.g. `envs/dev/ecs/terragrunt.hcl`) resolves
`env/<environment>/<service>.yml` via a `try()` guard, so a missing file simply
means an empty environment.

## Format

Flat key-value YAML, injected into the task definition
`container_definitions.environment` as `{ name, value }` pairs:

    TEMPORAL_HOST_PORT: "temporal.example.com:443"
    OPENAI_API_KEY: "sk-..."

## How it gets populated

`ecs-deploy.yml` in the GitLab CI/CD template converts a GitLab CI variable
(the app's `.env` content) to YAML and writes it to `env/<environment>/<service>.yml`
before running `terragrunt apply`. Because the file is written per environment,
parallel dev/staging/prod deploys never race on the same file.

## Manual usage

You can also drop a YAML file here and run terragrunt directly - the per-env
config picks it up automatically.
