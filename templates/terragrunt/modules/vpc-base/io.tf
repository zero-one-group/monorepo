###################
###### INPUT #####
###################

variable "vpc_name" {
  description = "Name of the existing VPC"
  type        = string
}

variable "common_tags" {
  description = "Common tags"
  type        = map(string)
}

variable "cidr_blocks" {
  description = "CIDR block"
  type        = string
}

variable "region" {
  description = "Region names"
  type        = string
}

###################
###### OUTPUT #####
###################

output "vpc_id" {
  description = "ID of the existing VPC"
  value       = data.aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "CIDR block of the existing VPC"
  value       = data.aws_vpc.main.cidr_block
}

output "internet_gateway_id" {
  description = "ID of the existing Internet Gateway"
  value       = data.aws_internet_gateway.main.id
}

output "vpc_arn" {
  description = "ARN of the existing VPC"
  value       = data.aws_vpc.main.arn
}
