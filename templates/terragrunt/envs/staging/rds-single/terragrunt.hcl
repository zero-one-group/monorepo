include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/rds-single"
}

inputs = {
  enabled = false
}
