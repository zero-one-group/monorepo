###################
###### INPUT #####
###################

variable "vpc_name" {
  description = "Name prefix for VPC resources"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway for private subnets"
  type        = bool
  default     = true
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
}

variable "region" {
  description = "AWS region"
  type        = string
}

###################
###### OUTPUT #####
###################

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_arn" {
  description = "ARN of the VPC"
  value       = aws_vpc.main.arn
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "vpc_default_network_acl_id" {
  description = "The ID of the default network ACL"
  value       = aws_vpc.main.default_network_acl_id
}

output "vpc_default_route_table_id" {
  description = "The ID of the default route table"
  value       = aws_vpc.main.default_route_table_id
}

output "vpc_main_route_table_id" {
  description = "The ID of the main route table associated with this VPC"
  value       = aws_vpc.main.main_route_table_id
}

# Internet Gateway Outputs
output "internet_gateway_id" {
  description = "ID of the Internet Gateway"
  value       = aws_internet_gateway.main.id
}

output "internet_gateway_arn" {
  description = "ARN of the Internet Gateway"
  value       = aws_internet_gateway.main.arn
}

# Public Subnet Outputs
output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = [aws_subnet.public_a.id, aws_subnet.public_b.id, aws_subnet.public_c.id]
}

output "public_subnet_cidr_blocks" {
  description = "List of public subnet CIDR blocks"
  value       = [aws_subnet.public_a.cidr_block, aws_subnet.public_b.cidr_block, aws_subnet.public_c.cidr_block]
}

output "public_subnet_a_id" {
  description = "ID of public subnet A"
  value       = aws_subnet.public_a.id
}

output "public_subnet_b_id" {
  description = "ID of public subnet B"
  value       = aws_subnet.public_b.id
}

output "public_subnet_c_id" {
  description = "ID of public subnet C"
  value       = aws_subnet.public_c.id
}

output "public_subnet_a_arn" {
  description = "ARN of public subnet A"
  value       = aws_subnet.public_a.arn
}

output "public_subnet_b_arn" {
  description = "ARN of public subnet B"
  value       = aws_subnet.public_b.arn
}

output "public_subnet_c_arn" {
  description = "ARN of public subnet C"
  value       = aws_subnet.public_c.arn
}

# Private Subnet Outputs
output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = [aws_subnet.private_a.id, aws_subnet.private_b.id, aws_subnet.private_c.id]
}

output "private_subnet_cidr_blocks" {
  description = "List of private subnet CIDR blocks"
  value       = [aws_subnet.private_a.cidr_block, aws_subnet.private_b.cidr_block, aws_subnet.private_c.cidr_block]
}

output "private_subnet_a_id" {
  description = "ID of private subnet A"
  value       = aws_subnet.private_a.id
}

output "private_subnet_b_id" {
  description = "ID of private subnet B"
  value       = aws_subnet.private_b.id
}

output "private_subnet_c_id" {
  description = "ID of private subnet C"
  value       = aws_subnet.private_c.id
}

output "private_subnet_a_arn" {
  description = "ARN of private subnet A"
  value       = aws_subnet.private_a.arn
}

output "private_subnet_b_arn" {
  description = "ARN of private subnet B"
  value       = aws_subnet.private_b.arn
}

output "private_subnet_c_arn" {
  description = "ARN of private subnet C"
  value       = aws_subnet.private_c.arn
}

# Route Table Outputs
output "public_route_table_id" {
  description = "ID of the public route table"
  value       = aws_route_table.public.id
}

output "private_route_table_id" {
  description = "ID of the private route table"
  value       = aws_route_table.private.id
}

output "public_route_table_arn" {
  description = "ARN of the public route table"
  value       = aws_route_table.public.arn
}

output "private_route_table_arn" {
  description = "ARN of the private route table"
  value       = aws_route_table.private.arn
}

# NAT Gateway Outputs
output "nat_gateway_id" {
  description = "ID of the NAT Gateway"
  value       = var.enable_nat_gateway ? aws_nat_gateway.main[0].id : null
}

output "nat_gateway_public_ip" {
  description = "Public IP of NAT Gateway"
  value       = var.enable_nat_gateway ? aws_nat_gateway.main[0].public_ip : null
}

output "nat_gateway_private_ip" {
  description = "Private IP of NAT Gateway"
  value       = var.enable_nat_gateway ? aws_nat_gateway.main[0].private_ip : null
}

output "nat_eip_id" {
  description = "ID of the NAT EIP"
  value       = var.enable_nat_gateway ? aws_eip.nat[0].id : null
}

output "nat_eip_public_ip" {
  description = "Public IP of NAT EIP"
  value       = var.enable_nat_gateway ? aws_eip.nat[0].public_ip : null
}

output "nat_eip_allocation_id" {
  description = "Allocation ID of NAT EIP"
  value       = var.enable_nat_gateway ? aws_eip.nat[0].allocation_id : null
}

# All Subnet IDs (for convenience)
output "all_subnet_ids" {
  description = "List of all subnet IDs (public + private)"
  value       = concat(
    [aws_subnet.public_a.id, aws_subnet.public_b.id, aws_subnet.public_c.id],
    [aws_subnet.private_a.id, aws_subnet.private_b.id, aws_subnet.private_c.id]
  )
}

output "all_public_subnet_ids" {
  description = "List of all public subnet IDs"
  value       = [aws_subnet.public_a.id, aws_subnet.public_b.id, aws_subnet.public_c.id]
}

output "all_private_subnet_ids" {
  description = "List of all private subnet IDs"
  value       = [aws_subnet.private_a.id, aws_subnet.private_b.id, aws_subnet.private_c.id]
}

# Availability Zones
output "availability_zones" {
  description = "List of availability zones used"
  value       = [
    aws_subnet.public_a.availability_zone,
    aws_subnet.public_b.availability_zone,
    aws_subnet.public_c.availability_zone
  ]
}

# Additional useful outputs for other modules
output "public_subnet_route_table_association_ids" {
  description = "List of public subnet route table association IDs"
  value       = [
    aws_route_table_association.public_a.id,
    aws_route_table_association.public_b.id,
    aws_route_table_association.public_c.id
  ]
}

output "private_subnet_route_table_association_ids" {
  description = "List of private subnet route table association IDs"
  value       = [
    aws_route_table_association.private_a.id,
    aws_route_table_association.private_b.id,
    aws_route_table_association.private_c.id
  ]
}
