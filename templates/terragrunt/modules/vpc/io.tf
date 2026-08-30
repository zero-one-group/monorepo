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

variable "cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_offset" {
  description = "Offset for subnet CIDR calculation. Public subnets use offset+1..3, private subnets use offset+101..103."
  type        = number
  default     = 0
}

variable "enable_nat_gateway" {
  description = "Enable a NAT gateway so private subnets can reach the internet"
  type        = bool
  default     = true
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
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "internet_gateway_id" {
  description = "ID of the Internet Gateway"
  value       = aws_internet_gateway.main.id
}

output "subnet_public_all" {
  description = "List of public subnet IDs"
  value       = [aws_subnet.public_a.id, aws_subnet.public_b.id, aws_subnet.public_c.id]
}

output "subnet_public_cidr_blocks" {
  description = "List of public subnet CIDR blocks"
  value       = [aws_subnet.public_a.cidr_block, aws_subnet.public_b.cidr_block, aws_subnet.public_c.cidr_block]
}

output "subnet_private_all" {
  description = "List of private subnet IDs (one per AZ)"
  value       = [aws_subnet.private_a.id, aws_subnet.private_b.id, aws_subnet.private_c.id]
}

output "subnet_private_cidr_blocks" {
  description = "List of private subnet CIDR blocks"
  value       = [aws_subnet.private_a.cidr_block, aws_subnet.private_b.cidr_block, aws_subnet.private_c.cidr_block]
}

output "route_table_public_id" {
  description = "ID of the public route table"
  value       = aws_route_table.public.id
}

output "route_table_private_id" {
  description = "ID of the private route table"
  value       = var.enable_nat_gateway ? aws_route_table.private[0].id : null
}

output "nat_gateway_id" {
  description = "ID of the NAT gateway (null when enable_nat_gateway is false)"
  value       = var.enable_nat_gateway ? aws_nat_gateway.main[0].id : null
}

output "nat_public_ip" {
  description = "Public IP of the NAT gateway (null when enable_nat_gateway is false)"
  value       = var.enable_nat_gateway ? aws_eip.nat[0].public_ip : null
}
