# Terragrunt

Provisioning infrastructure as code via [Terragrunt](https://terragrunt.gruntwork.io/).

## Structure

```text
.
├── modules/                # Reusable Terraform modules
│   ├── ec2/                # EC2 instances (Swarm, monitoring, etc.)
│   ├── ecr/                # ECR repositories
│   ├── ecs/                # ECS cluster, task definitions, services, IAM, ALB
│   ├── full-vpc/           # Full VPC creation (non-production)
│   ├── keypair/            # SSH key pair
│   ├── s3/                 # S3 buckets
│   ├── security-group/     # Security groups
│   ├── vpc/                # Public subnets for a specific environment
│   └── vpc-base/           # References the existing base VPC and IGW
├── shared/                 # Configuration shared across environments
│   ├── ecr/                # ECR repository for application images
│   ├── elastic-ip/         # Elastic IPs
│   ├── s3/                 # Shared buckets
│   ├── variables/          # Shared variables (ec2, ecs)
│   └── vpc-base/           # Base VPC reference
├── dev/ staging/ prod/     # Per-environment units
│   ├── ec2/                # EC2 instances
│   ├── ec2-worker/         # Additional workers
│   ├── ecs/                # ECS Fargate cluster + services for the environment
│   ├── keypair/            # SSH key pair
│   ├── security-groups/    # Environment security groups
│   └── vpc/                # Environment subnets
├── root.hcl                # Root configuration (remote state, provider, common inputs)
└── template.yml            # Moon generator template configuration
```

## ECS Fargate

Everything runs through a single module, **`modules/ecs`**, one instance per
environment (`dev/ecs`, `staging/ecs`, `prod/ecs`):

- ECS cluster with Fargate/Fargate Spot capacity providers.
- Per-entry task definitions, ECS services, IAM task roles (with optional
  S3/SSM access) and CloudWatch log groups.
- **`services`** map: always-on workloads (task definition + ECS service).
  Each service can optionally enable:
  - **ALB** (`enable_alb = true` + `container_port`) — the module creates one
    ALB, a target group, a listener rule and the task/ALB security groups, so
    no separate ALB unit is needed.
  - **Service discovery** (`enable_service_discovery = true`) — a private Cloud
    Map namespace so tasks reach each other by a stable DNS name, e.g.
    `<service>.<dns_namespace_name>`.
- **`jobs`** map: on-demand workloads (task definition only, no ECS service),
  launched with `ecs:RunTask` when needed.
- Images are pulled from the **ECR** repository created by `shared/ecr`
  (`${project_name}-app`). The default image reference is
  `<ecr_repository_url>:<image_tag>-<service>`, matching the
  `<tag>-{{ app_name }}` tag pushed by the build pipeline.

### Applying

Apply the units in order:

```bash
terragrunt apply --terragrunt-working-dir shared/vpc-base
terragrunt apply --terragrunt-working-dir dev/vpc
terragrunt apply --terragrunt-working-dir shared/ecr
terragrunt apply --terragrunt-working-dir dev/security-groups
terragrunt apply --terragrunt-working-dir dev/ecs
```

### Customizing

Tunables live in `shared/variables/ecs/variables.yaml`:

| Key | Description |
| --- | --- |
| `task_sizes` | CPU/memory per service (cpu in vCPU-1024ths, memory in MiB) |
| `dns_namespace_name` | Private Cloud Map namespace for service discovery |
| `image_tag` | Default image tag; overridden by CI via `ECS_IMAGE_TAG` |

Per-service environment variables come from
`shared/variables/ecs/env/<environment>/<service>.yml` (git-ignored, populated
by the CI pipeline). The per-env `ecs/terragrunt.hcl` reads them with a `try()`
guard, so a missing file just means an empty environment. A service can use the
`image` field to override the auto-built image reference.

Note: service keys must match the image tag suffix your build pushes
(`<tag>-<service>`). The example uses `{{ app_name }}`.

### CI/CD

The [gitlab-cicd](../gitlab-cicd) template includes an ECS deployment pipeline
(`app/ecs-deploy.yml`). Build jobs push the image to ECR with a version tag,
then deploy jobs run `terragrunt apply` targeted at the service's task
definition + ECS service — converting the app `.env` into the per-environment
env YAML and passing the image tag via `ECS_IMAGE_TAG` so parallel
dev/staging/prod deploys never race on the same files. Configure the `APP_ENV_FILE`
and AWS CI/CD variables documented in that file.
