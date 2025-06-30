###################
###### INPUT #####
###################

variable "vpc_name" {
  type = string
}

variable "common_tags" {
  description = "Common tags"
}

variable "cidr_block" {
  type = string
}

variable "region_name" {
  type = string
}

###################
###### OUTPUT #####
###################

output "vpc_id" {
  value = data.aws_vpc.main.id
}

# Production Outputs
output "prod_subnet_public_all" {
  description = "List of Production public subnet IDs"
  value       = [aws_subnet.prod_public_a.id, aws_subnet.prod_public_b.id, aws_subnet.prod_public_c.id]
}

output "prod_subnet_public_cidr_blocks" {
  description = "List of Production public subnet CIDR blocks"
  value       = [aws_subnet.prod_public_a.cidr_block, aws_subnet.prod_public_b.cidr_block, aws_subnet.prod_public_c.cidr_block]
}

output "prod_all_subnet_ids" {
  description = "List of all Production subnet IDs (public and private)"
  value       = [aws_subnet.prod_public_a.id, aws_subnet.prod_public_b.id, aws_subnet.prod_public_c.id]
}

output "prod_all_cidr_block_subnets" {
  description = "List of all Production CIDR blocks (public and private)"
  value       = [aws_subnet.prod_public_a.cidr_block, aws_subnet.prod_public_b.cidr_block, aws_subnet.prod_public_c.cidr_block]
}

# Miscellanous Outputs
output "misc_subnet_public_all" {
  description = "List of Miscellanous public subnet IDs"
  value       = [aws_subnet.misc_public_a.id, aws_subnet.misc_public_b.id, aws_subnet.misc_public_c.id]
}

output "misc_subnet_public_cidr_blocks" {
  description = "List of Monitoring public subnet CIDR blocks"
  value       = [aws_subnet.misc_public_a.cidr_block, aws_subnet.misc_public_b.cidr_block, aws_subnet.misc_public_c.cidr_block]
}

# Development Outputs
output "dev_subnet_public_all" {
  description = "List of Development public subnet IDs"
  value       = [aws_subnet.dev_public_a.id, aws_subnet.dev_public_b.id, aws_subnet.dev_public_c.id]
}

output "dev_subnet_public_cidr_blocks" {
  description = "List of Development public subnet CIDR blocks"
  value       = [aws_subnet.dev_public_a.cidr_block, aws_subnet.dev_public_b.cidr_block, aws_subnet.dev_public_c.cidr_block]
}

# Staging Outputs
output "staging_subnet_public_all" {
  description = "List of Staging public subnet IDs"
  value       = [aws_subnet.staging_public_a.id, aws_subnet.staging_public_b.id, aws_subnet.staging_public_c.id]
}

output "staging_subnet_public_cidr_blocks" {
  description = "List of Staging public subnet CIDR blocks"
  value       = [aws_subnet.staging_public_a.cidr_block, aws_subnet.staging_public_b.cidr_block, aws_subnet.staging_public_c.cidr_block]
}

output "eip_bastion_id" {
  description = "ID of the Elastic IP"
  value       = aws_eip.bastion.id
}

output "eip_bastion_public_ip" {
  description = "Public IP of the Elastic IP"
  value       = aws_eip.bastion.public_ip
}

output "eip_nginx_id" {
  description = "ID of the Elastic IP"
  value       = aws_eip.nginx.id
}

output "eip_nginx_public_ip" {
  description = "Public IP of the Elastic IP"
  value       = aws_eip.nginx.public_ip
}

