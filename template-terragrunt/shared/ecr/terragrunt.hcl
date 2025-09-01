include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/ecr"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
}

inputs = {
  repo_name                       = "${local.project_name}-app"
  untagged_image_expiration_days  = 1
}
