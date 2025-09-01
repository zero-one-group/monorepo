###################
###### INPUT #####
###################

variable "eip_name" {
  description = "Name prefix for Elastic IPs"
  type        = string
}

variable "common_tags" {
  description = "Common tags"
  type        = map(string)
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "elastic_ips" {
  description = "Map of Elastic IPs to create. Each key is the suffix, and value should be an empty object."
  type        = map(map(string))
  default     = {}
}

###################
###### OUTPUT #####
###################

output "eip_ids" {
  description = "Map of Elastic IP IDs keyed by suffix"
  value       = { for k, v in aws_eip.this : k => v.id }
}

output "eip_public_ips" {
  description = "Map of Elastic IP public IPs keyed by suffix"
  value       = { for k, v in aws_eip.this : k => v.public_ip }
}

# Backward compatibility outputs
output "eip_bastion_id" {
  description = "ID of the bastion Elastic IP if created"
  value       = lookup(aws_eip.this, "bastion", null) != null ? aws_eip.this["bastion"].id : null
}

output "eip_bastion_public_ip" {
  description = "Public IP of the bastion Elastic IP if created"
  value       = lookup(aws_eip.this, "bastion", null) != null ? aws_eip.this["bastion"].public_ip : null
}

output "eip_nginx_dev_staging_id" {
  description = "ID of the nginx Elastic IP Dev/Staging if created"
  value       = lookup(aws_eip.this, "nginx-dev-staging", null) != null ? aws_eip.this["nginx-dev-staging"].id : null
}

output "eip_nginx_dev_staging_public_ip" {
  description = "Public IP of the nginx Elastic IP Dev/Staging if created"
  value       = lookup(aws_eip.this, "nginx-dev-staging", null) != null ? aws_eip.this["nginx-dev-staging"].public_ip : null
}
