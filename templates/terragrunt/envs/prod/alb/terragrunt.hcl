include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/alb"
}

dependency "vpc" {
  config_path = "../vpc"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  project_name = local.project_name
  environment  = local.environment
  vpc_id       = dependency.vpc.outputs.vpc_id
  subnet_ids   = dependency.vpc.outputs.subnet_public_all
  # One entry per ECS service exposed through the ALB. The key must match the
  # ECS service key (and container_port the service's port), since the ECS
  # unit picks up the target group via dependency.alb.outputs.target_group_arns.
  services = {
    "{{ app_name }}" = {
      container_port        = 8000
      alb_health_check_path = "/health"
    }
  }
}