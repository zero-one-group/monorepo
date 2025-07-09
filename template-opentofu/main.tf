terraform {
  backend "s3" {
    bucket         = "lgtm-bucket-states"
    key            = "state/terraform.tfstate"
    region         = "ap-southeast-1"
    encrypt        = true
    dynamodb_table = "lgtm-locks"
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9.0"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_region" "current" {}

provider "time" {}

# Create the S3 bucket for storing state
resource "aws_s3_bucket" "state_bucket" {
  bucket = "${local.prefix}-bucket-states"

  lifecycle {
    prevent_destroy = true
  }

  tags = local.common_tags
}

# Configure server-side encryption using a separate resource
resource "aws_s3_bucket_server_side_encryption_configuration" "state_bucket_encryption" {
  bucket = aws_s3_bucket.state_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Create the DynamoDB table for state locking
resource "aws_dynamodb_table" "state_lock" {
  name         = "${local.prefix}-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  server_side_encryption {
    enabled = true
    # kms_key_arn = null  # Optional: Use AWS managed CMK (default behavior)
  }

  lifecycle {
    prevent_destroy = true
  }

  tags = local.common_tags
}

###################
####### VPC #######
###################
# Running once for new vpc use for all only in base state & lock
# module "vpc-base" {
#   source = "./modules/vpc-base"

#   vpc_name    = "${local.prefix}"
#   cidr_block  = var.vpc_cidr_block
#   region_name = data.aws_region.current.name
#   common_tags = local.common_tags
# }

module "vpc" {
  source = "./modules/vpc"

  vpc_name    = local.prefix
  cidr_block  = var.vpc_cidr_block
  region_name = data.aws_region.current.name
  common_tags = local.common_tags
}
##############################
####### SECURITY GROUP #######
##############################
module "sg_lb" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-lb"
  description = "Allow access to RDS PostgreSQL"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP"
    },
    {
      protocol    = "tcp"
      from_port   = 443
      to_port     = 443
      cidr_blocks = ["0.0.0.0/0"]
      description = "HTTPS"
      }, {
      protocol                 = "tcp"
      from_port                = 22
      to_port                  = 22
      source_security_group_id = module.sg_bastion.sg_id
      description              = "SSH access from bastion"
  }]
  egress_rule = [{
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }]
  common_tags = local.common_tags
}

module "sg_pgcluster" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-pgcluster"
  description = "Allow access to RDS PostgreSQL"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol    = "tcp"
    from_port   = 5432
    to_port     = 5432
    cidr_blocks = ["10.201.0.0/20"]
  }]
  egress_rule = [{
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }]
  common_tags = local.common_tags
}


module "sg_docker_swarm" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-swarm-cluster"
  description = "Swarm port"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol    = "tcp"
    from_port   = 2377
    to_port     = 2377
    cidr_blocks = ["10.201.0.0/20"]
    description = "Swarm port"
    }, {
    protocol    = "tcp"
    from_port   = 7946
    to_port     = 7946
    cidr_blocks = ["10.201.0.0/20"]
    description = "Swarm port"
    }, {
    protocol    = "udp"
    from_port   = 7946
    to_port     = 7946
    cidr_blocks = ["10.201.0.0/20"]
    description = "Swarm port"
    }, {
    protocol    = "udp"
    from_port   = 4789
    to_port     = 4789
    cidr_blocks = ["10.201.0.0/20"]
    description = "Swarm port"
  }]
  egress_rule = []
  common_tags = local.common_tags
}

