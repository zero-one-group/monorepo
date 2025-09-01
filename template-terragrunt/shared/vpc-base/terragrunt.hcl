include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/vpc-base"
}

inputs = {
  vpc_name = "base"
  cidr_blocks = "10.201.0.0/16"
}
