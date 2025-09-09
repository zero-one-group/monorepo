include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/elastic-ip"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
}

inputs = {
  eip_name = local.project_name
  elastic_ips = {
    "bastion" = {}
    "nginx-dev-staging" = {}
    # Add any other Elastic IPs you need here:
    # "nginx" = {}
    # "other-service" = {}
  }
}