module "sg_monitoring" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-monitoring"
  description = "Monitoring port"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol    = "tcp"
    from_port   = 9100
    to_port     = 9100
    cidr_blocks = ["10.201.0.0/20"]
    description = "Node Exporter"
    }, {
    # Postgresql Exporter
    protocol    = "tcp"
    from_port   = 9187
    to_port     = 9187
    cidr_blocks = ["10.201.0.0/20"]
    description = "PSQL Exporter"
    },
    {
      # Loki
      protocol    = "tcp"
      from_port   = 3100
      to_port     = 3100
      cidr_blocks = ["10.201.0.0/20"]
      description = "Loki"
    },
    {
      # Portainer Agent
      protocol    = "tcp"
      from_port   = 9001
      to_port     = 9001
      cidr_blocks = ["10.201.0.0/20"]
      description = "Portainer Agent"
    },
    {
      # Grafana
      protocol                 = "tcp"
      from_port                = 3000
      to_port                  = 3000
      source_security_group_id = module.sg_lb.sg_id
      description              = "Grafana"
    },
    {
      # Portainer
      protocol                 = "tcp"
      from_port                = 9000
      to_port                  = 9000
      source_security_group_id = module.sg_lb.sg_id
      description              = "Portainer"
  }]
  egress_rule = []
  common_tags = local.common_tags
}

module "sg_bastion" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-bastion"
  description = "Control bastion inbound and outbound access"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }]
  egress_rule = [{
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }]
  common_tags = local.common_tags
}

# Swarm Master Security Group
module "sg_swarm_master" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-master"
  description = "Control master inbound and outbound access"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol                 = "tcp"
    from_port                = 22
    to_port                  = 22
    source_security_group_id = module.sg_bastion.sg_id
    description              = "SSH access from bastion"
    }, {
    # Apps(port depend on needs)
    protocol                 = "tcp"
    from_port                = 8080
    to_port                  = 8080
    source_security_group_id = module.sg_lb.sg_id
    description              = "API"
  }]
  egress_rule = [{
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }]
  common_tags = local.common_tags
}

module "sg_swarm_worker" {
  source = "./modules/security_group"

  sg_name     = "${local.prefix}-worker"
  description = "Control worker inbound and outbound access"
  vpc_id      = module.vpc.vpc_id
  ingress_rule = [{
    protocol                 = "tcp"
    from_port                = 22
    to_port                  = 22
    source_security_group_id = module.sg_bastion.sg_id
    description              = "SSH access from bastion"
  }, ]
  egress_rule = [{
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }]
  common_tags = local.common_tags
}

#######################
####### KEYPAIR #######
#######################
module "ec2_keypair" {
  source = "./modules/keypair"

  keyname         = local.prefix
  public_key_path = var.public_key_path
  common_tags     = local.common_tags
}

########################
####### EC2 ############
########################
# Bastion one for all environment
module "ec2_bastion" {
  source = "./modules/ec2"

