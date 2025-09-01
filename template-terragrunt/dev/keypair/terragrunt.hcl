include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/keypair"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  keyname         = "${local.project_name}-${local.environment}-keypair"
  public_key_path = "~/.ssh/id_rsa.pub"
}
