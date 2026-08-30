# Terragrunt

Provisioning infrastructure as code via [Terragrunt](https://terragrunt.gruntwork.io/).

## Structure

```text
.
├── modules/                # Reusable Terraform modules
│   ├── acm/                # ACM certificates (placeholder - not implemented yet)
│   ├── alb/                # Application Load Balancer (public/private, HTTP/HTTPS)
│   ├── ec2/                # EC2 instances (Swarm, monitoring, etc. - legacy)
│   ├── ecr/                # ECR repositories
│   ├── ecs/                # ECS Fargate cluster, task definitions, services, IAM
│   ├── elastic-ip/         # Elastic IPs
│   ├── keypair/            # SSH key pair
│   ├── rds-cluster/        # RDS Aurora cluster (placeholder - not implemented yet)
│   ├── rds-single/         # RDS single instance (placeholder - not implemented yet)
│   ├── s3/                 # S3 buckets
│   ├── security-group/     # Security groups
│   └── vpc/                # VPC with public + private subnets and NAT gateway
├── variables/              # Shared variables (ec2, ecs)
├── envs/                   # Per-environment units
│   ├── shared/             # Cross-environment units (not tied to one env)
│   │   ├── ecr/            # ECR repository for application images
│   │   └── s3/             # S3 buckets (monitoring, public/private assets)
│   ├── dev/  staging/  prod/
│   │   ├── acm/            # (placeholder - not implemented yet)
│   │   ├── alb/            # ALB for the environment
│   │   ├── ec2/            # EC2 instances (legacy)
│   │   ├── ec2-worker/     # Additional workers (dev only)
│   │   ├── ecs/            # ECS Fargate cluster + services for the environment
│   │   ├── elastic-ip/     # Elastic IPs
│   │   ├── keypair/        # SSH key pair
│   │   ├── rds-cluster/    # (placeholder - not implemented yet)
│   │   ├── rds-single/     # (placeholder - not implemented yet)
│   │   ├── security-groups/# Environment security groups
│   │   └── vpc/            # Environment VPC (public + private subnets, NAT)
├── root.hcl                # Root configuration (remote state, provider, common inputs)
└── template.yml            # Moon generator template configuration
```

## ECS Fargate

Everything runs through a single module, **`modules/ecs`**, one instance per
environment (`envs/dev/ecs`, `envs/staging/ecs`, `envs/prod/ecs`):

- ECS cluster with Fargate/Fargate Spot capacity providers.
- Per-entry task definitions, ECS services, IAM roles (execution + per-entry
  task roles, least-privilege) and CloudWatch log groups.
- Tasks run in the **private subnets** of the environment VPC; outbound
  internet access goes through the NAT gateway.
- **Public traffic** reaches services through the environment **ALB**
  (`envs/<env>/alb` + `modules/alb`): the ALB unit creates the load balancer
  and one target group per service, and the ECS unit wires each service to its
  target group via `alb_target_group_arn`. The tasks security group only opens
  the service's container port to the ALB security group.
- **`services`** map: always-on workloads (task definition + ECS service).
  Each service can optionally enable:
  - **ALB** (`alb_target_group_arn` + `container_port`) — the module attaches
    the service to the target group created by the ALB unit.
  - **Service discovery** (`enable_service_discovery = true`) — a private Cloud
    Map namespace so tasks reach each other by a stable DNS name, e.g.
    `<service>.<dns_namespace_name>`.
- **`jobs`** map: on-demand workloads (task definition only, no ECS service),
  launched with `ecs:RunTask` when needed.
- Images are pulled from the **ECR** repository created by the
  `envs/shared/ecr` unit (`${project_name}-app`). The default image reference is
  `<ecr_repository_url>:<image_tag>-<service>`, matching the
  `<tag>-{{ app_name }}` tag pushed by the build pipeline.

### Applying

Apply the units in order (states live under `envs/` in the state bucket).
`envs/shared/` holds units that are not tied to a single environment, so they
are applied once and reused by every environment:

```bash
terragrunt apply --working-dir envs/dev/vpc
terragrunt apply --working-dir envs/shared/ecr
terragrunt apply --working-dir envs/dev/security-groups
terragrunt apply --working-dir envs/dev/elastic-ip
terragrunt apply --working-dir envs/dev/alb
terragrunt apply --working-dir envs/dev/ecs
```

### Customizing

Tunables live in `variables/ecs/variables.yaml`:

| Key | Description |
| --- | --- |
| `task_sizes` | CPU/memory per service (cpu in vCPU-1024ths, memory in MiB) |
| `dns_namespace_name` | Private Cloud Map namespace for service discovery |
| `image_tag` | Default image tag; overridden by CI via `ECS_IMAGE_TAG` |

Per-service environment variables come from
`variables/ecs/env/<environment>/<service>.yml` (git-ignored, populated
by the CI pipeline). The per-env `ecs/terragrunt.hcl` reads them with a `try()`
guard, so a missing file just means an empty environment. A service can use the
`image` field to override the auto-built image reference.

Note: service keys must match the image tag suffix your build pushes
(`<tag>-<service>`). The example uses `{{ app_name }}`. The same key must exist
in the `envs/<env>/alb` unit with the service's `container_port`, so the ECS
unit can pick up the matching `target_group_arns["<service>"]`.

### CI/CD

The [gitlab-cicd](../gitlab-cicd) template includes an ECS deployment pipeline
(`app/ecs-deploy.yml`). Build jobs push the image to ECR with a version tag,
then deploy jobs run `terragrunt apply` targeted at the service's task
definition + ECS service — converting the app `.env` into the per-environment
env YAML and passing the image tag via `ECS_IMAGE_TAG` so parallel
dev/staging/prod deploys never race on the same files. Configure the `APP_ENV_FILE`
and AWS CI/CD variables documented in that file.
