include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/vpc"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  vpc_name           = local.project_name
  environment        = local.environment
  cidr_block         = "10.201.0.0/16"
  subnet_offset      = 3 # Staging subnets: public 10.201.4.0/24, 10.201.5.0/24, 10.201.6.0/24
  enable_nat_gateway = true
}