  ami_instance   = var.ubuntu_noble_ami
  instance_name  = "${local.prefix}-bastion"
  instance_type  = var.instance_type["bastion"]
  subnet         = module.vpc.misc_subnet_public_all[0]
  security_group = [module.sg_bastion.sg_id]
  keyname        = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["bastion"]
    vars = {
      aws_region = data.aws_region.current.name
      priv_key   = file(pathexpand(var.priv_key))
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["bastion"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }
}

resource "aws_eip_association" "bastion_eip" {
  instance_id   = module.ec2_bastion.ec2_instance_id
  allocation_id = module.vpc.eip_bastion_id
}
########################################
#### DEVELOPMENT & STAGING BLOCK #######
########################################
module "ec2_nginx" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-nginx"
  instance_type = var.instance_type["nginx-dev-staging"]
  subnet        = module.vpc.misc_subnet_public_all[0]
  security_group = [
    module.sg_monitoring.sg_id,
    module.sg_lb.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["nginx"]
    vars = {
      aws_region = data.aws_region.current.name
      priv_key   = file(pathexpand(var.priv_key))
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["nginx"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }
}
resource "aws_eip_association" "nginx_eip" {
  instance_id   = module.ec2_nginx.ec2_instance_id
  allocation_id = module.vpc.eip_nginx_id
}

module "ec2_master_1" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-master-1"
  instance_type = var.instance_type["master-dev-staging"]
  subnet        = module.vpc.dev_subnet_public_all[0]
  security_group = [
    module.sg_swarm_master.sg_id,
    module.sg_docker_swarm.sg_id,
    module.sg_monitoring.sg_id,
    module.sg_pgcluster.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["master"]
    vars = {
      aws_region   = data.aws_region.current.name
      project_name = local.prefix
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["master"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }
  depends_on = [module.ec2_monitoring_1]
}

resource "time_sleep" "wait_for_master" {
  depends_on = [module.ec2_master_1]

  create_duration = "120s"
}

module "ec2_worker_dev" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-worker-dev"
  instance_type = var.instance_type["worker-dev-staging"]
  subnet        = module.vpc.dev_subnet_public_all[0]
  security_group = [
    module.sg_swarm_worker.sg_id,
    module.sg_docker_swarm.sg_id,
    module.sg_monitoring.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["worker"]
    vars = {
      aws_region   = data.aws_region.current.name
      project_name = local.prefix
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["worker"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }

  depends_on = [time_sleep.wait_for_master]
}

module "ec2_worker_staging" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-worker-staging"
  instance_type = var.instance_type["worker-dev-staging"]
  subnet        = module.vpc.staging_subnet_public_all[0]
  security_group = [
    module.sg_swarm_worker.sg_id,
    module.sg_docker_swarm.sg_id,
    module.sg_monitoring.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["worker"]
    vars = {
      aws_region   = data.aws_region.current.name
      project_name = local.prefix
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["worker"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }

  depends_on = [time_sleep.wait_for_master]
}

# One for all environment
module "ec2_monitoring_1" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-monitoring-1"
  instance_type = var.instance_type["monitoring"]
  subnet        = module.vpc.misc_subnet_public_all[0]
  security_group = [
    module.sg_swarm_worker.sg_id,
    module.sg_docker_swarm.sg_id,
    module.sg_monitoring.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["master"]
    vars = {
      aws_region   = data.aws_region.current.name
      project_name = local.prefix
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["monitoring"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }
}

module "ec2_db_1" {
  source = "./modules/ec2"

  ami_instance  = var.ubuntu_noble_ami
  instance_name = "${local.prefix}-db-1"
  instance_type = var.instance_type["database-dev-staging"]
  subnet        = module.vpc.misc_subnet_public_all[0]
  security_group = [
    module.sg_swarm_worker.sg_id,
    module.sg_docker_swarm.sg_id,
    module.sg_monitoring.sg_id,
    module.sg_pgcluster.sg_id
  ]
  keyname = module.ec2_keypair.keyname
  user_data = {
    filename = var.user_data_ec2["worker"]
    vars = {
      aws_region   = data.aws_region.current.name
      project_name = local.prefix
    }
  }
  common_tags = local.common_tags

  root_block_device = {
    volume_type           = "gp3"
    volume_size           = var.volume["database-dev-staging"]
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }

  depends_on = [time_sleep.wait_for_master]
}
###########################
####PRODUCTION BLOCK#######
###########################
# module "ec2_nginx_prod" {
#   source = "./modules/ec2"

#   ami_instance   = var.ubuntu_noble_ami
#   instance_name  = "${local.prefix}-nginx-prod"
#   instance_type  = var.instance_type["nginx"]
#   subnet         = module.vpc.misc_subnet_public_all[0]
#   security_group = [module.sg_lb.sg_id]
#   keyname        = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["nginx"]
#     vars = {
#       aws_region = data.aws_region.current.name
#       priv_key   = file(pathexpand(var.priv_key))
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["nginx"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }
# }
# resource "aws_eip_association" "nginx_eip" {
#   instance_id   = module.ec2_nginx.ec2_instance_id
#   allocation_id = module.vpc.eip_nginx_id
# }

# module "ec2_master_1_prod" {
#   source = "./modules/ec2"

#   ami_instance  = var.ubuntu_noble_ami
#   instance_name = "${local.prefix}-master-1-prod"
#   instance_type = var.instance_type["master"]
#   subnet        = module.vpc.prod_subnet_public_all[0]
#   security_group = [
#     module.sg_swarm_master.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id,
#     module.sg_pgcluster.sg_id
#   ]
#   keyname = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["master"]
#     vars = {
#       aws_region   = data.aws_region.current.name
#       project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["master"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }
# }

# resource "time_sleep" "wait_for_master" {
#   depends_on = [module.ec2_master_1]

#   create_duration = "120s"
# }

# module "ec2_master_2_prod" {
#   source = "./modules/ec2"

#   ami_instance   = var.ubuntu_noble_ami
#   instance_name  = "${local.prefix}-master-2-prod"
#   instance_type  = var.instance_type["master"]
#   subnet         = module.vpc.prod_subnet_public_all[1]
#   security_group = [
#     module.sg_swarm_master.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id
#   ]
#   keyname        = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["master_join"]
#       vars = {
#         aws_region = data.aws_region.current.name
#         project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["master"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }

#   depends_on = [time_sleep.wait_for_master]
# }

# module "ec2_master_3_prod" {
#   source = "./modules/ec2"

#   ami_instance   = var.ubuntu_noble_ami
#   instance_name  = "${local.prefix}-master-3-prod"
#   instance_type  = var.instance_type["master"]
#   subnet         = module.vpc.prod_subnet_public_all[2]
#   security_group = [
#     module.sg_swarm_master.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id
#   ]
#   keyname        = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["master_join"]
#       vars = {
#         aws_region = data.aws_region.current.name
#         project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["master"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }

#   depends_on = [time_sleep.wait_for_master]
# }

# module "ec2_worker_1_prod" {
#   source = "./modules/ec2"

#   ami_instance  = var.ubuntu_noble_ami
#   instance_name = "${local.prefix}-worker-1"
#   instance_type = var.instance_type["worker"]
#   subnet        = module.vpc.prod_subnet_public_all[0]
#   security_group = [
#     module.sg_swarm_worker.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id
#   ]
#   keyname = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["worker"]
#     vars = {
#       aws_region   = data.aws_region.current.name
#       project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["worker"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }

#   depends_on = [time_sleep.wait_for_master]
# }

# module "ec2_worker_2_prod" {
#   source = "./modules/ec2"

#   ami_instance  = var.ubuntu_noble_ami
#   instance_name = "${local.prefix}-worker-2"
#   instance_type = var.instance_type["worker"]
#   subnet        = module.vpc.prod_subnet_public_all[1]
#   security_group = [
#     module.sg_swarm_worker.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id
#   ]
#   keyname = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["worker"]
#     vars = {
#       aws_region   = data.aws_region.current.name
#       project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["worker"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }

#   depends_on = [time_sleep.wait_for_master]
# }

# module "ec2_db_1_prod" {
#   source = "./modules/ec2"

#   ami_instance  = var.ubuntu_noble_ami
#   instance_name = "${local.prefix}-db-1-prod"
#   instance_type = var.instance_type["database"]
#   subnet        = module.vpc.misc_subnet_public_all[0]
#   security_group = [
#     module.sg_swarm_worker.sg_id,
#     module.sg_docker_swarm.sg_id,
#     module.sg_monitoring.sg_id,
#     module.sg_pgcluster.sg_id
#   ]
#   keyname = module.ec2_keypair.keyname
#   user_data = {
#     filename = var.user_data_ec2["worker"]
#     vars = {
#       aws_region   = data.aws_region.current.name
#       project_name = local.prefix
#     }
#   }
#   common_tags = local.common_tags

#   root_block_device = {
#     volume_type           = "gp3"
#     volume_size           = var.volume["database"]
#     delete_on_termination = true
#     encrypted             = true
#     kms_key_id            = ""
#   }

#   depends_on = [time_sleep.wait_for_master]
# }
###################
####### ECR #######
###################
# module "ecr" {
#   source = "./modules/ecr"

#   repo_name                      = "${local.prefix}-image"
#   common_tags                    = local.common_tags
#   untagged_image_expiration_days = 1
# }
###################
####### S3 ########
###################
# module "lgtm_s3_bucket" {
#   source = "./modules/s3"

#   bucket_name           = "${local.prefix}-buckets-monitoring"
#   enable_versioning     = false
#   enable_lifecycle_rule = true
#   transition_days       = 30
#   expiration_days       = 90

#   tags = local.common_tags
# }
