include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/ec2"
}

dependency "vpc" {
  config_path = "../vpc"
}

dependency "security_groups" {
  config_path = "../security-groups"
}

dependency "keypair" {
  config_path = "../keypair"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  variables = yamldecode(file("${get_parent_terragrunt_dir()}/shared/variables/ec2/variables.yaml"))
  priv_key = file(pathexpand("~/.ssh/id_rsa"))
}

inputs = {
  keyname        = dependency.keypair.outputs.keyname
  aws_account_id = get_aws_account_id()
  project_name   = local.project_name

  instances = {
    worker-staging-1 = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["worker-dev-staging"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.worker,dependency.security_groups.outputs.sg_ids.swarm,dependency.security_groups.outputs.sg_ids.monitoring,dependency.security_groups.outputs.sg_ids.postgresql]
      user_data_filename = local.variables.user_data_files["worker"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["worker"]
      s3_bucket_arns = [
        "arn:aws:s3:::lgtm-monitoring"
      ]
      ssm_parameter_paths = [
        "/lgtm/swarm/*",
      ]
      cluster_identifier = "app-cluster"
    }
  }
}
