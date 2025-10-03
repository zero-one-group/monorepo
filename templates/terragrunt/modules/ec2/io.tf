variable "instances" {
  description = "Map of EC2 instance configurations"
  type = map(object({
    ami_instance           = string
    instance_type          = string
    subnet_id              = string
    security_group_ids     = list(string)
    user_data_filename     = string
    user_data_vars         = optional(map(string), {})
    volume_size            = number
    volume_type            = optional(string, "gp3")
    enable_eip             = optional(bool, false)
    eip_allocation_id      = optional(string, "")
    hostname               = optional(string, "")
    s3_bucket_arns         = optional(list(string), [])
    ssm_parameter_paths    = optional(list(string), [])
    cluster_identifier     = optional(string, "default-cluster")
  }))
}

variable "keyname" {
  type        = string
  description = "Key name for instances"
}

variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

# Outputs
output "instance_public_ips" {
  description = "Map of instance public IPs"
  value       = { for k, v in aws_instance.instances : k => v.public_ip }
}

output "instance_private_ips" {
  description = "Map of instance private IPs"
  value       = { for k, v in aws_instance.instances : k => v.private_ip }
}

output "instance_ids" {
  description = "Map of instance IDs"
  value       = { for k, v in aws_instance.instances : k => v.id }
}

output "instance_public_dns" {
  description = "Map of instance public DNS names"
  value       = { for k, v in aws_instance.instances : k => v.public_dns }
}
