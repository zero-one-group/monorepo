include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/full-vpc"
}

inputs = {
  vpc_name           = "lgtm"
  environment        = "dev"
  cidr_block         = "10.10.0.0/16"
  enable_nat_gateway = true
}
