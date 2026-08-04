include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/ecs"
}

dependency "vpc" {
  config_path = "../vpc"
}

dependency "security_groups" {
  config_path = "../security-groups"
}

dependency "ecr" {
  config_path = "../../shared/ecr"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
  variables    = yamldecode(file("${get_parent_terragrunt_dir()}/shared/variables/ecs/variables.yaml"))
  env_dir      = "${get_parent_terragrunt_dir()}/shared/variables/ecs/env/${local.environment}"

  # CI/CD overrides this per deployment (see infra/app/ecs-deploy.yml) so
  # parallel dev/staging/prod deploys never race on the shared variables.yaml.
  image_tag = get_env("ECS_IMAGE_TAG", local.variables.image_tag)

  # Per-service env vars from YAML files (populated by CI pipeline)
  env_app = try(yamldecode(file("${local.env_dir}/{{ app_name }}.yml")), {})
}

inputs = {
  aws_account_id     = get_aws_account_id()
  project_name       = local.project_name
  vpc_id             = dependency.vpc.outputs.vpc_id
  subnet_ids         = dependency.vpc.outputs.subnet_public_all
  security_group_ids = [dependency.security_groups.outputs.sg_ids.ecs]
  ecr_repository_url = dependency.ecr.outputs.repository_url
  image_tag          = local.image_tag
  dns_namespace_name = local.variables.dns_namespace_name

  # Optional ALB: created once, one target group + listener rule per
  # service with enable_alb = true. The module also creates the task
  # security group that allows traffic from the ALB.
  alb_subnet_ids = dependency.vpc.outputs.subnet_public_all

  services = {
    # Service key must match the tag suffix pushed by the build pipeline
    # (<tag>-{{ app_name }}). Extend this map for more services, e.g.:
    #   worker = {
    #     cpu                      = local.variables.task_sizes["worker"].cpu
    #     memory                   = local.variables.task_sizes["worker"].memory
    #     enable_service_discovery = true
    #   }
    "{{ app_name }}" = {
      cpu                   = local.variables.task_sizes["{{ app_name }}"].cpu
      memory                = local.variables.task_sizes["{{ app_name }}"].memory
      environment           = { for k, v in local.env_app : k => tostring(v) }
      assign_public_ip      = true
      enable_alb            = true
      container_port        = 8000
      alb_health_check_path = "/health"
    }
  }

  # On-demand workloads: task definition only, no ECS service.
  # Launch them with ecs:RunTask when needed.
  jobs = {}
}
