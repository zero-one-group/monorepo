locals {
  environment  = basename(dirname(get_terragrunt_dir()))
  project_name = "{{ project_name }}"
  aws_region   = "{{ region }}"
}

remote_state {
  backend = "s3"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
  config = {
    bucket       = "${local.project_name}-bucket-states"
    key          = "${path_relative_to_include()}/terraform.tfstate"
    region       = local.aws_region
    encrypt      = true
    use_lockfile = true
  }
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents = <<EOF
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.13.0"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_region" "current" {}
EOF
}

inputs = {
  project_name = local.project_name
  region       = local.aws_region

  common_tags = {
    Project     = local.project_name
    Environment = local.environment
    Contact     = "{{ author }}"
    ManagedBy   = "Terragrunt"
    Version     = "1.0.0"
    CreatedBy   = "zero-one-group"
    LastModified = formatdate("YYYY-MM-DD", timestamp())
  }
}
