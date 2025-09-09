include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/vpc"
}

dependency "vpc_base" {
  config_path = "../../shared/vpc-base"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  vpc_name             = local.project_name
  environment          = local.environment
  vpc_id               = dependency.vpc_base.outputs.vpc_id
  internet_gateway_id  = dependency.vpc_base.outputs.internet_gateway_id
  subnet_offset        = 6 # Dev subnets will be 10.201.7.0/24, 10.201.8.0/24, 10.201.9.0/24
}
