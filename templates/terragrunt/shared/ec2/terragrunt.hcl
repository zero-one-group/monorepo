include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/ec2"
}

dependency "vpc" {
  config_path = "../../dev/vpc"
}

dependency "security_groups" {
  config_path = "../../dev/security-groups"
}

dependency "keypair" {
  config_path = "../../dev/keypair"
}

dependency "elastic_ip" {
  config_path = "../elastic-ip"
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
    bastion = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["bastion"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.bastion]
      user_data_filename = local.variables.user_data_files["bastion"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["bastion"]
      enable_eip        = true
      eip_allocation_id = dependency.elastic_ip.outputs.eip_bastion_id
    }

    master-dev-staging-1 = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["master-dev-staging"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.master,dependency.security_groups.outputs.sg_ids.swarm,dependency.security_groups.outputs.sg_ids.monitoring,dependency.security_groups.outputs.sg_ids.postgresql]
      user_data_filename = local.variables.user_data_files["master"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["master"]
      s3_bucket_arns = [
        "arn:aws:s3:::lgtm-monitoring"
      ]
      ssm_parameter_paths = [
        "/lgtm/swarm/*",
      ]
      cluster_identifier = "app-cluster"
    }

    nginx-1 = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["nginx-dev-staging"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.lb,dependency.security_groups.outputs.sg_ids.swarm,dependency.security_groups.outputs.sg_ids.monitoring]
      user_data_filename = local.variables.user_data_files["nginx"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["nginx"]
      enable_eip        = true
      eip_allocation_id = dependency.elastic_ip.outputs.eip_nginx_dev_staging_id
      s3_bucket_arns = [
        "arn:aws:s3:::lgtm-monitoring"
      ]
      ssm_parameter_paths = [
        "/lgtm/swarm/*",
      ]
    }

    monitoring-1 = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["monitoring"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.master,dependency.security_groups.outputs.sg_ids.swarm,dependency.security_groups.outputs.sg_ids.monitoring,dependency.security_groups.outputs.sg_ids.postgresql]
      user_data_filename = local.variables.user_data_files["master"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["monitoring"]
      s3_bucket_arns = [
        "arn:aws:s3:::lgtm-monitoring"
      ]
      ssm_parameter_paths = [
        "/lgtm/swarm/*",
      ]
      cluster_identifier = "monitoring-cluster"
    }

    db-dev-staging-1 = {
      ami_instance       = local.variables.ubuntu_noble_ami
      instance_type      = local.variables.instance_types["database-dev-staging"]
      subnet_id          = dependency.vpc.outputs.subnet_public_all[0]
      security_group_ids = [dependency.security_groups.outputs.sg_ids.worker,dependency.security_groups.outputs.sg_ids.swarm,dependency.security_groups.outputs.sg_ids.monitoring,dependency.security_groups.outputs.sg_ids.postgresql]
      user_data_filename = local.variables.user_data_files["worker"]
      user_data_vars = {
        priv_key = local.priv_key
      }
      volume_size       = local.variables.volumes["database-dev-staging"]
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
