###################
###### INPUT #####
###################
variable "vpc_name" {
  description = "VPC name"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "ID of the existing VPC"
  type        = string
}

variable "internet_gateway_id" {
  description = "ID of the existing Internet Gateway"
  type        = string
}

variable "subnet_offset" {
  description = "Offset for subnet CIDR calculation"
  type        = number
}

variable "common_tags" {
  description = "Common tags"
  type        = map(string)
}

variable "region" {
  description = "AWS region"
  type        = string
}

###################
###### OUTPUT #####
###################

output "vpc_id" {
  description = "ID of the VPC"
  value       = data.aws_vpc.main.id
}

output "subnet_public_all" {
  description = "List of public subnet IDs"
  value       = [aws_subnet.public_a.id, aws_subnet.public_b.id, aws_subnet.public_c.id]
}

output "subnet_public_cidr_blocks" {
  description = "List of public subnet CIDR blocks"
  value       = [aws_subnet.public_a.cidr_block, aws_subnet.public_b.cidr_block, aws_subnet.public_c.cidr_block]
}

output "route_table_public_id" {
  description = "ID of the public route table"
  value       = aws_route_table.public.id
}
