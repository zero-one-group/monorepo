include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "../../modules/security-group"
}

dependency "vpc" {
  config_path = "../vpc"
}

locals {
  root_config  = read_terragrunt_config(find_in_parent_folders("root.hcl"))
  project_name = local.root_config.locals.project_name
  environment  = basename(dirname(get_terragrunt_dir()))
}

inputs = {
  vpc_id = dependency.vpc.outputs.vpc_id

  security_groups = {
    bastion = {
      name        = "${local.project_name}-${local.environment}-bastion"
      description = "Security group for bastion host in development"
      ingress_rule = [
        {
          protocol    = "tcp"
          from_port   = 22
          to_port     = 22
          cidr_blocks = ["0.0.0.0/0"]
          description = "SSH access"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    worker = {
      name        = "${local.project_name}-${local.environment}-worker"
      description = "Security group for worker nodes in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 22
          to_port                     = 22
          source_security_group_key   = "bastion"
          description                 = "SSH access from bastion"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    master = {
      name        = "${local.project_name}-${local.environment}-master"
      description = "Security group for master nodes in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 22
          to_port                     = 22
          source_security_group_key   = "bastion"
          description                 = "SSH access from bastion"
        },
        {
          protocol                    = "tcp"
          from_port                   = 8080
          to_port                     = 8080
          source_security_group_key   = "lb"
          description                 = "Apps ports adjust based on needs"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    swarm = {
      name        = "${local.project_name}-${local.environment}-swarm"
      description = "Security group for swarm in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 2377
          to_port                     = 2377
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Swarm"
        },
        {
          protocol                    = "tcp"
          from_port                   = 7946
          to_port                     = 7946
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Swarm"
        },
        {
          protocol                    = "udp"
          from_port                   = 7946
          to_port                     = 7946
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Swarm"
        },
        {
          protocol                    = "udp"
          from_port                   = 4789
          to_port                     = 4789
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Swarm"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    postgresql = {
      name        = "${local.project_name}-${local.environment}-postgresql"
      description = "Security group for postgresql in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 5432
          to_port                     = 5432
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Postgresql"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    lb = {
      name        = "${local.project_name}-${local.environment}-lb"
      description = "Security group for lb in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 80
          to_port                     = 80
          cidr_blocks                 = ["0.0.0.0/0"]
          description                 = "HTTP"
        },
        {
          protocol                    = "tcp"
          from_port                   = 443
          to_port                     = 443
          cidr_blocks                 = ["0.0.0.0/0"]
          description                 = "HTTPS"
        },
        {
          protocol                    = "tcp"
          from_port                   = 22
          to_port                     = 22
          source_security_group_key   = "bastion"
          description                 = "SSH access from bastion"
        }
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }

    monitoring = {
      name        = "${local.project_name}-${local.environment}-monitoring"
      description = "Security group for monitoring in development"
      ingress_rule = [
        {
          protocol                    = "tcp"
          from_port                   = 9100
          to_port                     = 9100
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Node exporter"
        },
        {
          protocol                    = "tcp"
          from_port                   = 3100
          to_port                     = 3100
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Loki"
        },
        {
          protocol                    = "tcp"
          from_port                   = 9187
          to_port                     = 9187
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Postgresql exporter"
        },
        {
          protocol                    = "tcp"
          from_port                   = 9001
          to_port                     = 9001
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Portainer agent"
        },
        {
          protocol                    = "tcp"
          from_port                   = 3000
          to_port                     = 3000
          source_security_group_key   = "lb"
          description                 = "Grafana"
        },
        {
          protocol                    = "tcp"
          from_port                   = 9000
          to_port                     = 9000
          source_security_group_key   = "lb"
          description                 = "Portainer"
        },
        {
          protocol                    = "tcp"
          from_port                   = 4040
          to_port                     = 4040
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Pyroscope"
        },
        {
          protocol                    = "tcp"
          from_port                   = 4317
          to_port                     = 4317
          cidr_blocks                 = ["10.201.0.0/20"]
          description                 = "Tempo"
        },
      ]
      egress_rule = [
        {
          protocol    = "-1"
          from_port   = 0
          to_port     = 0
          cidr_blocks = ["0.0.0.0/0"]
          description = "Allow all outbound traffic"
        }
      ]
    }
  }
}
