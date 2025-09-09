include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/s3"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
}

inputs = {
  buckets = {
    monitoring = {
      bucket_name            = "${local.project_name}-monitoring"
      enable_versioning      = false
      enable_lifecycle_rule  = true
      transition_days        = 30
      expiration_days        = 90
      enable_lb_logging      = false
    }

    public_assets = {
      bucket_name            = "${local.project_name}-public-assets"
      enable_versioning      = false
      enable_lifecycle_rule  = false
      enable_public_read     = true
      cors_origin           = ["https://*.zero-one.cloud"]
    }

    private_assets = {
      bucket_name         = "${local.project_name}-private-assets"
      enable_versioning   = false
      cors_origin        = ["https://*.zero-one.cloud"]
      allow_public_policy = true
      custom_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
          {
            Effect = "Allow"
            Principal = "*"
            Action = ["s3:GetObject"]
            Resource = "BUCKET_ARN_PLACEHOLDER/*"
            Condition = {
              StringLike = {
                "aws:Referer" = [
                  "https://*.zero-one.cloud/*"
                ]
              }
            }
          }
        ]
      })
    }
  }
}

