include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/vpc"
}

dependency "vpc_base" {
  config_path = "../../shared/vpc-base"
}

inputs = {
  vpc_name             = "lgtm"
  environment          = "prod"
  vpc_id               = dependency.vpc_base.outputs.vpc_id
  internet_gateway_id  = dependency.vpc_base.outputs.internet_gateway_id
  subnet_offset        = 0 # Prod subnets will be 10.201.1.0/24, 10.201.2.0/24, 10.201.3.0/24
}
