include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../../modules/vpc"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  vpc_name         = local.project_name
  environment      = local.environment
  cidr_block       = "10.201.0.0/16"
  subnet_offset    = 6 # Public subnets: 10.201.7.0/24, 10.201.8.0/24, 10.201.9.0/24. Private: 10.201.107.0/24, ...
  enable_nat_gateway = true
}
