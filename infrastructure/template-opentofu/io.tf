variable "project" {
  default     = "lgtm"
  type        = string
  description = "Project name"
}

variable "region" {
  default     = "ap-southeast-1"
  type        = string
  description = "Region use for all resource"
}

variable "vpc_cidr_block" {
  type        = string
  default     = "10.201.0.0/16"
  description = "VPC cidr block"
}

variable "db_username" {
  sensitive   = true
  description = "Username for RDS Postgres instance"
  type        = string
}

variable "db_password" {
  sensitive   = true
  description = "Password for RDS Postgres instance"
  type        = string
}

variable "deletion_protection" {
  description = "Enable/disable deletion protection for RDS"
  type        = bool
  default     = true
}

variable "skip_final_snapshot" {
  description = "Determine if final snapshot should be created before RDS deletion"
  type        = bool
  default     = false
}

variable "final_snapshot_identifier" {
  description = "The name of the final snapshot when destroying instance"
  type        = string
  default     = null
}

variable "multi_az" {
  description = "When true, enables Multi-AZ for RDS"
  type        = bool
  default     = false
}

variable "priv_key" {
  description = "The private key for SSH access"
  type        = string
}

variable "public_key_path" {
  type        = string
  default     = "~/.ssh/id_rsa.pub"
  description = "Path of public key for ec2 instance"
}

variable "ubuntu_jammy_ami" {
  // Please consult https://cloud-images.ubuntu.com/locator/ec2/
  default     = "ami-0e0ddf453092e1e37"
  type        = string
  description = "Ubuntu AMI on singapore"
}

variable "ubuntu_noble_ami" {
  // Please consult https://cloud-images.ubuntu.com/locator/ec2/
  default     = "ami-0b874c2ac1b5e9957"
  type        = string
  description = "Ubuntu AMI on singapore"
}

variable "instance_type" {
  description = "Instance types for different EC2 instances"
  type        = map(string)
  default = {
    bastion              = "t3a.micro"
    master               = "t3a.small"
    worker               = "t3a.small"
    master-dev-staging   = "t3a.micro"
    worker-dev-staging   = "t3a.micro"
    monitoring           = "t3a.medium"
    database             = "t3a.medium"
    database-dev-staging = "t3a.small"
    nginx                = "t3a.medium"
    nginx-dev-staging    = "t3a.small"
  }
}

variable "volume" {
  description = "Instance type"
  type        = map(number)
  default = {
    bastion              = 20
    master               = 20
    worker               = 20
    monitoring           = 50
    database             = 50
    database-dev-staging = 25
    nginx                = 25
  }
}

variable "user_data_ec2" {
  type = map(string)
  default = {
    bastion     = "user-data.sh"
    nginx       = "user-data-nginx.sh"
    master      = "user-data-master.sh"
    master_join = "user-data-master-join.sh"
    worker      = "user-data-worker.sh"
  }
  description = "Template file name for ec2 instance apps"
}

variable "root_domain" {
  type        = string
  default     = "zero-one.cloud"
  description = "Root domain of the project"
}

locals {
  prefix = var.project
  common_tags = {
    Project     = var.project
    Environment = terraform.workspace
    Contact     = "rafiqul@zero-one-group.com"
    ManagedBy   = "Terraform"
    Version     = "1.x.x"
    created-by  = "zero-one-group"
  }
}
